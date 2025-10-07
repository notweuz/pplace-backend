package model

type Pixel struct {
	ID     uint   `gorm:"primaryKey"`
	X      uint   `gorm:"not null"`
	Y      uint   `gorm:"not null"`
	Color  string `gorm:"default:'#FFFFFF'"`
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
}

func NewPixel(id, x, y uint, color string) *Pixel {
	return &Pixel{
		ID:    id,
		X:     x,
		Y:     y,
		Color: color,
	}
}
