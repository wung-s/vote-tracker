package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Recruiter struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
	Name                string    `json:"name" db:"name"`
	PhoneNo             string    `json:"phoneNo" db:"phone_no"`
	Invited             bool      `json:"invited" db:"invited"`
	NotificationEnabled bool      `json:"notificationEnabled" db:"notification_enabled"`
}

// String is not required by pop and may be deleted
func (r Recruiter) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Recruiters is not required by pop and may be deleted
type Recruiters []Recruiter

// String is not required by pop and may be deleted
func (r Recruiters) String() string {
	jr, _ := json.Marshal(r)
	return string(jr)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *Recruiter) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: r.Name, Name: "Name"},
		&validators.StringIsPresent{Field: r.PhoneNo, Name: "PhoneNo"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (r *Recruiter) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (r *Recruiter) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
