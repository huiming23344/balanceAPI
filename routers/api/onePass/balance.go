package onePass

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

type batchPayJson struct {
	BatchPayId string  `json:"batchPayId"`
	Uids       []int64 `json:"uids"`
}

func BatchPay(c *gin.Context) {
	var body batchPayJson
	if err := c.ShouldBind(&body); err != nil {
		// 如果解析失败，返回错误信息。
		c.JSON(400, gin.H{
			"error": "Invalid JSON",
		})
		return
	}
	// TODO: make sure a batchPayId will only do once
	mp := map[int64]int64{}
	wg := sync.WaitGroup{}
	for _, uid := range body.Uids {
		wg.Add(1)
		go func(uid int64) {
			amount, err := getAllFund(uid)
			if err != nil {
				wg.Done()
				return
			}
			mp[uid] = amount
			wg.Done()
		}(uid)
	}
	wg.Wait()
	fmt.Println(mp)
	// TODO: call POST http://120.92.116.60/thirdpart/onePass/batchPayFinish
}

func UserTrade(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}

func QueryUserAmount(c *gin.Context) {

}
