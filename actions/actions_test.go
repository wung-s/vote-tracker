package actions

import (
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/suite"
	"github.com/wung-s/gotv/models"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	action, err := suite.NewActionWithFixtures(App(), packr.NewBox("../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ActionSuite{
		Action: action,
	}
	suite.Run(t, as)
}

func (as ActionSuite) createUser() error {
	u := models.User{Password: "ffffff", UserName: "test1@test.com"}
	return as.DB.Create(&u)
}

func (as ActionSuite) setAuthorization() {
	as.Willie.Headers["Authorization"] = "bearer " + "eyJhbGciOiJSUzI1NiIsImtpZCI6ImMxYTg1OWFmNjkxNTZjODMwMGY2NzllMGMxODJlMGJkMjBmNzA4MDEifQ.eyJpc3MiOiJodHRwczovL3NlY3VyZXRva2VuLmdvb2dsZS5jb20vZ290di1kZXYiLCJhdWQiOiJnb3R2LWRldiIsImF1dGhfdGltZSI6MTUxODYwMDg0MCwidXNlcl9pZCI6IkpOTklYN1BTemlSbndUQVRsMHRUUWg1cnM5RzIiLCJzdWIiOiJKTk5JWDdQU3ppUm53VEFUbDB0VFFoNXJzOUcyIiwiaWF0IjoxNTE4NjAwODQwLCJleHAiOjE1MTg2MDQ0NDAsImVtYWlsIjoidGVzdDFAdGVzdC5jb20iLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsiZW1haWwiOlsidGVzdDFAdGVzdC5jb20iXX0sInNpZ25faW5fcHJvdmlkZXIiOiJwYXNzd29yZCJ9fQ.i2-ct9uDlU0_KfyCESoKCQNgRxgJuaxT0a4pRFGwfk09mY1PiwaGVj9ekVsc-qnmbQKNH6QZCYh6RUn2RDwTkCePipW4_6OCPHaVql7p94Wy306f0mIySoaZaVjZs6D2wXtL4WDB3gk20DinDw9nYusSL_xaa9K7oVPGcRGt3g02d9_lc-IVTlevyAOrVJ95IfqMz67f-n4H0Jf581Yx1e2yTa_SQZ1KrYunn-5Sv4K6SieoT3WYAJmxaE73-3sZWj8p7M_4R7RP0TrZdEtZ9kO7HKmcKeDZeTW_aq2yvEPJR-bMJvbg9sEKYyrxUw4dlWiWd3Vn-lVLUZXwPW1mWA"
}
