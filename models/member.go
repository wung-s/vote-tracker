package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	pgTypes "github.com/mc2soft/pq-types"
)

type Member struct {
	ID                uuid.UUID            `json:"id" db:"id"`
	CreatedAt         time.Time            `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time            `json:"updatedAt" db:"updated_at"`
	FirstName         string               `json:"firstName" db:"first_name"`
	LastName          string               `json:"lastName" db:"last_name"`
	VoterID           string               `json:"voterId" db:"voter_id"`
	UnitNumber        string               `json:"unitNumber" db:"unit_number"`
	StreetNumber      string               `json:"streetNumber" db:"street_number"`
	StreetName        string               `json:"streetName" db:"street_name"`
	City              string               `json:"city" db:"city"`
	State             string               `json:"state" db:"state"`
	PostalCode        string               `json:"postalCode" db:"postal_code"`
	HomePhone         string               `json:"homePhone" db:"home_phone"`
	CellPhone         string               `json:"cellPhone" db:"cell_phone"`
	Recruiter         string               `json:"recruiter" db:"recruiter"`
	RecruiterID       uuid.UUID            `json:"recruiterId" db:"recruiter_id"`
	PollID            uuid.UUID            `json:"pollId" db:"poll_id"`
	PollingDivisionID nulls.UUID           `json:"pollingDivisionId" db:"polling_division_id"`
	Supporter         bool                 `json:"supporter" db:"supporter"`
	Voted             bool                 `json:"voted" db:"voted"`
	RecruiterPhone    string               `json:"recruiterPhone" db:"recruiter_phone"`
	LatLng            pgTypes.PostGISPoint `json:"latlng" db:"latlng"`
}

// String is not required by pop and may be deleted
func (m Member) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Address returns the derived full address
func (m Member) Address() string {
	addr := m.UnitNumber + " "
	addr += m.StreetNumber + " "
	addr += m.StreetName + " "
	addr += m.City + " "
	addr += m.State + " "
	addr += m.PostalCode
	return addr
}

// Members is not required by pop and may be deleted
type Members []Member

// String is not required by pop and may be deleted
func (m Members) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// InitializeFromMapString initialized member from the supplied map of string with string type keys
func (m *Member) InitializeFromMapString(arr map[string]string) {
	m.VoterID = arr["voter_id"]
	m.LastName = arr["last_name"]
	m.FirstName = arr["first_name"]
	m.UnitNumber = arr["unit_number"]
	m.StreetNumber = arr["street_number"]
	m.StreetName = arr["street_name"]
	m.City = arr["city"]
	m.State = arr["state"]
	m.PostalCode = arr["postal_code"]
	m.HomePhone = arr["home_phone"]
	m.CellPhone = arr["cell_phone"]
	m.Recruiter = arr["recruiter"]
	m.RecruiterPhone = arr["recruiter_phone"]
	m.Supporter = arr["supporter"] == ("TRUE") || arr["supporter"] == ("true")
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *Member) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: m.LastName, Name: "LastName"},
		&validators.StringIsPresent{Field: m.StreetNumber, Name: "StreetNumber"},
		&validators.StringIsPresent{Field: m.StreetName, Name: "StreetName"},
		&validators.StringIsPresent{Field: m.City, Name: "City"},
		&validators.StringIsPresent{Field: m.State, Name: "State"},
		&validators.StringIsPresent{Field: m.PostalCode, Name: "PostalCode"},
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
