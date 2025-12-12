package services

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"

	"virtual-cuppa-be/models"
	"virtual-cuppa-be/repositories"
)

var (
	ErrAdminNoOrganisation = errors.New("admin must be assigned to an organisation to import users")
	ErrInvalidCSVFormat    = errors.New("invalid CSV format, expected: firstName,lastName,email")
	ErrEmptyCSV            = errors.New("CSV file is empty")
	ErrEmailExists         = errors.New("user with this email already exists")
)

type UserService interface {
	ImportUsersFromCSV(adminID uint, csvContent io.Reader) (int, error)
	ConfirmUser(adminID uint, userID uint) error
	GetUsersByOrganisation(organisationID uint) ([]*models.User, error)
	GetUserByID(userID uint) (*models.User, error)
	UpdateUser(user *models.User) error
	CreateUser(adminID uint, input *models.CreateUserInput) (*models.User, error)
	DeleteUser(adminID uint, userID uint) error
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

	if admin.OrganisationID == nil {
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
			FirstName:      firstName,
			LastName:       lastName,
			Email:          email,
			AccountType:    models.AccountTypeUser,
			OrganisationID: admin.OrganisationID,
			IsConfirmed:    false,
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

	if admin.OrganisationID == nil {
		return ErrAdminNoOrganisation
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if user.OrganisationID == nil || *user.OrganisationID != *admin.OrganisationID {
		return errors.New("user does not belong to admin's organisation")
	}

	user.IsConfirmed = true
	return s.userRepo.Update(user)
}

func (s *userService) GetUsersByOrganisation(organisationID uint) ([]*models.User, error) {
	return s.userRepo.FindByOrganisation(organisationID)
}

func (s *userService) GetUserByID(userID uint) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *userService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *userService) CreateUser(adminID uint, input *models.CreateUserInput) (*models.User, error) {
	admin, err := s.userRepo.FindByID(adminID)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrUserNotFound
	}

	if admin.OrganisationID == nil {
		return nil, ErrAdminNoOrganisation
	}

	// Check if user with this email already exists
	existingUser, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailExists
	}

	user := &models.User{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		AccountType:    models.AccountTypeUser,
		OrganisationID: admin.OrganisationID,
		IsConfirmed:    false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(adminID uint, userID uint) error {
	admin, err := s.userRepo.FindByID(adminID)
	if err != nil {
		return err
	}
	if admin == nil {
		return ErrUserNotFound
	}

	if admin.OrganisationID == nil {
		return ErrAdminNoOrganisation
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Check if user belongs to the same organisation
	if user.OrganisationID == nil || *user.OrganisationID != *admin.OrganisationID {
		return errors.New("user does not belong to your organisation")
	}

	return s.userRepo.Delete(userID)
}
