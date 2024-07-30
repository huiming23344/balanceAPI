package onePass

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huiming23344/balanceapi/db"
	"io"
	"net/http"
	"sync"
	"time"
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
	c.JSON(200, gin.H{
		"msg":       "ok",
		"code":      200,
		"requestId": c.Request.Header.Get("X-KSY-REQUEST-ID"),
	})

	go doBatchPay(body)
	return
}

func UserTrade(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})
}

type queryUserAmountResponse struct {
	Code      int                   `json:"code"`
	Msg       string                `json:"msg"`
	RequestID string                `json:"requestId"`
	Data      []queryUserAmountData `json:"data"`
}

type queryUserAmountData struct {
	Uid    int64   `json:"uid"`
	Amount float64 `json:"amount"`
}

func QueryUserAmount(c *gin.Context) {
	var body []int64
	if err := c.ShouldBind(&body); err != nil {
		// 如果解析失败，返回错误信息。
		c.JSON(400, gin.H{
			"error": "Invalid JSON",
		})
		return
	}
	var data []queryUserAmountData
	for _, uid := range body {
		amount, err := db.GetBalance(uid)
		if err != nil {
			amount = 0
		}
		data = append(data, queryUserAmountData{
			Uid:    uid,
			Amount: float64(amount) / 100,
		})
	}
	c.JSON(200, queryUserAmountResponse{
		Code:      200,
		Msg:       "ok",
		RequestID: c.Request.Header.Get("X-KSY-REQUEST-ID"),
		Data:      data,
	})
	return
}

type finishJson struct {
	BatchPayId string `json:"batchPayId"`
}

func batchPayFinish(reqUuid, requestId string, ch chan int, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			ch <- 400
			return
		default:
			data := finishJson{
				BatchPayId: requestId,
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Println("Error marshalling JSON: ", err)
			}
			reqBoday := bytes.NewBuffer(jsonData)
			req, err := http.NewRequest("POST", "http://120.92.116.60/thirdpart/onePass/batchPayFinish", reqBoday)
			if err != nil {
				fmt.Println("Error creating request: ", err)
			}

			req.Header.Set("X-KSY-REQUEST-ID", reqUuid)
			req.Header.Set("X-KSY-KINGSTAR-ID", "20004")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request: ", err)
				ch <- 0
				return
			}
			defer resp.Body.Close()
			// 读取响应体
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body: ", err)
				ch <- 0
				return
			}
			ch <- resp.StatusCode
			fmt.Println("Response status code:", resp.Status)
			fmt.Println("Response body:", string(body))
		}
	}
}

func doBatchPay(body batchPayJson) {
	// TODO: make sure each batchPayId will only do once
	wg := sync.WaitGroup{}
	for _, uid := range body.Uids {
		wg.Add(1)
		go func(uid int64) {
			amount, err := getAllFund(uid)
			if err != nil {
				wg.Done()
				return
			}
			db.AddMoney(uid, amount)
			wg.Done()
		}(uid)
	}
	wg.Wait()
	fmt.Println(db.GetAllBalance())
	ch := make(chan int)

	uniqueId := uuid.New().String()
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel() // 确保在函数退出时取消上下文
	for {
		go batchPayFinish(uniqueId, body.BatchPayId, ch, ctx)
		select {
		case code := <-ch:
			switch code {
			case 200:
				return
			default:
				continue
			}
		}
	}
}
