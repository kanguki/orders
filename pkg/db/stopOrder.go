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

func (sor StopOrderRepo) Get(so *[]StopOrder) error {
	return sor.Db.Find(so).Error
}

func (sor StopOrderRepo) GetId(so *StopOrder, id string) error {
	return sor.Db.Where("id = ?", id).First(so).Error
}

type StopOrder struct {
	Id             int64  `gorm:"primaryKey,unique"`
	Code           string `gorm:"size:20"`
	Quantity       int32
	SellBuyType    string //BUY, SELL
	StopPrice      float32
	StopVolume     float32
	OrderPrice     float32
	OrderType      string `gorm:"size:100"` //STOP, STOP_LIMIT
	Status         string `gorm:"size:50"`  //PENDING, COMPLETED, CANCELLED, FAILED, SENDING
	FromDate       time.Time
	ToDate         time.Time
	Username       string `gorm:"size:255"`
	AccountNumber  string `gorm:"size:255"`
	OrderNumber    string `gorm:"size:255"`
	FailReason     string
	SecuritiesType string //STOCK, FUND, BOND, ETF, CW, FUTURES
	OrderedAt      time.Time
	CancelledAt    time.Time
	CancelledBy    string `gorm:"size:255"`
	Header         string
	SourceIp       string    `gorm:"size:255"`
	TradingAccSeq  string    `gorm:"size:10"`
	CreatedBy      string    `gorm:"size:255"`
	UpdatedBy      string    `gorm:"size:255"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}
