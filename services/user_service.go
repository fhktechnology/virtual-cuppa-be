package services

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"

	"virtual-cuppa-be/models"
	"virtual-cuppa-be/repositories"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAdminNoOrganisation = errors.New("admin must be assigned to an organisation to import users")
	ErrInvalidCSVFormat    = errors.New("invalid CSV format, expected: firstName,lastName,email")
	ErrEmptyCSV            = errors.New("CSV file is empty")
)

type UserService interface {
	ImportUsersFromCSV(adminID uint, csvContent io.Reader) (int, error)
	ConfirmUser(adminID uint, userID uint) error
	GetUsersByOrganisation(organisation string) ([]*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) ImportUsersFromCSV(adminID uint, csvContent io.Reader) (int, error) {
	admin, err := s.userRepo.FindByID(adminID)
	if err != nil {
		return 0, err
	}
	if admin == nil {
		return 0, ErrUserNotFound
	}

	if admin.Organisation == nil || *admin.Organisation == "" {
		return 0, ErrAdminNoOrganisation
	}

	reader := csv.NewReader(csvContent)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, err
	}

	if len(records) == 0 {
		return 0, ErrEmptyCSV
	}

	var users []*models.User
	defaultPassword := "ChangeMe123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	startIndex := 0
	if len(records) > 0 {
		firstRow := records[0]
		if len(firstRow) == 3 && 
		   strings.ToLower(strings.TrimSpace(firstRow[0])) == "firstname" &&
		   strings.ToLower(strings.TrimSpace(firstRow[1])) == "lastname" &&
		   strings.ToLower(strings.TrimSpace(firstRow[2])) == "email" {
			startIndex = 1
		}
	}

	for i := startIndex; i < len(records); i++ {
		record := records[i]
		if len(record) != 3 {
			continue
		}

		firstName := strings.TrimSpace(record[0])
		lastName := strings.TrimSpace(record[1])
		email := strings.TrimSpace(record[2])

		if firstName == "" || lastName == "" || email == "" {
			continue
		}

		existingUser, _ := s.userRepo.FindByEmail(email)
		if existingUser != nil {
			continue
		}

		user := &models.User{
			FirstName:    firstName,
			LastName:     lastName,
			Email:        email,
			Password:     string(hashedPassword),
			AccountType:  models.AccountTypeUser,
			Organisation: admin.Organisation,
			IsConfirmed:  false,
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		return 0, nil
	}

	err = s.userRepo.CreateBatch(users)
	if err != nil {
		return 0, err
	}

	return len(users), nil
}

func (s *userService) ConfirmUser(adminID uint, userID uint) error {
	admin, err := s.userRepo.FindByID(adminID)
	if err != nil {
		return err
	}
	if admin == nil {
		return ErrUserNotFound
	}

	if admin.Organisation == nil || *admin.Organisation == "" {
		return ErrAdminNoOrganisation
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if user.Organisation == nil || *user.Organisation != *admin.Organisation {
		return errors.New("user does not belong to admin's organisation")
	}

	user.IsConfirmed = true
	return s.userRepo.Update(user)
}

func (s *userService) GetUsersByOrganisation(organisation string) ([]*models.User, error) {
	return s.userRepo.FindByOrganisation(organisation)
}
