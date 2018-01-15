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
		addRole("admin")
		addRole("recruiter")
		addRole("manager")
		return nil
	})

})

func addRole(role string) {
	r := &models.Role{
		Name: role,
	}

	models.DB.Create(r)
}
