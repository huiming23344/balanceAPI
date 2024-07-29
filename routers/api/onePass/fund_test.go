package onePass

import (
	"fmt"
	"testing"
)

func TestGetPay(t *testing.T) {
	var uid int64 = 600001
	var amount int64 = 1
	getPay(uid, amount, "aaaaa")
}

func TestGetAllFund(t *testing.T) {
	iF := []Fund{
		{
			Uid:    600001,
			Amount: 88.91,
		},
		{
			Uid:    600002,
			Amount: 10000.93,
		},
	}
	initFunds(iF)
	ans, _ := getAllFund(600002)
	fmt.Println(ans)
}

func TestInitFund(t *testing.T) {
	iF := []Fund{
		{
			Uid:    600001,
			Amount: 88.91,
		},
		{
			Uid:    600002,
			Amount: 10000.93,
		},
	}
	initFunds(iF)
}
