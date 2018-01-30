package grifts

import (
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

	models.DB.Create(u)

	rs := &models.Roles{}
	models.DB.Where("name = ?", role).All(rs)

	ur := &models.UserRole{
		UserID: u.ID,
		RoleID: (*rs)[0].ID,
	}

	models.DB.Create(ur)
}

func addRole(role string) {
	r := &models.Role{
		Name: role,
	}

	models.DB.Create(r)
}
