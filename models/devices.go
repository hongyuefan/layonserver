package models

import (
	"layonserver/util/qmsql"
	"time"
)

type Devices struct {
	Id       int64     `gorm:"column:id;"`
	UId      int64     `gorm:"column:uid;"`
	DevId    string    `gorm:"column:device_id"`
	OutDate  int64     `gorm:"column:out_date";`
	CreateTm time.Time `gorm:"column:created;"`
	Modified time.Time `gorm:"column:updated;"`
	Status   int       `gorm:"column:status;"`
}

func init() {
	qmsql.DEFAULTDB.AutoMigrate(&Devices{})
}

func (u Devices) TableName() string {
	return "devices"
}
