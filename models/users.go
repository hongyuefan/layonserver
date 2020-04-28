package models

import (
	"layonserver/util/qmsql"
	"time"

	"github.com/jinzhu/gorm"
)

type Users struct {
	gorm.Model
	UserName string `gorm:"column:user_name;"`
	PassWord string `gorm:"column:password;"`
	Type     int    `gorm:"column:type;"`
	FsId     string `gorm:"column:father";`
}

func init() {
	qmsql.DEFAULTDB.AutoMigrate(&Users{})
}

func (u Users) TableName() string {
	return "users"
}

func (u Users) InsertUser() error {
	return qmsql.DEFAULTDB.Model(new(Users)).Create(&u).Error
}

func (u Users) GetUserByNamePass() (user Users, err error) {
	err = qmsql.DEFAULTDB.Model(new(Users)).Where("user_name = ? and password = ?", u.UserName, u.PassWord).Find(&user).Error
	return
}

func (u Users) GetUsers(page, limit int, tp int) (users []Users, err error) {
	offset := limit * (page - 1)
	db := qmsql.DEFAULTDB.Model(new(Users))
	if tp != 0 {
		db = db.Where("type = ?", tp)
	}
	err = db.Order("id desc").Offset(offset).Limit(limit).Find(&users).Error
	return
}
