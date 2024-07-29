package onePass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
)

type Fund struct {
	TransactionId string  `json:"transactionid"`
	Uid           int64   `json:"uid"`
	Amount        big.Rat `json:"amount"`
}

func getPay(uid int64, amount big.Rat) {
	data := Fund{
		TransactionId: "aaa",
		Uid:           uid,
		Amount:        amount,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON: ", err)
		return
	}
	reqBody := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", "http://example.com", reqBody)
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-KSY-REQUEST-ID", "aaa")
	req.Header.Set("X-KSY-KINGSTAR-ID", "20004")

	// 设置请求头，指明发送的是 JSON 格式的数据
	req.Header.Set("Content-Type", "application/json")

	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return
	}

	// 打印响应体
	fmt.Println("Response status code:", resp.Status)
	fmt.Println("Response body:", string(body))
}
