package onePass

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huiming23344/balanceapi/config"
	"github.com/huiming23344/balanceapi/db"
	"github.com/huiming23344/balanceapi/uuidCache"
	"io"
	"log"
	"math"
	"net/http"
	"sync"
	"time"
)

type batchPayJson struct {
	BatchPayId string  `json:"batchPayId"`
	Uids       []int64 `json:"uids"`
}

type queryUserAmountResponse struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestID string `json:"requestId"`
	Data      []Fund `json:"data"`
}

type finishJson struct {
	BatchPayId string `json:"batchPayId"`
}

type userTradeJson struct {
	SourceUid int64   `json:"sourceUid"`
	TargetUid int64   `json:"targetUid"`
	Amount    float64 `json:"amount"`
}

var timeStart time.Time

func BatchPay(c *gin.Context) {
	timeStart = time.Now()
	var body batchPayJson
	if err := c.ShouldBind(&body); err != nil {
		// 如果解析失败，返回错误信息。
		c.JSON(400, gin.H{
			"error": "Invalid JSON",
		})
		return
	}
	if !uuidCache.CheckAndAddBatchPay(body.BatchPayId) {
		c.JSON(400, gin.H{
			"error": "batchPayId already exist",
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
	// TODO: make sure each requestId will only do once
	if !uuidCache.CheckAndAddTrade(c.Request.Header.Get("X-KSY-REQUEST-ID")) {
		c.JSON(400, gin.H{
			"error": "requestId already exist",
		})
		return
	}
	var body userTradeJson
	if err := c.ShouldBind(&body); err != nil {
		// 如果解析失败，返回错误信息。
		c.JSON(400, gin.H{
			"error": "Invalid JSON",
		})
		return
	}
	amount := int64(math.Round(body.Amount * 100))
	err := db.Transfer(body.SourceUid, body.TargetUid, amount)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":       "ok",
		"code":      200,
		"requestId": c.Request.Header.Get("X-KSY-REQUEST-ID"),
	})
	return
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
	var data []Fund
	for _, uid := range body {
		amount, err := db.GetBalance(uid)
		if err != nil {
			amount = 0
		}
		data = append(data, Fund{
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

func batchPayFinish(reqUuid, requestId string, ch chan int) {
	cfg := config.GlobalConfig()
	data := finishJson{
		BatchPayId: requestId,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
	}
	reqBoday := bytes.NewBuffer(jsonData)
	url := cfg.Server.ServerAddr + "/thirdpart/onePass/batchPayFinish"
	req, err := http.NewRequest("POST", url, reqBoday)
	if err != nil {
		log.Println("Error creating request: ", err)
	}

	req.Header.Set("X-KSY-REQUEST-ID", reqUuid)
	req.Header.Set("X-KSY-KINGSTAR-ID", "20004")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request: ", err)
		ch <- 0
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing response body: ", err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body: ", err)
		ch <- 0
		return
	}
	ch <- resp.StatusCode
	log.Println("Response status code:", resp.Status)
	log.Println("Response body:", string(body))
}

func doBatchPay(body batchPayJson) {
	payFunds(body.Uids)
	// call batchPayFinish when all user finish
	ch := make(chan int)
	uniqueId := uuid.New().String()
	for {
		go batchPayFinish(uniqueId, body.BatchPayId, ch)
		select {
		case code := <-ch:
			switch code {
			case 200:
				log.Printf("use time: %v\n", time.Since(timeStart))
				return
			default:
				continue
			}
		}
	}
}

func payFunds(uids []int64) {
	wg := sync.WaitGroup{}
	for _, uid := range uids {
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
}
