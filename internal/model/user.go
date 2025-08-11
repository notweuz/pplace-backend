package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username   string `gorm:"unique"`
	Password   string
	LastPlaced time.Time
	Active     bool
}
