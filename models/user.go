package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username   string     `sql:"size:255" form:"username" binding:"required"`
	Password   string     `sql:"size:255" form:"password" binding:"required"`
	Realname   string     `sql:"size:255" form:"realname"`
	Department string     `sql:"size:255" form:"department"`
}
