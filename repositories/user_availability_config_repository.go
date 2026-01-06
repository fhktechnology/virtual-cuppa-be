package repositories

import (
	"errors"

	"virtual-cuppa-be/models"

	"gorm.io/gorm"
)

type UserAvailabilityConfigRepository interface {
	Create(config *models.UserAvailabilityConfig) error
	FindByUserID(userID uint) (*models.UserAvailabilityConfig, error)
	FindByUserIDs(userIDs []uint) ([]*models.UserAvailabilityConfig, error)
	Update(config *models.UserAvailabilityConfig) error
	Delete(userID uint) error
	Exists(userID uint) (bool, error)
}

type userAvailabilityConfigRepository struct {
	db *gorm.DB
}

func NewUserAvailabilityConfigRepository(db *gorm.DB) UserAvailabilityConfigRepository {
	return &userAvailabilityConfigRepository{db: db}
}

func (r *userAvailabilityConfigRepository) Create(config *models.UserAvailabilityConfig) error {
	return r.db.Create(config).Error
}

func (r *userAvailabilityConfigRepository) FindByUserID(userID uint) (*models.UserAvailabilityConfig, error) {
	var config models.UserAvailabilityConfig
	err := r.db.Where("user_id = ?", userID).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

func (r *userAvailabilityConfigRepository) FindByUserIDs(userIDs []uint) ([]*models.UserAvailabilityConfig, error) {
	var configs []*models.UserAvailabilityConfig
	err := r.db.Where("user_id IN ?", userIDs).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *userAvailabilityConfigRepository) Update(config *models.UserAvailabilityConfig) error {
	return r.db.Save(config).Error
}

func (r *userAvailabilityConfigRepository) Delete(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserAvailabilityConfig{}).Error
}

func (r *userAvailabilityConfigRepository) Exists(userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserAvailabilityConfig{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
