package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type Member struct {
	ID             uuid.UUID `json:"id" db:"id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	FirstName      string    `json:"first_name" db:"first_name"`
	LastName       string    `json:"last_name" db:"last_name"`
	VoterID        string    `json:"voter_id" db:"voter_id"`
	UnitNumber     string    `json:"unit_number" db:"unit_number"`
	StreetNumber   string    `json:"street_number" db:"street_number"`
	StreetName     string    `json:"street_name" db:"street_name"`
	City           string    `json:"city" db:"city"`
	State          string    `json:"state" db:"state"`
	PostalCode     string    `json:"postal_code" db:"postal_code"`
	HomePhone      string    `json:"home_phone" db:"home_phone"`
	CellPhone      string    `json:"cell_phone" db:"cell_phone"`
	Recruiter      string    `json:"recruiter" db:"recruiter"`
	PollID         uuid.UUID `json:"poll_id" db:"poll_id"`
	Supporter bool `json:"supporter" db:"supporter"`
	Voted bool `json:"voted" db:"voted"`
	RecruiterPhone string    `json:"recruiter_phone" db:"recruiter_phone"`
}

// String is not required by pop and may be deleted
func (m Member) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Members is not required by pop and may be deleted
type Members []Member

// String is not required by pop and may be deleted
func (m Members) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *Member) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: m.LastName, Name: "LastName"},
		&validators.StringIsPresent{Field: m.VoterID, Name: "VoterID"},
		&validators.StringIsPresent{Field: m.StreetNumber, Name: "StreetNumber"},
		&validators.StringIsPresent{Field: m.StreetName, Name: "StreetName"},
		&validators.StringIsPresent{Field: m.City, Name: "City"},
		&validators.StringIsPresent{Field: m.State, Name: "State"},
		&validators.StringIsPresent{Field: m.PostalCode, Name: "PostalCode"},
		&validators.StringIsPresent{Field: m.HomePhone, Name: "HomePhone"},
		&validators.StringIsPresent{Field: m.CellPhone, Name: "CellPhone"},
		&validators.StringIsPresent{Field: m.Recruiter, Name: "Recruiter"},
		&validators.StringIsPresent{Field: m.RecruiterPhone, Name: "RecruiterPhone"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *Member) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *Member) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
