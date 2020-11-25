# UCAS-Parallel-And-Distributed-Computing-Labs
国科大研一课程并行与分布式计算配套实验。

## Lab 1
实现全局状态的快照算法，并监控下列程序：两个进程 P 和 Q 用两个通道连成一个环，它们不断地轮转消息 m。在任何一个时刻，系统中仅有一份 m 的拷贝。

每个进程的状态是指由它接收到 m 的次数。P 首先发送 m。在某一点，P 得到消息且它的状态是 101。在发送 m 之后，P 启动快照算法，要求记录由快照算法报告的全局状态。

## Lab 2
在 Linux 平台，利用开源 RPC 代码，实现锁服务(Lock Service)，锁服务包括两个模块：锁客户，锁服务器，两者通过 RPC 通信。

需要实现两个功能：
1. 客户发 acquire 请求，从锁服务器请求一个特定的锁，用 release 释放锁，锁服务器一次将锁授予一个客户。
2. 扩充 RPC 库，实现 at-most-once 执行语义，即消除重复的 RPC 请求。

## Environment
Ubuntu 18.04

`export GOPATH=~/snapshot`
