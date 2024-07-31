package onePass

import (
	"context"
	"fmt"
	"testing"
)

func TestGetPay(t *testing.T) {
	var uid int64 = 600001
	var amount int64 = 1
	ch := make(chan int)
	ctx := context.Background()
	getPay(uid, amount, "aaaaa", ctx, ch)
	fmt.Println(<-ch)
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
	ans, err := getAllFund(600004)
	if err != nil {
		fmt.Println(err)
	}
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
