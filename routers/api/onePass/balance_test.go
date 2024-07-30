package onePass

import "testing"

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
	}
	initFunds(iF)
}

func TestBatchPayFinish(t *testing.T) {
	batchPayFinish("aaaaaaa", "bnnnnaaaaaaan")
}
