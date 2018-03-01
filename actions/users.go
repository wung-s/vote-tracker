package actions

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"firebase.google.com/go/auth"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	uuid "github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/wung-s/gotv/models"
	"golang.org/x/net/context"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (User)
// DB Table: Plural (users)
// Resource: Plural (Users)
// Path: Plural (/users)
// View Template Folder: Plural (/templates/users/)

type UserParams struct {
	AuthID  uuid.UUID  `json:"authId" db:"auth_id"`
	RoleID  uuid.UUID  `json:"roleId" db:"role_id"`
	Email   string     `json:"email" db:"email"`
	PollID  nulls.UUID `json:"pollId" db:"poll_id"`
	PhoneNo string     `json:"phoneNo" db:"phone_no"`
	Invited nulls.Bool `json:"invited" db:"invited"`
	Pw      string     `json:"pw" db:"-"`
}

// UsersShow gets the data for one User. This function is mapped to
// the path GET /users/{user_id}
func UsersShow(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}

	// To find the User the parameter user_id is used.
	if err := tx.Find(user, c.Param("user_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(user))
}

// UsersCurrent will return the authenticated user's information
func UsersCurrent(c buffalo.Context) error {
	return c.Render(200, r.JSON(c.Value("currentUser")))
}

type auth0Succ struct {
	Auth0ID string `json:"user_id"`
}

// UsersList gets the data for all User. This function is mapped to
// the path GET /users
func UsersList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	users := &models.Users{}

	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Users from the DB
	if err := q.All(users); err != nil {
		return errors.WithStack(err)
	}

	type UserWithRoles struct {
		models.User
		Roles []string `json:"roles"`
	}

	type UsersWithRoles []UserWithRoles

	uwr := UsersWithRoles{}

	for _, u := range *users {
		roles := &models.Roles{}
		sql := "SELECT roles.* FROM user_roles INNER JOIN roles ON user_roles.role_id = roles.id WHERE user_roles.user_id = ?"
		q := tx.RawQuery(sql, u.ID)
		err := q.All(roles)

		if err != nil {
			return errors.WithStack(err)
		}

		roleNames := []string{}
		// obtain only the name of the roles
		for _, v := range *roles {
			roleNames = append(roleNames, v.Name)
		}

		tmp := UserWithRoles{
			u,
			roleNames,
		}

		uwr = append(uwr, tmp)
	}

	result := struct {
		UsersWithRoles     `json:"users"`
		Page               int `json:"page"`
		PerPage            int `json:"perPage"`
		Offset             int `json:"offset"`
		TotalEntriesSize   int `json:"totalEntriesSize"`
		CurrentEntriesSize int `json:"currentEntriesSize"`
		TotalPages         int `json:"totalPages"`
	}{
		uwr,
		q.Paginator.Page,
		q.Paginator.PerPage,
		q.Paginator.Offset,
		q.Paginator.TotalEntriesSize,
		q.Paginator.CurrentEntriesSize,
		q.Paginator.TotalPages,
	}

	return c.Render(200, r.JSON(result))
}

// UsersCreate adds a User to the DB. This function is mapped to the
// path POST /users
func UsersCreate(c buffalo.Context) error {
	// Allocate an empty User
	user := &models.User{}
	userRole := &models.UserRole{}
	userParams := &UserParams{}

	// Bind user to the html form elements
	if err := c.Bind(userParams); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}
	q := tx.Where("id = ?", userParams.RoleID)

	exist, err := q.Exists("roles")
	if err != nil {
		return errors.WithStack(err)
	}

	if userParams.Pw == "" {
		return c.Render(http.StatusBadRequest, r.JSON("password cannot be blank"))
	}

	if !exist {
		return c.Render(http.StatusBadRequest, r.JSON("role not found"))
	}

	role := &models.Role{}
	if err := tx.Find(role, userParams.RoleID); err != nil {
		return errors.WithStack(err)
	}

	if role.Name == "captain" {
		user.PollID = userParams.PollID
		user.PhoneNo = userParams.PhoneNo
	}

	user.Email = userParams.Email

	// token := getAuth0Token()

	// auth0User, err := createAuth0User(token, user.Email, userParams.Pw)
	fbUserID, err := createFbUser(user.Email, userParams.Pw)
	if err != nil {
		fmt.Println(err)
		return errors.WithStack(err)
	}

	user.AuthID = fbUserID

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(user)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Render errors as JSON
		return c.Render(400, r.JSON(verrs))
	}

	userRole.UserID = user.ID
	userRole.RoleID = userParams.RoleID

	verrs, err = tx.ValidateAndCreate(userRole)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		// return verrs
		return c.Error(400, verrs)
	}

	return c.Render(201, r.JSON(user))
}

func createFbUser(email string, password string) (string, error) {
	client, err := FirebaseApp.Auth(context.Background())
	params := (&auth.UserToCreate{}).
		Email(email).
		EmailVerified(false).
		Password(password).
		Disabled(false)
	u, err := client.CreateUser(context.Background(), params)
	if err != nil {
		fmt.Printf("error creating user: %v\n", err)
	}
	log.Printf("Successfully created user: %v\n", u)
	return u.UID, nil
}

// UsersUpdate changes a User in the DB. This function is mapped to
// the path PUT /users/{user_id}
func UsersUpdate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty User
	user := &models.User{}
	userParams := &UserParams{}

	if err := tx.Find(user, c.Param("id")); err != nil {
		return c.Error(404, err)
	}

	// Bind User to the html form elements
	if err := c.Bind(userParams); err != nil {
		return errors.WithStack(err)
	}

	if !uuid.Equal(userParams.RoleID, uuid.UUID{}) {
		exist, err := tx.Where("id = ?", userParams.RoleID).Exists("roles")
		if err != nil {
			return errors.WithStack(err)
		}

		if !exist {
			return c.Render(500, r.JSON("Role not found"))
		}

		if err := (*user).DeleteAllRoles(tx); err != nil {
			return errors.WithStack(err)
		}

		userRole := &models.UserRole{}
		userRole.UserID = user.ID
		userRole.RoleID = userParams.RoleID

		_, err = tx.ValidateAndCreate(userRole)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	// set PollID if persent
	if !uuid.Equal(userParams.PollID.UUID, uuid.UUID{}) {
		user.PollID = userParams.PollID
	}

	// set Invited if persent
	if v, _ := userParams.Invited.Value(); v != nil {
		user.Invited = userParams.Invited
	}

	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		return c.Render(400, r.JSON(verrs))
	}

	if userParams.Invited.Bool == true {
		err := SendSms(
			user.PhoneNo,
			os.Getenv("TWILIO_NO"),
			"Hello "+user.Email+", you've been invited")

		if err != nil {
			return errors.WithStack(err)
		}
	}

	return c.Render(200, r.JSON(user))
}
