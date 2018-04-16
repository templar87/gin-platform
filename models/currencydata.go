package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type CurrencyData struct {
	gorm.Model
	Time   time.Time     `sql:"type:date;"`
	High   float64    `sql:"type:decimal(10,6);"`
	Low    float64    `sql:"type:decimal(10,6);"`
	Name   string     `sql:"-"`
}

func (cd CurrencyData) TableName() string {

	return cd.Name
}