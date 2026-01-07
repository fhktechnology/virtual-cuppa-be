package repositories

import (
	"virtual-cuppa-be/models"

	"gorm.io/gorm"
)

type MatchRepository interface {
	Create(match *models.Match) error
	FindByID(id uint) (*models.Match, error)
	FindCurrentByUserID(userID uint) (*models.Match, error)
	FindByOrganisation(organisationID uint) ([]*models.Match, error)
	FindByUserID(userID uint) ([]*models.Match, error)
	Update(match *models.Match) error
	Delete(id uint) error
	HasPendingMatch(userID uint) (bool, error)
	
	// Availability methods
	CreateAvailability(availability *models.MatchAvailability) error
	UpdateAvailability(availability *models.MatchAvailability) error
	FindAvailabilityByMatchAndUser(matchID, userID uint) (*models.MatchAvailability, error)
	FindAvailabilitiesByMatch(matchID uint) ([]*models.MatchAvailability, error)
}

type matchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) Create(match *models.Match) error {
	return r.db.Create(match).Error
}

func (r *matchRepository) FindByID(id uint) (*models.Match, error) {
	var match models.Match
	err := r.db.Preload("User1").Preload("User2").Preload("Availabilities").Preload("Feedbacks.User").First(&match, id).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *matchRepository) FindCurrentByUserID(userID uint) (*models.Match, error) {
	var match models.Match
	err := r.db.Preload("User1.Tags").Preload("User2.Tags").
		Preload("User1.AvailabilityConfig").Preload("User2.AvailabilityConfig").
		Preload("Availabilities").
		Preload("Feedbacks.User").
		Where(`(user1_id = ? OR user2_id = ?) AND 
		       (status = ? OR status = ?)`, 
			userID, userID, models.MatchStatusPending, models.MatchStatusWaitingForFeedback).
		Order("created_at DESC").
		First(&match).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *matchRepository) FindByOrganisation(organisationID uint) ([]*models.Match, error) {
	var matches []*models.Match
	err := r.db.Preload("User1.Tags").Preload("User2.Tags").
		Preload("User1.AvailabilityConfig").Preload("User2.AvailabilityConfig").
		Preload("Availabilities.User").Preload("Feedbacks.User").
		Where("organisation_id = ?", organisationID).
		Order("created_at DESC").
		Find(&matches).Error
	return matches, err
}

func (r *matchRepository) FindByUserID(userID uint) ([]*models.Match, error) {
	var matches []*models.Match
	err := r.db.Preload("User1.Tags").Preload("User2.Tags").
		Preload("User1.AvailabilityConfig").Preload("User2.AvailabilityConfig").
		Preload("Availabilities.User").Preload("Feedbacks.User").
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Order("created_at DESC").
		Find(&matches).Error
	return matches, err
}

func (r *matchRepository) Update(match *models.Match) error {
	return r.db.Save(match).Error
}

func (r *matchRepository) Delete(id uint) error {
	return r.db.Delete(&models.Match{}, id).Error
}

func (r *matchRepository) HasPendingMatch(userID uint) (bool, error) {
	var count int64
	// User has pending match if:
	// 1. Match status is pending OR waiting_for_feedback
	// 2. User hasn't completed their part yet (either not accepted, or accepted but not given feedback)
	err := r.db.Model(&models.Match{}).
		Where(`((user1_id = ? AND (status = ? OR status = ?)) OR 
		       (user2_id = ? AND (status = ? OR status = ?)))`, 
			userID, models.MatchStatusPending, models.MatchStatusWaitingForFeedback,
			userID, models.MatchStatusPending, models.MatchStatusWaitingForFeedback).
		Count(&count).Error
	return count > 0, err
}

func (r *matchRepository) CreateAvailability(availability *models.MatchAvailability) error {
	return r.db.Create(availability).Error
}

func (r *matchRepository) UpdateAvailability(availability *models.MatchAvailability) error {
	return r.db.Save(availability).Error
}

func (r *matchRepository) FindAvailabilityByMatchAndUser(matchID, userID uint) (*models.MatchAvailability, error) {
	var availability models.MatchAvailability
	err := r.db.Where("match_id = ? AND user_id = ?", matchID, userID).First(&availability).Error
	if err != nil {
		return nil, err
	}
	return &availability, nil
}

func (r *matchRepository) FindAvailabilitiesByMatch(matchID uint) ([]*models.MatchAvailability, error) {
	var availabilities []*models.MatchAvailability
	err := r.db.Preload("User").Where("match_id = ?", matchID).Find(&availabilities).Error
	return availabilities, err
}
