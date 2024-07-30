package uuidCache

import "sync"

type uuidCache struct {
	batchPay sync.Map
	trade    sync.Map
}

var uuidCacheInstance *uuidCache

func init() {
	uuidCacheInstance = &uuidCache{
		batchPay: sync.Map{},
		trade:    sync.Map{},
	}
}

func CheckAndAddBatchPay(uuid string) bool {
	_, ok := uuidCacheInstance.batchPay.LoadOrStore(uuid, struct{}{})
	return !ok
}

func CheckAndAddTrade(uuid string) bool {
	_, ok := uuidCacheInstance.trade.LoadOrStore(uuid, struct{}{})
	return !ok
}
