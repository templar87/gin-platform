package services

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitGormDb() {
	db, err := gorm.Open("mysql", "username:password!@tcp(127.0.0.1:3306)/fishballcurrency?charset=utf8&parseTime=True")
	if err != nil {
		log.Panic("Failed to connect to database: %v\n", err)
	}

	db.DB()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	// Disable table name's pluralization
	db.SingularTable(true)
	DB = db
}
