# mydocker
自己动手写 docker，从零到 1 实现一个 mini docker

## 3.1 run
注意这里有一个坑，书中代码在 mount 一个全新的 procfs 到 /proc 后并没有 unmount，这就导致覆盖了系统上的 /proc 文件系统。这也就导致了 mydocker run 命令只能执行一次，第二次执行会报错 `/proc/self/exe` 路径不存在。


在 urfave cli 框架钩子实现 `umount /proc` 或者暂时移除书中这段代码。
```go
app.After = func(context *cli.Context) error {
    var err error
    if err = syscall.Unmount("/proc", syscall.MNT_DETACH); err != nil {
        logrus.Errorf("failed to umount /proc: %v", err)
    }
    return err
}
```

> Tips: 等我们实现到 ch4，会将这部分不合理的地方 fix 掉。目前可以使用这种不够优雅方式规避。

## 3.2 cgroup limit
cgroup v2 下的对应接口
* 内存限制: `memory.max`
* cpu限制: `cpu.max`
* cpuset限制:`cpuset.cpus`
* 进程加入 cgroup: `cgroup.procs`

具体默认值用法，参考内核文档cgroup-v2 文档: https://docs.kernel.org/admin-guide/cgroup-v2.html

这节也有一个坑：内存、cpu资源限制不生效？

书中实例是通过 stress 进行的。因为 stress 是多进程模型，默认通过创建多个子进程来模拟不同类型的负载，而主进程则调用 wait() 等待子进程退出。

代码是先通过 cmd.Start() 将 stress 进程运行起来，再将其 `parnet.Process.Pid`(主进程) 加入 cgroup 限制资源。而子进程并没有被加入 cgroup，所以会发现资源限制被没有生效。作者使用 stress 作为演示程序不是最好的选择，换成单进程或者多线程模型演示程序更好。

临时解决方案：

手动将另外的 stress 进程加入 cgroup，此时发现 cpu、cpuset 限制生效了
```shell
# pidof stress
72821 72816

# 将不在 cgroup 中的另一个进程加入 cgroup
# echo 72821 >> /sys/fs/cgroup/mydocker-cgroup/cgroup.procs
# cat /sys/fs/cgroup/mydocker-cgroup/cgroup.procs
72816
72821
```

内存不生效的原因则是，cgroup 内存限制指的新分配内存，也就是增量，因为 stress 再启动时已经完成内存分配，运行时不再分配内存，所以会看到内存限制不生效情况。

> Tips: 等我们实现到 ch3.3，会将这部分不合理的地方 fix 掉。目前我们理解资源限制原理和限制不生效的原因即可。

## 3.3 Pipe & Env
测试命令，
```shell
sudo ./mydocker run -ti --cpushare 8000 --cpuset 0-1 -m 50m /usr/bin/stress --vm-bytes 30m --vm-keep --vm 1
```

需要注意的点，必须关闭写端以触发 EOF，否则 io.ReadAll 将永久阻塞等待数据结束。​因为我们在读端调用的是 io.ReadAll() API，只有在读取到 EOF 才返回，这和调用 os.Read() API 不一样。
```golang
// reader
	data, err := io.ReadAll(pipe)


//writer
    pipe.Close()
```

**可喜可贺！** 本节通过让子进程 pending 住的方式（**此时进程已经在存在了，已经有 pid 了，虽然 `ps` 命令还看不到，可以通过 /proc 文件系统查看**））。在 pending 这段时间设置好容器 cgroup，相当于在容器进程没启动之前就提前限制好 cgroup 资源，接着子进程从管道中接受到父进程发来的民工。通过 syscall.Exec 替换自身，实现启动容器进程。这种方式解决了 3.2 提到的 stress worker 进程不会被纳入 cgroup 管理的问题。


## 4.1 rootfs
目标：需要为容器提供 rootfs, 让容器使用自己的 rootfs 启动。

本章节有一个可以优化的点，在执行 `pivotRoot` 后，推荐将容器内的挂载传播属性设为 MS_PRIVATE，隔离容器与宿主机的挂载点交互，避免容器内操作影响宿主机（如覆盖 /proc、/dev）

