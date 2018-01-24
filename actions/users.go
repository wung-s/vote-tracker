package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/wung-s/gotv/models"
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

	if !exist {
		return c.Render(500, r.JSON("Role not found"))
	}

	role := &models.Role{}
	if err := tx.Find(role, userParams.RoleID); err != nil {
		return errors.WithStack(err)
	}

	if role.Name == "scrutineer" {
		user.PollID = userParams.PollID
		user.PhoneNo = userParams.PhoneNo
	}

	user.Email = userParams.Email

	token := getAuth0Token()

	auth0User, err := createAuth0User(token, user.Email)
	if err != nil {
		fmt.Println(err)
		return errors.WithStack(err)
	}

	user.AuthID = auth0User.Auth0ID

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

func createAuth0User(token string, email string) (auth0Succ, error) {
	url := "https://wung.auth0.com/api/v2/users"
	payload := strings.NewReader(`{"connection":"Username-Password-Authentication","email":"` + email + `","password": "ffffff","email_verified": true}`)
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	type auth0Res struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
	}

	r := auth0Res{}
	if err := json.Unmarshal([]byte(string(body)), &r); err != nil {
		fmt.Println(err)
	}

	if r.StatusCode != 0 {
		return auth0Succ{}, errors.New(r.Message)
	}

	as := auth0Succ{}
	if err := json.Unmarshal([]byte(string(body)), &as); err != nil {
		fmt.Println(err)
	}
	return as, nil
}

func getAuth0Token() string {
	url := "https://wung.auth0.com/oauth/token"

	payload := strings.NewReader("{\"grant_type\":\"client_credentials\",\"client_id\": \"63WJajY4AcAWcj4fJWQEHB3KdKM5co4q\",\"client_secret\": \"Hl5BDf_Xy7PDWpfYLgs1pdg4MGuklH2Efq2Z6TcT03RRhT5TUKRF3iHHib9vNbCQ\",\"audience\": \"https://wung.auth0.com/api/v2/\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	type auth0Res struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
	}

	ar := auth0Res{}
	if err := json.Unmarshal([]byte(string(body)), &ar); err != nil {
		fmt.Println(err)
	}

	return string(ar.AccessToken)
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
	}

	// Delete all existing roles
	sql := "DELETE FROM user_roles as user_roles WHERE user_roles.user_id = ?"
	err := tx.RawQuery(sql, user.ID).Exec()
	if err != nil {
		return errors.WithStack(err)
	}

	userRole := &models.UserRole{}
	userRole.UserID = user.ID
	userRole.RoleID = userParams.RoleID

	verrs, err := tx.ValidateAndCreate(userRole)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Render errors as JSON
		return c.Render(400, r.JSON(verrs))
	}

	return c.Render(200, r.JSON(user))
}
