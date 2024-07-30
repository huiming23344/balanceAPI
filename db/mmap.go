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
	uidMap sync.Map
}

func NewMMapEngine() *MMap {
	mp := sync.Map{}
	return &MMap{
		uidMap: mp,
	}
}

// Add will add balance to uid account
// if account do not exist just add a new account
func (m *MMap) addMoney(uid int64, amount int64) {
	if account, ok := m.uidMap.Load(uid); ok {
		account := account.(*balanceAccount)
		account.RwMutex.Lock()
		account.Balance += amount
		account.RwMutex.Unlock()
	} else {
		account := &balanceAccount{
			Uid:     uid,
			Balance: amount,
			RwMutex: sync.RWMutex{},
		}
		m.uidMap.Store(uid, account)
	}
}

func (m *MMap) getBalance(uid int64) (int64, error) {
	if account, ok := m.uidMap.Load(uid); ok {
		account := account.(*balanceAccount)
		account.RwMutex.RLock()
		defer account.RwMutex.RUnlock()
		return account.Balance, nil

	} else {
		return 0, errors.New("can not find the account")
	}
}

// Transfer will transfer amount from 'from' account to 'to' account
func (m *MMap) transfer(from, to, amount int64) error {
	if fromAccount, ok := m.uidMap.Load(from); ok {
		fromAccount := fromAccount.(*balanceAccount)
		fromAccount.RwMutex.Lock()
		defer fromAccount.RwMutex.Unlock()
		if fromAccount.Balance < amount {
			return errors.New("insufficient balance")
		}
		fromAccount.Balance -= amount
	} else {
		return errors.New("can not find the account")
	}

	if toAccount, ok := m.uidMap.Load(to); ok {
		toAccount := toAccount.(*balanceAccount)
		toAccount.RwMutex.Lock()
		defer toAccount.RwMutex.Unlock()
		toAccount.Balance += amount
	} else {
		return errors.New("can not find the account")
	}
	return nil
}
