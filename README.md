# BalanceAPI
[![wakatime](https://wakatime.com/badge/user/2a98216d-462c-465e-b3a8-fcfb22e79aac/project/0da5d80c-4904-487d-86c6-79bfc06c51df.svg)](https://wakatime.com/badge/user/2a98216d-462c-465e-b3a8-fcfb22e79aac/project/0da5d80c-4904-487d-86c6-79bfc06c51df)

- 使用`gin`框架进行API的开发
  - `gin`框架每一个请求都会开启一个goroutine进行处理，所以不需要额外的goroutine管理就可以保证较高性能
- 数据部分使用`sync.Map`进行索引
  - map的value为读写锁保护的结构体，降低锁的粒度，保证并发读写的安全
- payFund的每个账户都会开启一个goroutine进行处理
  - 对于每个开起的goroutine
    - 使用channel控制全局worker的数量
    - 根据账户的金额不同动态的调整goroutine的数量
    - 对于小于一万的部分，开启2个goroutine进行处理

## 性能分析
通过Flame Graph进行性能分析，可以看到占用cpu时间较大的主要是两部分，`sync.Waitgroup`同步等待的时间和`runtime.kevent`主要是网络epoll的等待时间
后续的性能优化可以考虑异步的方式进行处理，减少`sync.Waitgroup`的同步等待

## Benchmarks
### onePass/batchPay
#### 单个大数据测试
单个一亿的账户处理时间5s左右，`payFunds use time:  4.374602958s`

单个十亿的账户处理时间50s左右，`payFunds use time:  50.013073s`
```go
// 测试数据
iF := []Fund{
    {
        Uid:    100001,
        Amount: 100000000.53,
    },
}
```
#### 多账户测试
使用`test-file/ininFund100.json`初始化funds, 100个账户数据获取支付的时间为 `use time: 1.2026825s`
```go
// 计时方法
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

在超过500个用户数据时，服务端会返回成功的报文，但是不会注册账户，所以没有进行更多的测试。

### onePass/transfers
使用外部API进行转账数据源为：`test-file/ininFund100.json`, 100个账户数据转账到同一个账户 `Transfer time:  24.312625ms`，
测试结果正确，除了id为100001的账户其他全部账户余额均为0，10001账户金额多次相同，且数额正确
```go
// 计时方法
func transferFundsToOneAccount(funds []Fund) {
	timeStart := time.Now()
	// transfer the funds to one account
	for _, f := range funds {
		if f.Uid == 100001 {
			continue
		}
		err := transferApi(f.Uid, 100001, f.Amount)
		if err != nil {
			log.Fatalf("Error transfering fund: %s", err)
		}
	}
	fmt.Println("Transfer time: ", time.Since(timeStart))
}
```