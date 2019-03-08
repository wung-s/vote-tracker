package actions

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/wung-s/gotv/models"
	"golang.org/x/crypto/bcrypt"
)

// HomeHandler is a default handler to serve up
// a home page.

type LoginParams struct {
	UserName string `json:"userName" db:"-"`
	Pw       string `json:"pw" db:"-"`
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.JSON(map[string]string{"message": "Welcome to RallyCHQ"}))
}

func LoginHandler(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	params := &UserParams{}
	if err := c.Bind(params); err != nil {
		return errors.WithStack(err)
	}

	user := &models.User{}
	err := tx.Where("user_name = ?", params.UserName).First(user)
	if err != nil {
		return errors.WithStack(err)
	}

	match := checkPasswordHash(params.Pw, user.Password)

	if match {
		claims := jwt.StandardClaims{
			ExpiresAt: time.Now().Add(oneWeek()).Unix(),
			Id:        user.ID.String(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		box := packr.New("../config", "../config")
		signingKey := box.Bytes(os.Getenv("JWT_SIGN_KEY"))

		tokenString, err := token.SignedString(signingKey)
		if err != nil {
			return fmt.Errorf("could not sign token, %v", err)
		}

		return c.Render(200, r.JSON(map[string]string{"token": tokenString}))
	}

	return c.Render(401, r.JSON(map[string]string{"message": "Username/password mismatch"}))

}
