package model

import (
	"time"
)

//ISO 3166-1 alpha-2
const(
	KoreaCode string = "KR"
	USACode string = "US"
	JapanCode string = "JP"
)

type News struct {
	ID        uint `gorm:"primary_key"`
	Title           string    `gorm:"column:title"`
	Content         string    `gorm:"column:content"`
	Press           string    `gorm:"column:press"`
	Date            time.Time `gorm:"column:date"`
	Url             string    `gorm:"column:url"`
	Portal          string    `gorm:"column:portal"`
	CountryCode		string    `gorm:"column:country_code"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
