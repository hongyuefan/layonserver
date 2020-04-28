package models

import (
	"layonserver/util/qmsql"
	"time"
)

type Users struct {
	Id       int64     `gorm:"column:id;"`
	UserName string    `gorm:"column:user_name;"`
	CreateTm time.Time `gorm:"column:created"`
	Modified time.Time `gorm:"column:updated"`
	IsDel    string    `gorm:"column:is_deleted"`
}

func init() {
	qmsql.DEFAULTDB.AutoMigrate(&Users{})
}

func (u Users) TableName() string {
	return "user"
}

func GetUsers(offset int64, limit int64) (users []Users, err error) {
	err = qmsql.DEFAULTDB.Model(new(Users)).Where("is_deleted = ?", "N").Order("id desc").Offset(offset).Limit(limit).Find(&users).Error
	return
}
