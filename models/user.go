package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/nulls"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	// AuthID    string     `json:"authId" db:"auth_id"`
	// Email    string     `json:"email" db:"email"`
	Name     nulls.String `json:"name" db:"name"`
	UserName string       `json:"userName" db:"user_name"`
	PhoneNo  string       `json:"phoneNo" db:"phone_no"`
	Password string       `json:"password" db:"password"`
	PollID   nulls.UUID   `json:"pollId" db:"poll_id"`
	Invited  nulls.Bool   `json:"invited" db:"invited"`
}

// DeleteAllRoles removes all associated roles of a user
func (u User) DeleteAllRoles(tx *pop.Connection) error {
	sql := "DELETE FROM user_roles as user_roles WHERE user_roles.user_id = ?"
	return tx.RawQuery(sql, u.ID).Exec()
}

// Roles gives back the roles of a user
func (u User) Roles(tx *pop.Connection) ([]string, error) {
	roles := &Roles{}
	sql := "SELECT roles.* FROM user_roles INNER JOIN roles ON user_roles.role_id = roles.id WHERE user_roles.user_id = ?"
	if err := tx.RawQuery(sql, u.ID).All(roles); err != nil {
		return []string{}, err
	}

	rnames := []string{}
	for _, r := range *roles {
		rnames = append(rnames, r.Name)
	}
	return rnames, nil
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.UserName, Name: "Username"},
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
