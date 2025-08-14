package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username    string `gorm:"unique"`
	Password    []byte
	LastPlaced  time.Time
	Active      bool `gorm:"default:true"`
	Admin       bool `gorm:"not null,default:false"`
	PixelsStock uint `gorm:"default:1"`
}
