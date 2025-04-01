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

可喜可贺！本节通过先创建子进程，然后 pending 住(从管道中接收需要执行的命令是阻塞的)，接着子进程通过 syscall.Exec 替换自己的方式实现创建真正的子进程。这种方式解决了 3.2 提到的 stress worker 进程不会被纳入 cgroup 管理的问题。
