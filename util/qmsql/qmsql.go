package qmsql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DEFAULTDB *gorm.DB

func InitMysql(conn string) error {
	if db, err := gorm.Open("mysql", conn); err != nil {
		return err
	} else {
		DEFAULTDB = db
		DEFAULTDB.DB().SetMaxIdleConns(10)
		DEFAULTDB.DB().SetMaxOpenConns(10)
		DEFAULTDB.LogMode(true)
	}
	return nil
}
