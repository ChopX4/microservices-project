package model

type ShipAssembled struct {
	EventUuid    string
	OrderUuid    string
	UserUuid     string
	BuildTimeSec int64
}

type OrderPaid struct {
	EventUuid       string
	OrderUuid       string
	UserUuid        string
	TransactionUuid string
}
