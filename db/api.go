package db

type engine interface {
	addMoney(uid int64, amount int64)
	getBalance(uid int64) (int64, error)
	transfer(from, to, amount int64) error
}

var myEngine engine

func init() {
	// init global db engine
	myEngine = NewMMapEngine()
}

func AddMoney(uid int64, amount int64) {
	myEngine.addMoney(uid, amount)
}

func GetBalance(uid int64) (int64, error) {
	return myEngine.getBalance(uid)
}

func Transfer(from, to, amount int64) error {
	return myEngine.transfer(from, to, amount)
}
