package actions

import (
	"encoding/json"
	"fmt"

	"github.com/markbates/pop"

	"github.com/wung-s/gotv/models"
)

func seedRoles(db *pop.Connection) {
	roles := models.Roles{
		{Name: "captain"},
		{Name: "scrutineer"},
		{Name: "manager"},
	}

	for _, role := range roles {
		db.Create(&role)
	}
}

func (as *ActionSuite) Test_RolesList() {
	seedRoles(as.DB)
	user := models.User{AuthID: "JNNIX7PSziRnwTATl0tTQh5rs9G2", Email: "test1@test.com"}
	as.DB.Create(&user)

	as.setAuthorization()
	dbRoles := models.Roles{}
	as.DB.All(&dbRoles)
	dbNames := dbRoles.Names()

	resp := as.JSON("/roles").Get()
	result := models.Roles{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		fmt.Println("error unmarshalling result", err)
	}

	names := result.Names()
	as.Equal(200, resp.Code)
	as.Equal(dbNames, names)
}
