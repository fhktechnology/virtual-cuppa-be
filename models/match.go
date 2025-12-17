package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type MatchStatus string

const (
	MatchStatusPending  MatchStatus = "pending"
	MatchStatusAccepted MatchStatus = "accepted"
	MatchStatusRejected MatchStatus = "rejected"
	MatchStatusExpired  MatchStatus = "expired"
)

// Availability represents user's available time slots
// Format: {"2025-02-18": ["09:30", "10:30"], "2025-02-19": ["14:00"]}
type Availability map[string][]string

// Scan implements sql.Scanner interface
func (a *Availability) Scan(value interface{}) error {
	if value == nil {
		*a = make(Availability)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

// Value implements driver.Valuer interface
func (a Availability) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

type Match struct {
	ID                  uint           `gorm:"primarykey" json:"id"`
	OrganisationID      uint           `gorm:"not null;index" json:"organisationId"`
	User1ID             uint           `gorm:"not null;index" json:"user1Id"`
	User2ID             uint           `gorm:"not null;index" json:"user2Id"`
	User1               *User          `gorm:"foreignKey:User1ID" json:"user1,omitempty"`
	User2               *User          `gorm:"foreignKey:User2ID" json:"user2,omitempty"`
	MatchScore          float64        `gorm:"type:decimal(5,2)" json:"matchScore"`
	Status              MatchStatus    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	User1Accepted       bool           `gorm:"default:false" json:"user1Accepted"`
	User2Accepted       bool           `gorm:"default:false" json:"user2Accepted"`
	User1AcceptedAt     *time.Time     `json:"user1AcceptedAt,omitempty"`
	User2AcceptedAt     *time.Time     `json:"user2AcceptedAt,omitempty"`
	ExpiresAt           *time.Time     `json:"expiresAt,omitempty"`
	Availabilities      []*MatchAvailability `gorm:"foreignKey:MatchID" json:"availabilities,omitempty"`
	Feedbacks           []*MatchFeedback `gorm:"foreignKey:MatchID" json:"feedbacks,omitempty"`
	ScheduledDate       time.Time      `json:"scheduledDate"`
	ScheduledTime       string         `gorm:"type:varchar(10)" json:"scheduledTime"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

type MatchHistory struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	User1ID   uint      `gorm:"not null;index:idx_match_history" json:"user1Id"`
	User2ID   uint      `gorm:"not null;index:idx_match_history" json:"user2Id"`
	MatchedAt time.Time `gorm:"not null" json:"matchedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type MatchAvailability struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	MatchID      uint           `gorm:"not null;index;uniqueIndex:idx_match_user" json:"matchId"`
	UserID       uint           `gorm:"not null;index;uniqueIndex:idx_match_user" json:"userId"`
	Availability Availability   `gorm:"type:jsonb;not null" json:"availability"`
	Match        *Match         `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	User         *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type MatchFeedback struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	MatchID   uint           `gorm:"not null;index;uniqueIndex:idx_feedback_match_user" json:"matchId"`
	UserID    uint           `gorm:"not null;index;uniqueIndex:idx_feedback_match_user" json:"userId"`
	Rating    int            `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment   string         `gorm:"type:text" json:"comment"`
	Match     *Match         `gorm:"foreignKey:MatchID" json:"match,omitempty"`
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
