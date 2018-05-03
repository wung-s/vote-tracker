package grifts

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/markbates/grift/grift"
	"github.com/wung-s/gotv/actions"
	"github.com/wung-s/gotv/models"
)

var masterUserEmail = os.Getenv("MASTER_USER_EMAIL")
var masterUserPw = os.Getenv("MASTER_USER_PW")
var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		// Add DB seeding stuff here
		ctx := context.Background()
		addRolesAndMasterUser(ctx)

		return nil
	})

	grift.Desc("reset", "Truncates all the tables except for electoral_district and polling_divisions")
	grift.Add("reset", func(c *grift.Context) error {
		// Truncates all data in the database except for electoral_district and polling_divisions table
		// A manager user is created

		// Note: Any new table created should also be listed here
		sql := `TRUNCATE dispositions, ride_requests, members, polls, recruiters,
					  roles, user_roles, users CASCADE;`
		if err := models.DB.RawQuery(sql).Exec(); err != nil {
			log.Fatalf("error truncating tables: %v", err)
		} else {
			log.Println("tables truncated successfully")

			ctx := context.Background()
			addRolesAndMasterUser(ctx)
		}
		return nil
	})
})

func addRolesAndMasterUser(ctx context.Context) {
	addRole("captain")
	addRole("scrutineer")
	addRole("manager")

	addUser(ctx, masterUserEmail, masterUserPw, "manager")
}

func addUser(ctx context.Context, email string, pw string, role string) {
	u := &models.User{
		Email: email,
	}

	u.Password, _ = actions.HashPassword(pw)

	if err := models.DB.Create(u); err != nil {
		fmt.Println("could not add role:", err)
	}

	rs := &models.Role{}
	if err := models.DB.Where("name = ?", role).First(rs); err != nil {
		fmt.Println("role NOT found", err)
	}

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
	} else {
		log.Printf("%v role created successfully", role)
	}
}
