package onePass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	time2 "time"
)

func TestBatchPay(t *testing.T) {
	iF := []Fund{
		{
			Uid:    100001,
			Amount: 88.91,
		},
		{
			Uid:    100042,
			Amount: 10000.93,
		},
		{
			Uid:    403131,
			Amount: 2345.35,
		},
		{
			Uid:    100052,
			Amount: 88.93,
		},
	}
	initFunds(iF)
}

func TestBatchPayOnce(t *testing.T) {
	iF := []Fund{
		{
			Uid:    100001,
			Amount: 36.73,
		},
	}
	initFunds(iF)
	uids := []int64{}
	for _, f := range iF {
		uids = append(uids, f.Uid)
	}
	payFundsAPI(uids)
}

func TestBatchPayFromFile(t *testing.T) {
	iF := []Fund{}
	jsonData, err := os.ReadFile("../../../testfile/initFund100.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %s", err)
	}
	err = json.Unmarshal(jsonData, &iF)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %s", err)
	}
	initFunds(iF)
	//// pay all funds
	//uids := []int64{}
	//for _, f := range iF {
	//	uids = append(uids, f.Uid)
	//}
	//payFunds(uids)
	//fmt.Println(db.GetBalance(100002)) // 2302047
}

var iF []Fund

func TestUserTrade(t *testing.T) {
	// init the funds
	iF = []Fund{}
	jsonData, err := os.ReadFile("../../../testfile/initFund100.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %s", err)
	}
	err = json.Unmarshal(jsonData, &iF)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %s", err)
	}
	initFunds(iF)
	// pay all funds
	uids := []int64{}
	for _, f := range iF {
		uids = append(uids, f.Uid)
	}

	payFundsAPI(uids)
	// transfer the funds
	//transferFundsToOneAccount(iF)
	//getFundAccount([]int64{100001})
}

func getFundAccount(uids []int64) {
	uniqueId := uuid.New().String()
	jsonData, err := json.Marshal(uids)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %s", err)
	}
	reqBody := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", "http://127.0.0.1:20004/onePass/queryUserAmount", reqBody)
	if err != nil {
		log.Fatalf("Error creating request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-KSY-REQUEST-ID", uniqueId)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %s", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
	}
	fmt.Println("Response status code:", resp.Status)
	fmt.Println("Response body:", string(body))
}

func payFundsAPI(uids []int64) {
	uniqueId := uuid.New().String()
	data := batchPayJson{
		BatchPayId: uniqueId,
		Uids:       uids,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %s", err)
	}
	reqBody := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", "http://127.0.0.1:20004/onePass/batchPay", reqBody)
	if err != nil {
		log.Fatalf("Error creating request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-KSY-REQUEST-ID", uniqueId)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %s", err)
	}
	defer resp.Body.Close()
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println("Error reading response body: ", err)
	//}
	//fmt.Println("Response status code:", resp.Status)
	//fmt.Println("Response body:", string(body))
}

func transferFundsToOneAccount(funds []Fund) {
	timeStart := time2.Now()
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
	fmt.Println("Transfer time: ", time2.Since(timeStart))
}

func transferApi(from, to int64, amount float64) error {
	data := userTradeJson{
		SourceUid: from,
		TargetUid: to,
		Amount:    amount,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %s", err)
	}
	reqBody := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", "http://127.0.0.1:20004/onePass/userTrade", reqBody)
	if err != nil {
		log.Fatalf("Error creating request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	uniqueId := uuid.New().String()
	req.Header.Set("X-KSY-REQUEST-ID", uniqueId)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %s", err)
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return err
	}
	fmt.Println("Response status code:", resp.Status)
	fmt.Println("Response body:", string(body))
	return nil
}
