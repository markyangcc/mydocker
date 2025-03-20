# mydocker
自己动手写 docker，从零到 1 实现一个 mini docker

## 3.1 mydocker run
注意这里有一个坑，书中代码在 mount 一个全新的 procfs 到 /proc 后并没有 unmount，这就导致覆盖了系统上的 /proc 文件系统。这也就导致了 mydocker run 命令只能执行一次，第二次执行会报错 `/proc/self/exe` 路径不存在。

参考我的实现，在 urfave cli 框架钩子实现 `umount /proc`
```go
app.After = func(context *cli.Context) error {
    var err error
    if err = syscall.Unmount("/proc", syscall.MNT_DETACH); err != nil {
        logrus.Errorf("failed to umount /proc: %v", err)
    }
    return err
}
```
