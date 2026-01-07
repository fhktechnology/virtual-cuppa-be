package repositories

import (
	"time"
	"virtual-cuppa-be/models"

	"gorm.io/gorm"
)

type MatchHistoryRepository interface {
	Create(history *models.MatchHistory) error
	WasRecentlyMatched(user1ID, user2ID uint, days int) (bool, error)
	WasEverMatched(user1ID, user2ID uint) (bool, error)
}

type matchHistoryRepository struct {
	db *gorm.DB
}

func NewMatchHistoryRepository(db *gorm.DB) MatchHistoryRepository {
	return &matchHistoryRepository{db: db}
}

func (r *matchHistoryRepository) Create(history *models.MatchHistory) error {
	return r.db.Create(history).Error
}

func (r *matchHistoryRepository) WasRecentlyMatched(user1ID, user2ID uint, days int) (bool, error) {
	var count int64
	cutoffDate := time.Now().AddDate(0, 0, -days)
	
	err := r.db.Model(&models.MatchHistory{}).
		Where("((user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)) AND matched_at > ?",
			user1ID, user2ID, user2ID, user1ID, cutoffDate).
		Count(&count).Error
	
	return count > 0, err
}

func (r *matchHistoryRepository) WasEverMatched(user1ID, user2ID uint) (bool, error) {
	var count int64
	
	err := r.db.Model(&models.MatchHistory{}).
		Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
			user1ID, user2ID, user2ID, user1ID).
		Count(&count).Error
	
	return count > 0, err
}
