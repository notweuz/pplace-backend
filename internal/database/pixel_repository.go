package database

import (
	"gorm.io/gorm"
	"pplace_backend/internal/model"
)

type PixelRepository struct {
	db *gorm.DB
}

func NewPixelRepository(db *gorm.DB) PixelRepository {
	return PixelRepository{db: db}
}

func (pr *PixelRepository) Create(pixel *model.Pixel) (*model.Pixel, error) {
	err := pr.db.Create(pixel).Error
	return pixel, err
}

func (pr *PixelRepository) Update(pixel *model.Pixel) (*model.Pixel, error) {
	err := pr.db.Save(pixel).Error
	return pixel, err
}

func (pr *PixelRepository) GetByCoordinates(x, y uint) (*model.Pixel, error) {
	var pixel model.Pixel
	err := pr.db.Where("x = ? AND y = ?", x, y).First(&pixel).Error
	return &pixel, err
}

func (pr *PixelRepository) GetById(id uint) (*model.Pixel, error) {
	var pixel model.Pixel
	err := pr.db.First(&pixel, id).Error
	return &pixel, err
}

func (pr *PixelRepository) GetAll() ([]model.Pixel, error) {
	var pixels []model.Pixel
	return pixels, pr.db.Find(&pixels).Error
}

func (pr *PixelRepository) Delete(id uint) error {
	return pr.db.Delete(&model.Pixel{}, id).Error
}
