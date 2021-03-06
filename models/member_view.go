package models

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	pgTypes "github.com/mc2soft/pq-types"

	uuid "github.com/gobuffalo/uuid"
)

// MemberView is the same as Member
type MemberView struct {
	ID             uuid.UUID            `json:"id" db:"id"`
	FirstName      string               `json:"firstName" db:"first_name"`
	LastName       string               `json:"lastName" db:"last_name"`
	VoterID        string               `json:"voterId" db:"voter_id"`
	UnitNumber     string               `json:"unitNumber" db:"unit_number"`
	StreetNumber   string               `json:"streetNumber" db:"street_number"`
	StreetName     string               `json:"streetName" db:"street_name"`
	City           string               `json:"city" db:"city"`
	State          string               `json:"state" db:"state"`
	PostalCode     string               `json:"postalCode" db:"postal_code"`
	HomePhone      string               `json:"homePhone" db:"home_phone"`
	CellPhone      string               `json:"cellPhone" db:"cell_phone"`
	Recruiter      string               `json:"recruiter" db:"recruiter"`
	RecruiterID    uuid.UUID            `json:"recruiterId" db:"recruiter_id"`
	PollID         uuid.UUID            `json:"pollId" db:"poll_id"`
	Supporter      bool                 `json:"supporter" db:"supporter"`
	Voted          bool                 `json:"voted" db:"voted"`
	RecruiterPhone string               `json:"recruiterPhone" db:"recruiter_phone"`
	LatLng         pgTypes.PostGISPoint `json:"latlng" db:"latlng"`
}

type MembersView []MemberView

// TableName allows a different table tname to be specified
func (MemberView) TableName() string {
	return "members_view"
}

// FilterFromParam will apply filter from the query parameters
func (msv MembersView) FilterFromParam(q *pop.Query, c buffalo.Context) error {
	if c.Param("q") != "" {
		q = q.Where("(address ILIKE ? OR full_name ILIKE ?)", "%"+c.Param("q")+"%", "%"+c.Param("q")+"%")
	}

	if c.Param("address") != "" {
		q = q.Where("address ILIKE ?", "%"+c.Param("address")+"%")
	}

	if c.Param("street_name") != "" {
		q = q.Where("street_name = ?", c.Param("street_name"))
	}

	if c.Param("street_number") != "" {
		q = q.Where("street_number = ?", c.Param("street_number"))
	}

	if c.Param("poll_id") != "" {
		q = q.Where("poll_id = ?", c.Param("poll_id"))
	}

	if c.Param("voted") != "" {
		q = q.Where("voted = ?", c.Param("voted"))
	}

	if c.Param("voter_id") != "" {
		q = q.Where("voter_id = ?", c.Param("voter_id"))
	}

	if c.Param("unit_number") != "" {
		q = q.Where("unit_number = ?", c.Param("unit_number"))
	}

	if c.Param("recruiter_id") != "" {
		q = q.Where("recruiter_id = ?", c.Param("recruiter_id"))
	}

	if c.Param("supporter") != "" {
		q = q.Where("supporter = ?", c.Param("supporter"))
	}

	return nil

}
