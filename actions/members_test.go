package actions

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gobuffalo/pop"

	"github.com/wung-s/gotv/models"
)

var recruiter = models.Recruiter{
	Name:    "Recruiter 1",
	PhoneNo: "1111111111",
}

func createPoll(db *pop.Connection, p models.Poll) (models.Poll, error) {
	err := db.Create(&p)
	return p, err
}

func createRecruiter(db *pop.Connection, r models.Recruiter) models.Recruiter {
	if err := db.Create(&r); err != nil {
		fmt.Println("create recruiter error:::", err)
	}

	return r
}

func prepareData(db *pop.Connection, recruiter models.Recruiter) {
	poll := models.Poll{
		Name: "Poll1",
	}
	db.Create(&poll)

	for i := 1; i < 24; i++ {
		sr := strconv.Itoa(i)
		member := models.Member{
			FirstName:    "Fname" + sr,
			LastName:     "Lname" + sr,
			VoterID:      sr,
			UnitNumber:   "11",
			StreetNumber: "11",
			StreetName:   "some street",
			City:         "Pune",
			State:        "Maharashtra",
			PostalCode:   "11-111-11",
			HomePhone:    "",
			CellPhone:    "",
			Recruiter:    "",
			RecruiterID:  recruiter.ID,
			PollID:       poll.ID,
			Supporter:    i%2 == 0,
			Voted:        i%2 == 0,
		}
		db.Create(&member)
	}
}

func (as *ActionSuite) Test_MembersSearch_Without_Custom_Pagination() {
	as.createUser()
	as.setAuthorization()
	recruiter := createRecruiter(as.DB, recruiter)
	prepareData(as.DB, recruiter)
	resp := as.JSON("/members/search").Get()
	result := MembersViewSearchResult{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		fmt.Println("error unmarshalling result", err)
	}
	fmt.Println("members::::", result.Members)
	as.Equal(20, len(result.Members))
	as.Equal(20, result.PerPage)
	as.Equal(23, result.TotalEntriesSize)
	as.Equal(20, result.CurrentEntriesSize)
	as.Equal(2, result.TotalPages)
}

func (as *ActionSuite) Test_MembersSearch_With_Custom_Pagination() {
	as.createUser()
	as.setAuthorization()

	recruiter := createRecruiter(as.DB, recruiter)
	prepareData(as.DB, recruiter)
	resp := as.JSON("/members/search?page=1&per_page=4").Get()
	result := MembersViewSearchResult{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		fmt.Println("error unmarshalling result", err)
	}
	as.Equal(4, len(result.Members))
	as.Equal(4, result.PerPage)
	as.Equal(23, result.TotalEntriesSize)
	as.Equal(4, result.CurrentEntriesSize)
	as.Equal(6, result.TotalPages)
}

func (as *ActionSuite) Test_MembersSearch_By_Unknown_Address() {
	as.createUser()
	as.setAuthorization()

	recruiter := createRecruiter(as.DB, recruiter)
	prepareData(as.DB, recruiter)
	resp := as.JSON("/members/search?page=1&per_page=4&address=aa").Get()
	result := MembersViewSearchResult{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		fmt.Println("error unmarshalling result", err)
	}
	as.Equal(0, len(result.Members))
	as.Equal(4, result.PerPage)
	as.Equal(0, result.TotalEntriesSize)
	as.Equal(0, result.CurrentEntriesSize)
	as.Equal(0, result.TotalPages)
}

func (as *ActionSuite) Test_MembersSearch_By_Known_Address() {
	as.createUser()
	as.setAuthorization()

	recruiter := createRecruiter(as.DB, recruiter)
	prepareData(as.DB, recruiter)
	resp := as.JSON("/members/search?page=1&per_page=4&address=11").Get()
	result := MembersViewSearchResult{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		fmt.Println("error unmarshalling result", err)
	}
	as.Equal(4, len(result.Members))
	as.Equal(4, result.PerPage)
	as.Equal(23, result.TotalEntriesSize)
	as.Equal(4, result.CurrentEntriesSize)
	as.Equal(6, result.TotalPages)
}

