package model

import (
	"time"
)

type News struct {
	ID        uint `gorm:"primary_key"`
	Title           string    `gorm:"column:title"`
	Content         string    `gorm:"column:content"`
	Press           string    `gorm:"column:press"`
	Date            time.Time `gorm:"column:date"`
	Url             string    `gorm:"column:url"`
	Portal          string    `gorm:"column:portal"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
