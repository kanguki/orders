package db

import (
	"time"

	"gorm.io/gorm"
)

type StopOrderRepo struct {
	Db *gorm.DB
}

func GetStopOrderRepo() *StopOrderRepo {
	db := getDb()
	db.AutoMigrate(&StopOrder{})
	return &StopOrderRepo{
		Db: db,
	}
}

func (sor StopOrderRepo) Add(so StopOrder) error {
	return sor.Db.Create(so).Error
}

func (sor StopOrderRepo) Update(so StopOrder) error {
	return sor.Db.Save(so).Error
}

func (sor StopOrderRepo) GetId(so *[]StopOrder) error {
	return sor.Db.Find(so).Error
}

func (sor StopOrderRepo) Get(so *StopOrder, id string) error {
	return sor.Db.Where("id = ?", id).First(so).Error
}

type StopOrder struct {
	gorm.Model
	id             string
	code           string
	quantity       int
	sellBuyType    string //BUY, SELL
	stopPrice      float32
	stopVolume     float32
	orderPrice     float32
	orderType      string //STOP, STOP_LIMIT
	status         string //PENDING, COMPLETED, CANCELLED, FAILED, SENDING
	fromDate       time.Time
	toDate         time.Time
	username       string
	accountNumber  string
	orderNumber    string
	failReason     string
	securitiesType string //STOCK, FUND, BOND, ETF, CW, FUTURES
	orderedAt      time.Time
	cancelledAt    time.Time
	cancelledBy    string
	header         string
	sourceIp       string
	tradingAccSeq  string
	createdAt      time.Time
	updatedAt      time.Time
}
