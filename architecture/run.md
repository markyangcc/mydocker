# run

执行 `mydocker run` 架构如下，
```text
[Run命令] 
     │
     ├─ [Init命令] 
     │           │
     │           └─ syscall.exec → [Stress命令] → 运行结束 → 返回状态
     │
     └─ 等待 → 接收状态 → 结束
```


