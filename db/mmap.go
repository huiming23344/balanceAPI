package db

import (
	"errors"
	"sync"
)

type balanceAccount struct {
	Uid     int64
	RwMutex sync.RWMutex
	Balance int64
}

type MMap struct {
	uidMap map[int64]*balanceAccount
}

func NewMMapEngine() *MMap {
	return &MMap{
		uidMap: make(map[int64]*balanceAccount),
	}
}

// Add will add balance to uid account
// if account do not exist just add a new account
func (m *MMap) addMoney(uid int64, amount int64) {
	if account, ok := m.uidMap[uid]; ok {
		account.RwMutex.Lock()
		account.Balance += amount
		account.RwMutex.Unlock()
	} else {
		account := &balanceAccount{
			Uid:     uid,
			Balance: amount,
			RwMutex: sync.RWMutex{},
		}
		m.uidMap[uid] = account
	}
}

func (m *MMap) getBalance(uid int64) (int64, error) {
	if account, ok := m.uidMap[uid]; ok {
		account.RwMutex.RLock()
		defer account.RwMutex.RUnlock()
		return account.Balance, nil

	} else {
		return 0, errors.New("can not find the account")
	}
}

// Transfer will transfer amount from 'from' account to 'to' account
func (m *MMap) transfer(from, to, amount int64) error {
	if fromAccount, ok := m.uidMap[from]; ok {
		fromAccount.RwMutex.Lock()
		defer fromAccount.RwMutex.Unlock()
		if fromAccount.Balance < amount {
			return errors.New("insufficient balance")
		}
		fromAccount.Balance -= amount
	} else {
		return errors.New("can not find the account")
	}

	if toAccount, ok := m.uidMap[to]; ok {
		toAccount.RwMutex.Lock()
		defer toAccount.RwMutex.Unlock()
		toAccount.Balance += amount
	} else {
		return errors.New("can not find the account")
	}
	return nil
}

func (m *MMap) getAllBalance() map[int64]int64 {
	balanceMap := make(map[int64]int64)
	for uid, account := range m.uidMap {
		account.RwMutex.RLock()
		balanceMap[uid] = account.Balance
		account.RwMutex.RUnlock()
	}
	return balanceMap
}
