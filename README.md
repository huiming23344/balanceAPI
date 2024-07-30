# BalanceAPI



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
`testfile/ininFund1000.json`, 1000个账户数据 `use time: 1.5673455s`



