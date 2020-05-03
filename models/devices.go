package models

import (
	"layonserver/util/qmsql"

	"github.com/jinzhu/gorm"
)

type Devices struct {
	gorm.Model
	PageSt
	UId     uint   `json:"uid" gorm:"column:uid;"`
	DevId   string `json:"device_id" gorm:"column:device_id"`
	Memo    string `json:"memo" gorm:"column:memo"`
	OutDate int64  `json:"out_date" gorm:"column:out_date";`
	Status  int    `json:"status" gorm:"column:status;"`
}

func (Devices) TableName() string {
	return "devices"
}

func (u Devices) InsertDevice() error {
	var one Devices
	qmsql.DEFAULTDB.Model(new(Devices)).Where("uid = ? and device_id = ?", u.UId, u.DevId).Find(&one)
	if one.ID != 0 {
		return nil
	}
	return qmsql.DEFAULTDB.Model(new(Devices)).Create(&u).Error
}

func (u Devices) UpdateDevice(devId string) error {
	return qmsql.DEFAULTDB.Model(new(Devices)).Update("out_date", "status").Error
}

func (u Devices) GetDevices(page, limit int, uid uint) (devs []Devices, total int64, err error) {
	offset := limit * (page - 1)
	db := qmsql.DEFAULTDB.Model(new(Devices))
	if uid != 0 {
		db = db.Where("uid = ?", uid)
	}
	db = db.Count(&total)
	err = db.Order("id desc").Offset(offset).Limit(limit).Find(&devs).Error
	return
}

func (u Devices) Count(uid uint) int64 {
	var count int64
	qmsql.DEFAULTDB.Model(new(Devices)).Where("uid = ?", uid).Count(&count)
	return count
}

func (u Devices) DelDevice() {
	qmsql.DEFAULTDB.Model(new(Devices)).Where("uid = ? and device_id = ?", u.UId, u.DevId).Delete(&Devices{})
}
