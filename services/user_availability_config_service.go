package services

import (
	"errors"
	"virtual-cuppa-be/models"
	"virtual-cuppa-be/repositories"
)

var (
	ErrConfigNotFound      = errors.New("availability configuration not found")
	ErrConfigAlreadyExists = errors.New("availability configuration already exists for this user")
	ErrNoAvailabilitySet   = errors.New("at least one availability slot must be selected")
)

type UserAvailabilityConfigService interface {
	CreateConfig(userID uint, input models.CreateAvailabilityConfigInput) (*models.UserAvailabilityConfig, error)
	GetConfig(userID uint) (*models.UserAvailabilityConfig, error)
	UpdateConfig(userID uint, input models.UpdateAvailabilityConfigInput) (*models.UserAvailabilityConfig, error)
	DeleteConfig(userID uint) error
	HasConfig(userID uint) (bool, error)
}

type userAvailabilityConfigService struct {
	configRepo repositories.UserAvailabilityConfigRepository
	userRepo   repositories.UserRepository
}

func NewUserAvailabilityConfigService(
	configRepo repositories.UserAvailabilityConfigRepository,
	userRepo repositories.UserRepository,
) UserAvailabilityConfigService {
	return &userAvailabilityConfigService{
		configRepo: configRepo,
		userRepo:   userRepo,
	}
}

func (s *userAvailabilityConfigService) CreateConfig(userID uint, input models.CreateAvailabilityConfigInput) (*models.UserAvailabilityConfig, error) {
	// Check if user exists
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if config already exists
	exists, err := s.configRepo.Exists(userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrConfigAlreadyExists
	}

	// Validate that at least one slot is selected
	if !s.hasAtLeastOneSlot(input) {
		return nil, ErrNoAvailabilitySet
	}

	config := &models.UserAvailabilityConfig{
		UserID:             userID,
		MondayMorning:      input.MondayMorning,
		MondayAfternoon:    input.MondayAfternoon,
		TuesdayMorning:     input.TuesdayMorning,
		TuesdayAfternoon:   input.TuesdayAfternoon,
		WednesdayMorning:   input.WednesdayMorning,
		WednesdayAfternoon: input.WednesdayAfternoon,
		ThursdayMorning:    input.ThursdayMorning,
		ThursdayAfternoon:  input.ThursdayAfternoon,
		FridayMorning:      input.FridayMorning,
		FridayAfternoon:    input.FridayAfternoon,
		SaturdayMorning:    input.SaturdayMorning,
		SaturdayAfternoon:  input.SaturdayAfternoon,
		SundayMorning:      input.SundayMorning,
		SundayAfternoon:    input.SundayAfternoon,
	}

	if err := s.configRepo.Create(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (s *userAvailabilityConfigService) GetConfig(userID uint) (*models.UserAvailabilityConfig, error) {
	config, err := s.configRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, ErrConfigNotFound
	}
	return config, nil
}

func (s *userAvailabilityConfigService) UpdateConfig(userID uint, input models.UpdateAvailabilityConfigInput) (*models.UserAvailabilityConfig, error) {
	config, err := s.configRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, ErrConfigNotFound
	}

	// Update only provided fields
	if input.MondayMorning != nil {
		config.MondayMorning = *input.MondayMorning
	}
	if input.MondayAfternoon != nil {
		config.MondayAfternoon = *input.MondayAfternoon
	}
	if input.TuesdayMorning != nil {
		config.TuesdayMorning = *input.TuesdayMorning
	}
	if input.TuesdayAfternoon != nil {
		config.TuesdayAfternoon = *input.TuesdayAfternoon
	}
	if input.WednesdayMorning != nil {
		config.WednesdayMorning = *input.WednesdayMorning
	}
	if input.WednesdayAfternoon != nil {
		config.WednesdayAfternoon = *input.WednesdayAfternoon
	}
	if input.ThursdayMorning != nil {
		config.ThursdayMorning = *input.ThursdayMorning
	}
	if input.ThursdayAfternoon != nil {
		config.ThursdayAfternoon = *input.ThursdayAfternoon
	}
	if input.FridayMorning != nil {
		config.FridayMorning = *input.FridayMorning
	}
	if input.FridayAfternoon != nil {
		config.FridayAfternoon = *input.FridayAfternoon
	}
	if input.SaturdayMorning != nil {
		config.SaturdayMorning = *input.SaturdayMorning
	}
	if input.SaturdayAfternoon != nil {
		config.SaturdayAfternoon = *input.SaturdayAfternoon
	}
	if input.SundayMorning != nil {
		config.SundayMorning = *input.SundayMorning
	}
	if input.SundayAfternoon != nil {
		config.SundayAfternoon = *input.SundayAfternoon
	}

	// Validate that at least one slot is still selected
	if !s.hasAtLeastOneSlotInConfig(config) {
		return nil, ErrNoAvailabilitySet
	}

	if err := s.configRepo.Update(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (s *userAvailabilityConfigService) DeleteConfig(userID uint) error {
	config, err := s.configRepo.FindByUserID(userID)
	if err != nil {
		return err
	}
	if config == nil {
		return ErrConfigNotFound
	}

	return s.configRepo.Delete(userID)
}

func (s *userAvailabilityConfigService) HasConfig(userID uint) (bool, error) {
	return s.configRepo.Exists(userID)
}

// Helper function to check if at least one slot is selected in input
func (s *userAvailabilityConfigService) hasAtLeastOneSlot(input models.CreateAvailabilityConfigInput) bool {
	return input.MondayMorning || input.MondayAfternoon ||
		input.TuesdayMorning || input.TuesdayAfternoon ||
		input.WednesdayMorning || input.WednesdayAfternoon ||
		input.ThursdayMorning || input.ThursdayAfternoon ||
		input.FridayMorning || input.FridayAfternoon ||
		input.SaturdayMorning || input.SaturdayAfternoon ||
		input.SundayMorning || input.SundayAfternoon
}

// Helper function to check if at least one slot is selected in config
func (s *userAvailabilityConfigService) hasAtLeastOneSlotInConfig(config *models.UserAvailabilityConfig) bool {
	return config.MondayMorning || config.MondayAfternoon ||
		config.TuesdayMorning || config.TuesdayAfternoon ||
		config.WednesdayMorning || config.WednesdayAfternoon ||
		config.ThursdayMorning || config.ThursdayAfternoon ||
		config.FridayMorning || config.FridayAfternoon ||
		config.SaturdayMorning || config.SaturdayAfternoon ||
		config.SundayMorning || config.SundayAfternoon
}
