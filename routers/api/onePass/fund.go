package onePass

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
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

func getPay(uid int64, amount int64, uniqueId string, ch chan int) {
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
		fmt.Println("Error creating request: ", err)
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
	fmt.Println("Response status code:", resp.Status)
	fmt.Println("Response body:", string(body))

	if resp.StatusCode != 200 {
		if resp.StatusCode == 504 {
			ch <- 1
			return
		}
		ch <- 0
		return
	}

	var result getFundResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error unmarshalling json: ", err)
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

type Fund struct {
	Uid    int64   `json:"uid"`
	Amount float64 `json:"amount"`
}

func initFunds(list []Fund) {
	jsonData, err := json.Marshal(list)
	if err != nil {
		fmt.Println("Error marshalling JSON: ", err)
	}
	reqBody := bytes.NewBuffer(jsonData)
	fmt.Println(string(jsonData))
	req, err := http.NewRequest("POST", "http://120.92.116.60/thirdpart/onePass/initAccount", reqBody)
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-KSY-REQUEST-ID", "1")
	req.Header.Set("X-KSY-KINGSTAR-ID", "20004")

	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return
	}

	// 打印响应体
	fmt.Println("Response status code:", resp.Status)
	fmt.Println("Response body:", string(body))
}

func getAllFund(uid int64) (int64, error) {
	//cfg := config.GlobalConfig()
	var pre, ans int64 = 500000, 0
	uniqueId := uuid.New().String()
	// TODO: use timeout
	//timeOut := time.Duration(cfg.Server.RequestTimeOut) * time.Millisecond
	for pre >= 1 {
		ch := make(chan int)
		go getPay(uid, pre, uniqueId, ch)
		select {
		case code := <-ch:
			switch code {
			case 200:
				ans += pre
				uniqueId = uuid.New().String()
				continue
			case 501:
				pre = pre / 2
				uniqueId = uuid.New().String()
				continue
			case 1:
				continue
			case 404:
				return 0, errors.New(fmt.Sprintf("not found account by uid: %d", uid))
			}
			// TODO: use config
		case <-time.After(time.Duration(100) * time.Millisecond):
			continue
		}
	}
	return ans, nil
}