func (as *ActionSuite) Test_MembersSearch_By_Non_Supporter() {
	as.createUser()
	as.setAuthorization()

	recruiter := createRecruiter(as.DB, recruiter)
	prepareData(as.DB, recruiter)
	resp := as.JSON("/members/search?page=1&per_page=4&supporter=false").Get()
	result := MembersViewSearchResult{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		fmt.Println("error unmarshalling result", err)
	}

	allFalse := true
	for _, v := range result.Members {
		allFalse = v.Supporter == false
	}
	as.Equal(true, allFalse)
	as.Equal(4, len(result.Members))
	as.Equal(4, result.PerPage)
	as.Equal(12, result.TotalEntriesSize)
	as.Equal(4, result.CurrentEntriesSize)
	as.Equal(3, result.TotalPages)
}

func (as *ActionSuite) Test_MembersSearch_By_Voted() {
	as.createUser()
	as.setAuthorization()

	recruiter := createRecruiter(as.DB, recruiter)
	prepareData(as.DB, recruiter)
	resp := as.JSON("/members/search?page=1&per_page=4&voted=true").Get()
	result := MembersViewSearchResult{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		fmt.Println("error unmarshalling result", err)
	}

	allTrue := true
	for _, v := range result.Members {
		allTrue = v.Voted == true
	}
	as.Equal(true, allTrue)
	as.Equal(4, len(result.Members))
	as.Equal(4, result.PerPage)
	as.Equal(11, result.TotalEntriesSize)
	as.Equal(4, result.CurrentEntriesSize)
	as.Equal(3, result.TotalPages)
}

func (as *ActionSuite) Test_MembersSearch_By_Unknown_VoterId() {
	as.createUser()
	as.setAuthorization()

	recruiter := createRecruiter(as.DB, recruiter)
	prepareData(as.DB, recruiter)

	resp := as.JSON("/members/search?page=1&per_page=4&voter_id=5454").Get()
	result := MembersViewSearchResult{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		fmt.Println("error unmarshalling result", err)
	}

	as.Equal(0, len(result.Members))
	as.Equal(4, result.PerPage)
	as.Equal(0, result.TotalEntriesSize)
	as.Equal(0, result.CurrentEntriesSize)
	as.Equal(1, result.TotalPages)
}

func (as *ActionSuite) Test_MembersSearch_By_Known_VoterId() {
	as.createUser()
	as.setAuthorization()

	recruiter := createRecruiter(as.DB, recruiter)
	prepareData(as.DB, recruiter)

	resp := as.JSON("/members/search?page=1&per_page=4&voter_id=2").Get()
	result := MembersViewSearchResult{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		fmt.Println("error unmarshalling result", err)
	}

	as.Equal(1, len(result.Members))
	as.Equal(4, result.PerPage)
	as.Equal(1, result.TotalEntriesSize)
	as.Equal(1, result.CurrentEntriesSize)
	as.Equal(1, result.TotalPages)
}

func (as *ActionSuite) Test_MembersSearch_By_Known_RecruiterId() {
	as.createUser()
	as.setAuthorization()

	r1 := createRecruiter(as.DB, recruiter)
	prepareData(as.DB, r1)
	createRecruiter(as.DB, models.Recruiter{Name: "Recruiter 22", PhoneNo: "22-2233-2211"})

	resp := as.JSON("/members/search?page=1&per_page=4&recruiter_id=" + r1.ID.String()).Get()
	result := MembersViewSearchResult{}
	if err := json.Unmarshal(resp.Body.Bytes(), &result); err != nil {
		fmt.Println("error unmarshalling result", err)
	}
	as.Equal(4, len(result.Members))
	as.Equal(4, result.PerPage)
	as.Equal(23, result.TotalEntriesSize)
	as.Equal(4, result.CurrentEntriesSize)
	as.Equal(6, result.TotalPages)
}
