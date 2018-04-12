package grifts

import (
	"context"
	"fmt"
	"log"
	"os"

	"firebase.google.com/go/auth"
	"github.com/markbates/grift/grift"
	"github.com/wung-s/gotv/actions"
	"github.com/wung-s/gotv/models"
	"google.golang.org/api/iterator"
)

var masterUserEmail = os.Getenv("MASTER_USER_EMAIL")
var masterUserPw = os.Getenv("MASTER_USER_PW")
var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		// Add DB seeding stuff here
		ctx := context.Background()
		client, err := actions.FirebaseApp.Auth(ctx)
		if err != nil {
			log.Fatalf("error authenticating Firebase: %v\n", err)
		} else {
			addRolesAndMasterUser(ctx, client)
		}

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
			client, err := actions.FirebaseApp.Auth(ctx)
			if err != nil {
				log.Fatalf("error authenticating Firebase: %v\n", err)
			} else {
				deleteFirebaseUsers(ctx, client)
				addRolesAndMasterUser(ctx, client)
			}
		}
		return nil
	})
})

func deleteFirebaseUsers(ctx context.Context, client *auth.Client) {
	// Behind the scenes, the Users() iterator will retrive 1000 Users at a time through the API
	iter := client.Users(ctx, "")
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Fatalf("error listing users: %s\n", err)
		} else {
			err = client.DeleteUser(ctx, user.UID)
			if err != nil {
				log.Fatalf("error deleting user: %v\n", err)
			}
			log.Printf("Successfully deleted user: %s\n", user.UID)
		}
	}
}

func addRolesAndMasterUser(ctx context.Context, client *auth.Client) {
	addRole("captain")
	addRole("scrutineer")
	addRole("manager")

	addUser(ctx, client, masterUserEmail, masterUserPw, "manager")
}

func addUser(ctx context.Context, client *auth.Client, email string, pw string, role string) {
	uuid, _ := actions.CreateFbUser(email, pw)

	u := &models.User{
		Email:  email,
		AuthID: uuid,
	}

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
