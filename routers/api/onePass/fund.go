package onePass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type getFundJson struct {
	TransactionId string  `json:"transactionId"`
	Uid           int64   `json:"uid"`
	Amount        float64 `json:"amount"`
}

type getFundResponse struct {
	Code      int    `json:"code"`
	RequestId string `json:"requestId"`
	Msg       string `json:"msg"`
	Data      string `json:"data"`
}

type Fund struct {
	Uid    int64   `json:"uid"`
	Amount float64 `json:"amount"`
}

var maxRequestParallel chan struct{} = make(chan struct{}, 100)

func getPay(uid int64, amount int64, uniqueId string, ch chan int) {

	fmt.Printf("GETPAY uid: %d, amount: %d, uniqueID: %s\n", uid, amount, uniqueId)
	amt := float64(amount) / 100
	data := getFundJson{
		TransactionId: uniqueId,
		Uid:           uid,
		Amount:        amt,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(fmt.Sprintf("Error marshalling JSON: %s", err))
		ch <- 0
		return
	}
	reqBody := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", "http://120.92.116.60/thirdpart/onePass/pay", reqBody)
	if err != nil {
		log.Println("Error creating request: ", err)
		ch <- 0
		return
	}
	reqUuid := uuid.New().String()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-KSY-REQUEST-ID", reqUuid)
	req.Header.Set("X-KSY-KINGSTAR-ID", "20004")

	// 发起请求
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

	if resp.StatusCode != 200 {
		ch <- 1
		return
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body: ", err)
		ch <- 0
		return
	}

	var result getFundResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Error unmarshalling json: ", err)
		ch <- 0
		return
	}
	if result.RequestId != reqUuid {
		ch <- 1
		return
	}
	ch <- result.Code
	return
}

func initFunds(list []Fund) {
	jsonData, err := json.Marshal(list)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)
	}
	reqBody := bytes.NewBuffer(jsonData)
	log.Println(string(jsonData))
	req, err := http.NewRequest("POST", "http://120.92.116.60/thirdpart/onePass/initAccount", reqBody)
	if err != nil {
		log.Println("Error creating request: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-KSY-REQUEST-ID", "1")
	req.Header.Set("X-KSY-KINGSTAR-ID", "20004")

	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request: ", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing response body: ", err)
		}
	}(resp.Body)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body: ", err)
		return
	}

	// 打印响应体
	log.Println("Response status code:", resp.Status)
	log.Println("Response body:", string(body))
}

func getAllFund(uid int64) (int64, error) {
	//cfg := config.GlobalConfig()

	var pre, ans int64 = 500000, 0
	fmt.Println("before get all one amount")
	ans += getAllOneAmount(uid, int64(1000000), 500)

	// TODO: use timeout
	//timeOut := time.Duration(cfg.Server.RequestTimeOut) * time.Millisecond
	for pre >= 1 {
		ans += getAllOneAmount(uid, pre, 2)
		pre /= 2
	}
	return ans, nil
}

func getAllOneAmount(uid, amount int64, maxParallel int) int64 {
	ans := int64(0)
	wg := sync.WaitGroup{}
	isDone := false
	for i := 1; i <= maxParallel; i++ {
		if isDone {
			break
		}
		if maxParallel > 2 {
			if i < 30 && i != 1 {
				time.Sleep(time.Duration(10) * time.Millisecond)
			}
		}
		wg.Add(1)
		go func() {
			ans += singalGetPay(uid, amount)
			isDone = true
			wg.Done()
		}()
	}
	wg.Wait()
	return ans
}

func singalGetPay(uid, amount int64) int64 {
	var ans int64 = 0
	uniqueId := uuid.New().String()
	for {
		ch := make(chan int)
		go getPay(uid, amount, uniqueId, ch)
		select {
		case code := <-ch:
			switch code {
			case 200:
				ans += amount
				uniqueId = uuid.New().String()
				continue
			case 501:
				return ans
			case 1:
				continue
			case 404:
				return 0
			}
			// TODO: use config
		case <-time.After(time.Duration(800) * time.Millisecond):
			//fmt.Println("timeout")
			<-ch
			continue
		}
	}
}
