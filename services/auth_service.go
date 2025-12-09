package services

import (
	"errors"

	"virtual-cuppa-be/models"
	"virtual-cuppa-be/repositories"
	"virtual-cuppa-be/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type AuthService interface {
	Register(input *models.RegisterInput) (*models.AuthResponse, error)
	Login(input *models.LoginInput) (*models.AuthResponse, error)
	RefreshToken(refreshToken string) (*models.AuthResponse, error)
	GetUserByID(id uint) (*models.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(input *models.RegisterInput) (*models.AuthResponse, error) {
	if input.AccountType == "" {
		input.AccountType = models.AccountTypeUser
	}

	existingUser, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		Password:     string(hashedPassword),
		AccountType:  input.AccountType,
		Organisation: input.Organisation,
		IsConfirmed:  true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.AccountType))
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

	return &models.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *authService) Login(input *models.LoginInput) (*models.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.AccountType))
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

	token, err := utils.GenerateToken(user.ID, user.Email, string(user.AccountType))
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
