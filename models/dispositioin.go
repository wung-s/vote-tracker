package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
)

type Disposition struct {
	ID          uuid.UUID `json:"id" db:"id"`
	MemberID    uuid.UUID `json:"memberID" db:"member_id"`
	Intention   string    `json:"intention" db:"intention"`
	ContactType string    `json:"contactType" db:"contact_type"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (d Disposition) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Dispositions is not required by pop and may be deleted
type Dispositions []Disposition

// String is not required by pop and may be deleted
func (d Dispositions) String() string {
	jd, _ := json.Marshal(d)
	return string(jd)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (d *Disposition) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (d *Disposition) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (d *Disposition) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
