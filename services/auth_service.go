package services

import (
	"errors"
	"time"

	"virtual-cuppa-be/models"
	"virtual-cuppa-be/repositories"
	"virtual-cuppa-be/utils"
)

var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type AuthService interface {
	Register(input *models.RegisterInput) error
	RequestCode(input *models.RequestCodeInput) error
	Login(input *models.LoginInput) (*models.AuthResponse, error)
	RefreshToken(refreshToken string) (*models.AuthResponse, error)
	GetUserByID(id uint) (*models.User, error)
}

type authService struct {
	userRepo     repositories.UserRepository
	emailService EmailService
	matchService MatchService
}

func NewAuthService(userRepo repositories.UserRepository, emailService EmailService, matchService MatchService) AuthService {
	return &authService{
		userRepo:     userRepo,
		emailService: emailService,
		matchService: matchService,
	}
}

func (s *authService) Register(input *models.RegisterInput) error {
	input.AccountType = models.AccountTypeAdmin

	existingUser, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	confirmCode := utils.GenerateConfirmCode()

	user := &models.User{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		AccountType:    input.AccountType,
		OrganisationID: input.OrganisationID,
		IsConfirmed:    true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	cache := utils.GetConfirmCodeCache()
	cache.Set(input.Email, confirmCode, 5*time.Minute)

	fullName := input.Email
	if input.FirstName != "" && input.LastName != "" {
		fullName = input.FirstName + " " + input.LastName
	} else if input.FirstName != "" {
		fullName = input.FirstName
	}
	
	if err := s.emailService.SendConfirmCode(input.Email, fullName, confirmCode); err != nil {
		return err
	}

	return nil
}

func (s *authService) RequestCode(input *models.RequestCodeInput) error {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	confirmCode := utils.GenerateConfirmCode()

	cache := utils.GetConfirmCodeCache()
	cache.Set(input.Email, confirmCode, 5*time.Minute)

	fullName := user.Email
	if user.FirstName != "" && user.LastName != "" {
		fullName = user.FirstName + " " + user.LastName
	} else if user.FirstName != "" {
		fullName = user.FirstName
	}
	
	if err := s.emailService.SendConfirmCode(input.Email, fullName, confirmCode); err != nil {
		return err
	}

	return nil
}

func (s *authService) Login(input *models.LoginInput) (*models.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	cache := utils.GetConfirmCodeCache()
	cachedCode, exists := cache.Get(input.Email)
	if !exists || cachedCode != input.ConfirmCode {
		return nil, ErrInvalidCredentials
	}

	cache.Delete(input.Email)

	// If user is not confirmed yet (first login), confirm them
	wasUnconfirmed := !user.IsConfirmed
	if wasUnconfirmed {
		user.IsConfirmed = true
	}

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.AccountType), user.OrganisationID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	user.RefreshToken = &refreshToken
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// If user just became confirmed, try to generate a match for them
	if wasUnconfirmed && s.matchService != nil {
		go s.matchService.TryGenerateMatchForUser(user.ID)
	}

	return &models.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *authService) RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	user, err := s.userRepo.FindByRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	if user == nil || user.RefreshToken == nil || *user.RefreshToken != refreshToken {
		return nil, ErrInvalidRefreshToken
	}

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.AccountType), user.OrganisationID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	user.RefreshToken = &newRefreshToken
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token:        token,
		RefreshToken: newRefreshToken,
		User:         *user,
	}, nil
}

func (s *authService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
