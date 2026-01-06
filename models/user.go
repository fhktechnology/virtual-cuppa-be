package models

import (
	"time"

	"gorm.io/gorm"
)

type AccountType string

const (
	AccountTypeUser  AccountType = "User"
	AccountTypeAdmin AccountType = "Admin"
)

type User struct {
	ID                   uint                     `gorm:"primarykey" json:"id"`
	FirstName            string                   `gorm:"type:varchar(100)" json:"firstName,omitempty"`
	LastName             string                   `gorm:"type:varchar(100)" json:"lastName,omitempty"`
	Email                string                   `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	AccountType          AccountType              `gorm:"type:varchar(10);not null;default:'User'" json:"accountType"`
	OrganisationID       *uint                    `gorm:"index" json:"organisationId,omitempty"`
	Organisation         *Organisation            `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
	Tags                 []Tag                    `gorm:"many2many:user_tags;" json:"tags,omitempty"`
	IsConfirmed          bool                     `gorm:"default:false" json:"isConfirmed"`
	AverageRating        *float64                 `gorm:"type:decimal(3,2)" json:"averageRating,omitempty"`
	AvailabilityConfig   *UserAvailabilityConfig  `gorm:"foreignKey:UserID" json:"availabilityConfig,omitempty"`
	RefreshToken         *string                  `gorm:"type:text" json:"-"`
	CreatedAt            time.Time                `json:"createdAt"`
	UpdatedAt            time.Time                `json:"updatedAt"`
	DeletedAt            gorm.DeletedAt           `gorm:"index" json:"-"`
}

type RegisterInput struct {
	FirstName      string      `json:"firstName,omitempty"`
	LastName       string      `json:"lastName,omitempty"`
	Email          string      `json:"email" binding:"required,email"`
	AccountType    AccountType `json:"accountType" binding:"omitempty,oneof=User Admin"`
	OrganisationID *uint       `json:"organisationId,omitempty"`
}

type RequestCodeInput struct {
	Email string `json:"email" binding:"required,email"`
}

type LoginInput struct {
	Email       string `json:"email" binding:"required,email"`
	ConfirmCode string `json:"confirmCode" binding:"required,len=6"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type CSVUserInput struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

type ConfirmUserInput struct {
	UserID uint `json:"userId" binding:"required"`
}

type UpsertOrganisationInput struct {
	ID         uint   `json:"id"`
	Name       string `json:"name" binding:"required"`
	CompanyUrl string `json:"companyUrl,omitempty"`
}

type CreateUserInput struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

type UpdateTagsInput struct {
	Tags []string `json:"tags" binding:"required"`
}
