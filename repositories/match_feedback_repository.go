package repositories

import (
	"virtual-cuppa-be/models"

	"gorm.io/gorm"
)

type MatchFeedbackRepository interface {
	Create(feedback *models.MatchFeedback) error
	FindByID(id uint) (*models.MatchFeedback, error)
	FindByMatchAndUser(matchID, userID uint) (*models.MatchFeedback, error)
	FindByMatch(matchID uint) ([]*models.MatchFeedback, error)
	FindByUser(userID uint) ([]*models.MatchFeedback, error)
	HasFeedback(matchID, userID uint) (bool, error)
	CountFeedbacksByMatch(matchID uint) (int64, error)
	Update(feedback *models.MatchFeedback) error
	Delete(id uint) error
}

type matchFeedbackRepository struct {
	db *gorm.DB
}

func NewMatchFeedbackRepository(db *gorm.DB) MatchFeedbackRepository {
	return &matchFeedbackRepository{db: db}
}

func (r *matchFeedbackRepository) Create(feedback *models.MatchFeedback) error {
	return r.db.Create(feedback).Error
}

func (r *matchFeedbackRepository) FindByID(id uint) (*models.MatchFeedback, error) {
	var feedback models.MatchFeedback
	err := r.db.Preload("User").Preload("Match").First(&feedback, id).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *matchFeedbackRepository) FindByMatchAndUser(matchID, userID uint) (*models.MatchFeedback, error) {
	var feedback models.MatchFeedback
	err := r.db.Where("match_id = ? AND user_id = ?", matchID, userID).First(&feedback).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *matchFeedbackRepository) FindByMatch(matchID uint) ([]*models.MatchFeedback, error) {
	var feedbacks []*models.MatchFeedback
	err := r.db.Preload("User").Where("match_id = ?", matchID).Find(&feedbacks).Error
	return feedbacks, err
}

func (r *matchFeedbackRepository) FindByUser(userID uint) ([]*models.MatchFeedback, error) {
	var feedbacks []*models.MatchFeedback
	err := r.db.Preload("Match").Where("user_id = ?", userID).Find(&feedbacks).Error
	return feedbacks, err
}

func (r *matchFeedbackRepository) HasFeedback(matchID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.MatchFeedback{}).
		Where("match_id = ? AND user_id = ?", matchID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *matchFeedbackRepository) CountFeedbacksByMatch(matchID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.MatchFeedback{}).
		Where("match_id = ?", matchID).
		Count(&count).Error
	return count, err
}

func (r *matchFeedbackRepository) Update(feedback *models.MatchFeedback) error {
	return r.db.Save(feedback).Error
}

func (r *matchFeedbackRepository) Delete(id uint) error {
	return r.db.Delete(&models.MatchFeedback{}, id).Error
}
