package grifts

import (
	"fmt"

	"github.com/markbates/grift/grift"
	"github.com/wung-s/gotv/models"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		// Add DB seeding stuff here
		addRole("captain")
		addRole("scrutineer")
		addRole("manager")

		// IMPORTANT: Obtain the uuid from Firebase after manually creating the user
		addUser("test1@test.com", "eXe1tqzn26P7WICMt73Ozipzkw93", "manager")
		return nil
	})
})

func addUser(email string, auth0ID string, role string) {
	u := &models.User{
		Email:  email,
		AuthID: auth0ID,
	}

	if err := models.DB.Create(u); err != nil {
		fmt.Println("could not add role:", err)
	}

	rs := &models.Role{}
	if err := models.DB.Where("name = ?", role).First(rs); err != nil {
		fmt.Println("role NOT found", err)
	}
	fmt.Println("role is: ", *rs)
	ur := &models.UserRole{
		UserID: u.ID,
		RoleID: (*rs).ID,
	}

	if err := models.DB.Create(ur); err != nil {
		fmt.Println("could not add user role:", err)
	}
}

func addRole(role string) {
	r := &models.Role{
		Name: role,
	}

	if err := models.DB.Create(r); err != nil {
		fmt.Println("could not create role", err)
	}
}