```go
if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
	return fmt.Errorf("setting mount propagation to private failed: %v", err)
}
```

## 4.2 aufs/overlayfs
因为 aufs 在我的实验环境上(RockyLinux 9.5, RHEL Kernel 5.14)已经不被 `mount` 支持了，所以本节就替换为使用 overlayfs 代替。
本节内容很简单，我们先用 linux 命令理解 overlayfs 的原理，再转化为 golang 代码即可，

```shell
mount -t overlay overlay -o lowerdir=<lowerdir1>:<lowerdir2>:<...>,upperdir=<upperdir>,workdir=<workdir> <mountpoint>
```
大概介绍下这里的三个参数：

`lowerdir`：只读目录，多个目录用 ":" 分隔。

`upperdir`：可读写目录，所有对文件的修改（新增/覆盖/删除）会直接作用于该层。当文件被修改时，OverlayFS 会将原文件从 lowerdir 复制到 upperdir，并修改在 upperdir 中的同名文件。

> 若删除，则是在 upperdir 中标记同名文件删除，lowerdir 中文件始终是不变的。（在下文的例子中尝试一下）

`workdir`：必须是空目录，用于 overlayfs 元数据处理。

overlayfs 的原理是，lowerdir 是只读的，当我们将 lowerdir & upperdir 共同 mount 到一个 mountpoint 之后，那么在 mountpoint 目录的任何更改都会直接作用读写层 upperdir。

例如下面的例子：

```shell
mkdir low1 low2 low3 upper work overlaydir
touch low1/1.txt low2/2.txt low3/3.txt

mount -t overlay overlay -o lowerdir=./low1:./low2:./low3,upperdir=./upper,workdir=./work overlaydir
```
执行完上面的命令后：

```shell
# ls overlaydir
1.txt 2.txt 3.txt
```
可以看到三个 lower 目录下的内容被合并到 overlaydir 中。

如果这时在 overlaydir 中创建一个新文件，那么这个文件会出现在 upperdir 中，lowerdir 保持不变
```shell
touch overlaydir/4.txt
```

试一试修改下 3.txt 的内容，发现修改会直接作用于 upperdir 中的 3.txt 文件
```shell
echo "hello world" > overlaydir/3.txt
```

试一试删除 2.txt 文件，发现 upperdir 中新出现了一个 2.txt 文件, 但是一个特殊的 0/0 的字符设备，这就是标记删除的方式
``` shell
# file 2.txt
2.txt: character special (0/0)
```

那么，overlayfs 的简单用法就这样。总结下来就是，我们将 lowerdir 和 upperdir 同意 overlayfs 的方式 mount 到一个 mountpoint 时候，后续对 mountpoint 的修改都会作用到 upperdir 上，lowerdir 是永远不变的。毕竟 overlayfs 的设计里 uppperdir 才是读写层。


下面开始正文，直接参阅代码即可，

1、创建 4 个文件夹，分别是 overlayfs 需要的 lowerdir,upperdir,workdir 和我们 mount 之后的容器 rootfs 目录，这里叫做 "merged"

2、将 busybox rootfs 拷贝 overlayfs lowerdir，调用 `mount -t overlay overlay` 即可~

测试命令，
```shell
sduo ./mydocker run -ti sh
```

## 4.3 volume
本节要实现的功能是将主机上的目录挂载到容器里，来实现数据持久化功能。

和上一节一样，我们先理解原理，再转化为 golang 代码来实现，

实现原理很简单，我们只需要在 overlayfs 挂载成功后，在将节点的目录通过 bind mount 挂载到 overlayfs 合并后的目录中即可
``` shell
# mount --bind [源目录] [目标挂载点]
# /overlayfs/merged 为 overlayfs 合并之后的目录
mkdir /overlayfs/merged/host/tmp
mount --bind /tmp /overlayfs/merged/host/tmp
```
解释：命令将节点的 `/tmp` 目录挂载到 overlayfs 里的 `/host/tmp` 目录。这样实现容器的 -v 将 host 目录映射到容器中的功能。

测试命令，将节点的 /tmp 挂载到容器中的 /host/tmp 目录
```shell
sudo ./mydocker run -ti -v /tmp:/host/tmp sh
```
