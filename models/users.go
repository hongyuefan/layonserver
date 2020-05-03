package models

import (
	"layonserver/util/qmsql"

	"github.com/jinzhu/gorm"
)

type PageSt struct {
	Page  int `json:"page" gorm:"-"`
	Limit int `json:"limit" gorm:"-"`
	Total int `json:"total" gorm:"-"`
}

type Users struct {
	gorm.Model
	PageSt
	UserName   string `json:"user_name" gorm:"column:user_name;"`
	PassWord   string `json:"password" gorm:"column:password;"`
	Type       int    `json:"type" gorm:"column:type;commit:'0普通 1一级管理 2 二级管理';"`
	FsId       string `json:"father" gorm:"column:father";`
	Count      int64  `json:"count" gorm:"column:count;"`
	VerifyCode string `json:"verify_code" gorm:"-"`
	VerifyId   string `json:"verify_id" gorm:"-"`
	Token      string `json:"token" gorm:"-"`
}

func (Users) TableName() string {
	return "users"
}

func (u Users) InsertUser() error {
	return qmsql.DEFAULTDB.Model(new(Users)).Create(&u).Error
}

func (u Users) GetUserById(id uint) (user Users, err error) {
	err = qmsql.DEFAULTDB.Model(new(Users)).Where("id = ?", id).Find(&user).Error
	return
}

func (u Users) GetUserByNamePass() (user Users, err error) {
	err = qmsql.DEFAULTDB.Model(new(Users)).Where("user_name = ? and password = ?", u.UserName, u.PassWord).Find(&user).Error
	return
}

func (u Users) GetUsers(page, limit, tp int, fatherId string) (users []Users, err error) {
	offset := limit * (page - 1)
	db := qmsql.DEFAULTDB.Model(new(Users))
	if tp != 0 {
		db = db.Where("type = ?", tp)
	}
	if len(fatherId) > 0 {
		db = db.Where("father = ?", fatherId)
	}
	err = db.Order("id desc").Offset(offset).Limit(limit).Find(&users).Error
	return
}
