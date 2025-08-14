package model

type Pixel struct {
	ID     uint `gorm:"primaryKey"`
	X      uint `gorm:"not null"`
	Y      uint `gorm:"not null"`
	Color  string
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
}
