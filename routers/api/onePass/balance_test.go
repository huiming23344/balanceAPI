package onePass

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestBatchPay(t *testing.T) {
	iF := []Fund{
		{
			Uid:    100032,
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

func TestBatchPayFromFile(t *testing.T) {
	iF := []Fund{}
	jsonData, err := os.ReadFile("../../../testfile/initFund1000.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %s", err)
	}
	err = json.Unmarshal(jsonData, &iF)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %s", err)
	}
	initFunds(iF)
}
