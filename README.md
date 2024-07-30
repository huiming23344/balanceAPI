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




