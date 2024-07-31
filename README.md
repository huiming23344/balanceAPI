# BalanceAPI
[![wakatime](https://wakatime.com/badge/user/2a98216d-462c-465e-b3a8-fcfb22e79aac/project/0da5d80c-4904-487d-86c6-79bfc06c51df.svg)](https://wakatime.com/badge/user/2a98216d-462c-465e-b3a8-fcfb22e79aac/project/0da5d80c-4904-487d-86c6-79bfc06c51df)


## Benchmarks
### onePass/batchPay
计时方法：
```go
var timeStart time.Time
func BatchPay(c *gin.Context) {
    timeStart = time.Now()
    ...
    go BatchPay()
    return
}
 
func BatchPay() {
    ...
    batchPayFinish()
    fmt.Printf("use time: %v\n", time.Since(timeStart))
    return
}
```
`testfile/ininFund100.json`, 100个账户数据 `use time: 1.2026825s`

在超过500个用户数据时，服务端会返回成功的报文，但是不会注册账户，所以没有进行更多的测试。


## TODO
### payFund数据异常

不稳定复现

在当前的getFund中，会出现数据多加5000的情况，也就是二分法第一个尝试会返回成功，7.30日22时以后出现问题，7.31早2时以后无问题，后续无法复现

### 使用ctx的性能问题定位

为什么使用ctx后，性能会下降？

在 commit`86bc839fa0fb2bd6476a407f3f515fc35fa6db95` func`BatchPay`中，使用`ctx`后，性能会下降到30s左右

#### 定位问题

- 限制并发数量：只有在5个goroutine并发时才不会出现性能问题，但是因为并发度有限，时间还是在20s左右
- 经过测试不是db并发读写锁导致的性能问题，db的一次读写只有不到1ms。
- ctx的超时机制？在log中似乎观察到了ctx超时的情况，但没有退出goroutine的情况

### 使用pprof进行性能分析

