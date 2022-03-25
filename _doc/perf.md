# 性能测试数据

测试环境：

腾讯云，轻量应用服务器，CPU: 2核 内存: 4GB 80GB SSD云硬盘 CentOS 8.2 64bit

| 方式     | TPS   | 平均时延    | P99      |
|--------|-------|---------|----------|
| 正常记录日志 | 10907 | 9.161ms | 20.793ms |
| 关闭记录日志 | 56617 | 1.761ms | 7.917ms  |
| 异步记录日志 | 23709 | 4.199ms | 19.058ms |

## 正常记录日志

启动目标程序：

```sh
[d5k@VM-24-15-centos ~]$ fastrest
```

压测

```sh
[d5k@VM-24-15-centos ~]$ berf :14142 -d1m
Berf benchmarking http://127.0.0.1:14142/ for 1m0s using 100 goroutine(s), 2 GoMaxProcs.

Summary:
  Elapsed              1m0.001s
  Count/RPS    654452 10907.222
    200        654452 10907.222
  ReadWrite  28.882 13.789 Mbps

Statistics    Min      Mean    StdDev     Max
  Latency    115µs   9.161ms   4.154ms  54.769ms
  RPS       9785.59  10907.31   238.5   11350.36

Latency Percentile:
  P50        P75       P90      P95       P99      P99.9     P99.99
  8.871ms  11.377ms  14.62ms  16.528ms  20.793ms  26.785ms  39.438ms
```

## 关闭记录日志

启动目标程序：

```sh
[d5k@VM-24-15-centos ~]$ LOG_TYPE=off fastrest
```

压测

```sh
[d5k@VM-24-15-centos ~]$ berf :14142 -d1m
Berf benchmarking http://127.0.0.1:14142/ for 1m0s using 100 goroutine(s), 2 GoMaxProcs.

Summary:
  Elapsed                   1m0s
  Count/RPS    3397119 56617.918
    200        3397119 56617.918
  ReadWrite  149.924 71.567 Mbps

Statistics    Min      Mean    StdDev     Max
  Latency    21µs    1.761ms   1.708ms  95.142ms
  RPS       47624.3  56556.95  2065.36  59672.67

Latency Percentile:
  P50        P75      P90     P95      P99     P99.9    P99.99
  1.467ms  1.833ms  2.627ms  3.79ms  7.917ms  18.897ms  50.48ms
```

## 异步记录日志

启动目标程序：

```sh
[d5k@VM-24-15-centos ~]$ LOG_TYPE=async fastrest
```

压测

```sh
[d5k@VM-24-15-centos ~]$ berf :14142 -d1m
Berf benchmarking http://127.0.0.1:14142/ for 1m0s using 100 goroutine(s), 2 GoMaxProcs.

Summary:
  Elapsed              1m0.005s
  Count/RPS   1422704 23709.504
    200       1422704 23709.504
  ReadWrite  62.783 29.971 Mbps

Statistics    Min      Mean    StdDev     Max
  Latency     31µs    4.199ms  3.605ms  62.343ms
  RPS       17225.26  23697.4  1513.39  25870.2

Latency Percentile:
  P50        P75     P90      P95       P99      P99.9     P99.99
  3.133ms  3.776ms  8.25ms  11.328ms  19.058ms  32.761ms  46.517ms
```

