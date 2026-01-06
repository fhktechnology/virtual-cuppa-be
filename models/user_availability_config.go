package models

import (
	"time"

	"gorm.io/gorm"
)

// UserAvailabilityConfig represents user's system-wide availability configuration
// This configuration is used for matching users with compatible schedules
type UserAvailabilityConfig struct {
	ID                 uint           `gorm:"primarykey" json:"id"`
	UserID             uint           `gorm:"uniqueIndex;not null" json:"userId"`
	User               *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	MondayMorning      bool           `gorm:"default:false" json:"mondayMorning"`
	MondayAfternoon    bool           `gorm:"default:false" json:"mondayAfternoon"`
	TuesdayMorning     bool           `gorm:"default:false" json:"tuesdayMorning"`
	TuesdayAfternoon   bool           `gorm:"default:false" json:"tuesdayAfternoon"`
	WednesdayMorning   bool           `gorm:"default:false" json:"wednesdayMorning"`
	WednesdayAfternoon bool           `gorm:"default:false" json:"wednesdayAfternoon"`
	ThursdayMorning    bool           `gorm:"default:false" json:"thursdayMorning"`
	ThursdayAfternoon  bool           `gorm:"default:false" json:"thursdayAfternoon"`
	FridayMorning      bool           `gorm:"default:false" json:"fridayMorning"`
	FridayAfternoon    bool           `gorm:"default:false" json:"fridayAfternoon"`
	SaturdayMorning    bool           `gorm:"default:false" json:"saturdayMorning"`
	SaturdayAfternoon  bool           `gorm:"default:false" json:"saturdayAfternoon"`
	SundayMorning      bool           `gorm:"default:false" json:"sundayMorning"`
	SundayAfternoon    bool           `gorm:"default:false" json:"sundayAfternoon"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// GetAvailableSlots returns a list of available time slots in a readable format
func (uac *UserAvailabilityConfig) GetAvailableSlots() []string {
	slots := []string{}
	
	if uac.MondayMorning {
		slots = append(slots, "Monday morning")
	}
	if uac.MondayAfternoon {
		slots = append(slots, "Monday afternoon")
	}
	if uac.TuesdayMorning {
		slots = append(slots, "Tuesday morning")
	}
	if uac.TuesdayAfternoon {
		slots = append(slots, "Tuesday afternoon")
	}
	if uac.WednesdayMorning {
		slots = append(slots, "Wednesday morning")
	}
	if uac.WednesdayAfternoon {
		slots = append(slots, "Wednesday afternoon")
	}
	if uac.ThursdayMorning {
		slots = append(slots, "Thursday morning")
	}
	if uac.ThursdayAfternoon {
		slots = append(slots, "Thursday afternoon")
	}
	if uac.FridayMorning {
		slots = append(slots, "Friday morning")
	}
	if uac.FridayAfternoon {
		slots = append(slots, "Friday afternoon")
	}
	if uac.SaturdayMorning {
		slots = append(slots, "Saturday morning")
	}
	if uac.SaturdayAfternoon {
		slots = append(slots, "Saturday afternoon")
	}
	if uac.SundayMorning {
		slots = append(slots, "Sunday morning")
	}
	if uac.SundayAfternoon {
		slots = append(slots, "Sunday afternoon")
	}
	
	return slots
}

// ToAvailability converts UserAvailabilityConfig to Availability map format
func (uac *UserAvailabilityConfig) ToAvailability() Availability {
	availability := make(Availability)
	
	if uac.MondayMorning || uac.MondayAfternoon {
		availability["Monday"] = []string{}
		if uac.MondayMorning {
			availability["Monday"] = append(availability["Monday"], "morning")
		}
		if uac.MondayAfternoon {
			availability["Monday"] = append(availability["Monday"], "afternoon")
		}
	}
	
	if uac.TuesdayMorning || uac.TuesdayAfternoon {
		availability["Tuesday"] = []string{}
		if uac.TuesdayMorning {
			availability["Tuesday"] = append(availability["Tuesday"], "morning")
		}
		if uac.TuesdayAfternoon {
			availability["Tuesday"] = append(availability["Tuesday"], "afternoon")
		}
	}
	
	if uac.WednesdayMorning || uac.WednesdayAfternoon {
		availability["Wednesday"] = []string{}
		if uac.WednesdayMorning {
			availability["Wednesday"] = append(availability["Wednesday"], "morning")
		}
		if uac.WednesdayAfternoon {
			availability["Wednesday"] = append(availability["Wednesday"], "afternoon")
		}
	}
	
	if uac.ThursdayMorning || uac.ThursdayAfternoon {
		availability["Thursday"] = []string{}
		if uac.ThursdayMorning {
			availability["Thursday"] = append(availability["Thursday"], "morning")
		}
		if uac.ThursdayAfternoon {
			availability["Thursday"] = append(availability["Thursday"], "afternoon")
		}
	}
	
	if uac.FridayMorning || uac.FridayAfternoon {
		availability["Friday"] = []string{}
		if uac.FridayMorning {
			availability["Friday"] = append(availability["Friday"], "morning")
		}
		if uac.FridayAfternoon {
			availability["Friday"] = append(availability["Friday"], "afternoon")
		}
	}
	
	if uac.SaturdayMorning || uac.SaturdayAfternoon {
		availability["Saturday"] = []string{}
		if uac.SaturdayMorning {
			availability["Saturday"] = append(availability["Saturday"], "morning")
		}
		if uac.SaturdayAfternoon {
			availability["Saturday"] = append(availability["Saturday"], "afternoon")
		}
	}
	
	if uac.SundayMorning || uac.SundayAfternoon {
		availability["Sunday"] = []string{}
		if uac.SundayMorning {
			availability["Sunday"] = append(availability["Sunday"], "morning")
		}
		if uac.SundayAfternoon {
			availability["Sunday"] = append(availability["Sunday"], "afternoon")
		}
	}
	
	return availability
}

// HasCommonAvailability checks if two users have at least one common available time slot
func HasCommonAvailability(config1, config2 *UserAvailabilityConfig) bool {
	return (config1.MondayMorning && config2.MondayMorning) ||
		(config1.MondayAfternoon && config2.MondayAfternoon) ||
		(config1.TuesdayMorning && config2.TuesdayMorning) ||
		(config1.TuesdayAfternoon && config2.TuesdayAfternoon) ||
		(config1.WednesdayMorning && config2.WednesdayMorning) ||
		(config1.WednesdayAfternoon && config2.WednesdayAfternoon) ||
		(config1.ThursdayMorning && config2.ThursdayMorning) ||
		(config1.ThursdayAfternoon && config2.ThursdayAfternoon) ||
		(config1.FridayMorning && config2.FridayMorning) ||
		(config1.FridayAfternoon && config2.FridayAfternoon) ||
		(config1.SaturdayMorning && config2.SaturdayMorning) ||
		(config1.SaturdayAfternoon && config2.SaturdayAfternoon) ||
		(config1.SundayMorning && config2.SundayMorning) ||
		(config1.SundayAfternoon && config2.SundayAfternoon)
}

// GetCommonSlots returns a list of common available time slots between two users
func GetCommonSlots(config1, config2 *UserAvailabilityConfig) []string {
	slots := []string{}
	
	if config1.MondayMorning && config2.MondayMorning {
		slots = append(slots, "Monday morning")
	}
	if config1.MondayAfternoon && config2.MondayAfternoon {
		slots = append(slots, "Monday afternoon")
	}
	if config1.TuesdayMorning && config2.TuesdayMorning {
		slots = append(slots, "Tuesday morning")
	}
	if config1.TuesdayAfternoon && config2.TuesdayAfternoon {
		slots = append(slots, "Tuesday afternoon")
	}
	if config1.WednesdayMorning && config2.WednesdayMorning {
		slots = append(slots, "Wednesday morning")
	}
	if config1.WednesdayAfternoon && config2.WednesdayAfternoon {
		slots = append(slots, "Wednesday afternoon")
	}
	if config1.ThursdayMorning && config2.ThursdayMorning {
		slots = append(slots, "Thursday morning")
	}
	if config1.ThursdayAfternoon && config2.ThursdayAfternoon {
		slots = append(slots, "Thursday afternoon")
	}
	if config1.FridayMorning && config2.FridayMorning {
		slots = append(slots, "Friday morning")
	}
	if config1.FridayAfternoon && config2.FridayAfternoon {
		slots = append(slots, "Friday afternoon")
	}
	if config1.SaturdayMorning && config2.SaturdayMorning {
		slots = append(slots, "Saturday morning")
	}
	if config1.SaturdayAfternoon && config2.SaturdayAfternoon {
		slots = append(slots, "Saturday afternoon")
	}
	if config1.SundayMorning && config2.SundayMorning {
		slots = append(slots, "Sunday morning")
	}
	if config1.SundayAfternoon && config2.SundayAfternoon {
		slots = append(slots, "Sunday afternoon")
	}
	
	return slots
}

// CreateAvailabilityConfigInput represents input for creating availability configuration
type CreateAvailabilityConfigInput struct {
	MondayMorning      bool `json:"mondayMorning"`
	MondayAfternoon    bool `json:"mondayAfternoon"`
	TuesdayMorning     bool `json:"tuesdayMorning"`
	TuesdayAfternoon   bool `json:"tuesdayAfternoon"`
	WednesdayMorning   bool `json:"wednesdayMorning"`
	WednesdayAfternoon bool `json:"wednesdayAfternoon"`
	ThursdayMorning    bool `json:"thursdayMorning"`
	ThursdayAfternoon  bool `json:"thursdayAfternoon"`
	FridayMorning      bool `json:"fridayMorning"`
	FridayAfternoon    bool `json:"fridayAfternoon"`
	SaturdayMorning    bool `json:"saturdayMorning"`
	SaturdayAfternoon  bool `json:"saturdayAfternoon"`
	SundayMorning      bool `json:"sundayMorning"`
	SundayAfternoon    bool `json:"sundayAfternoon"`
}

// UpdateAvailabilityConfigInput represents input for updating availability configuration
type UpdateAvailabilityConfigInput struct {
	MondayMorning      *bool `json:"mondayMorning,omitempty"`
	MondayAfternoon    *bool `json:"mondayAfternoon,omitempty"`
	TuesdayMorning     *bool `json:"tuesdayMorning,omitempty"`
	TuesdayAfternoon   *bool `json:"tuesdayAfternoon,omitempty"`
	WednesdayMorning   *bool `json:"wednesdayMorning,omitempty"`
	WednesdayAfternoon *bool `json:"wednesdayAfternoon,omitempty"`
	ThursdayMorning    *bool `json:"thursdayMorning,omitempty"`
	ThursdayAfternoon  *bool `json:"thursdayAfternoon,omitempty"`
	FridayMorning      *bool `json:"fridayMorning,omitempty"`
	FridayAfternoon    *bool `json:"fridayAfternoon,omitempty"`
	SaturdayMorning    *bool `json:"saturdayMorning,omitempty"`
	SaturdayAfternoon  *bool `json:"saturdayAfternoon,omitempty"`
	SundayMorning      *bool `json:"sundayMorning,omitempty"`
	SundayAfternoon    *bool `json:"sundayAfternoon,omitempty"`
}
