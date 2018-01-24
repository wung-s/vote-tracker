package actions

import (
	"bufio"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sfreiberg/gotwilio"
	"github.com/wung-s/gotv/models"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Member)
// DB Table: Plural (members)
// Resource: Plural (Members)
// Path: Plural (/members)
// View Template Folder: Plural (/templates/members/)

// MembersResource is the resource for the Member model
type MembersResource struct {
	buffalo.Resource
}

// MembersList gets all Members. This function is mapped to the path
// GET /members
func MembersList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	members := &models.Members{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Members from the DB
	if err := q.All(members); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the headers so clients know how to paginate.
	c.Response().Header().Set("X-Pagination", q.Paginator.String())

	return c.Render(200, r.JSON(members))
}

// Show gets the data for one Member. This function is mapped to
// the path GET /members/{member_id}
func (v MembersResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Member
	member := &models.Member{}

	// To find the Member the parameter member_id is used.
	if err := tx.Find(member, c.Param("member_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.JSON(member))
}

// Create adds a Member to the DB. This function is mapped to the
// path POST /members
func (v MembersResource) Create(c buffalo.Context) error {
	// Allocate an empty Member
	member := &models.Member{}

	// Bind member to the html form elements
	if err := c.Bind(member); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(member)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Render errors as JSON
		return c.Render(400, r.JSON(verrs))
	}

	return c.Render(201, r.JSON(member))
}

// MembersUpload seeds a Members to the DB. This function is mapped to the
// path POST /members/upload
func MembersUpload(c buffalo.Context) error {
	// Allocate an empty Member

	type UploadParams struct {
		File string `db:"-"`
	}

	postParams := &UploadParams{}

	// Bind member to the html form elements
	if err := c.Bind(postParams); err != nil {
		return errors.WithStack(err)
	}

	if postParams.File == "" {
		return errors.Errorf("No file found")
	}

	// Decode the Base64 string
	dec, err := base64.StdEncoding.DecodeString(strings.Replace(postParams.File, `data:text/csv;base64,`, "", 1))
	if err != nil {
		panic(err)
	}
	fileName := uuid.Must(uuid.NewV4()).String()
	fileName += ".csv"
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}

	// Read the uploaded file content and insert into DB
	csvFile, _ := os.Open(fileName)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	i := 0
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		atEnd := true
		for _, v := range line {
			if strings.TrimSpace(v) != "" {
				atEnd = false
			}
		}

		// Exit the loop if the values in a row are all blank
		if atEnd {
			break
		}

		member := &models.Member{
			VoterID:        line[0],
			LastName:       line[1],
			FirstName:      line[2],
			UnitNumber:     line[3],
			StreetNumber:   line[4],
			StreetName:     line[5],
			City:           line[6],
			State:          line[7],
			PostalCode:     line[8],
			HomePhone:      line[9],
			CellPhone:      line[10],
			Recruiter:      line[11],
			RecruiterPhone: line[13],
			Supporter:      strings.TrimSpace(line[14]) == ("TRUE") || strings.TrimSpace(line[14]) == ("true"),
		}

		pollName := strings.TrimSpace(line[12])
		exist, err := tx.Where("name = ?", pollName).Exists(&models.Poll{})
		if err != nil {
			return errors.WithStack(err)
		}
		poll := &models.Poll{}
		if !exist {
			poll.Name = pollName
			if i != 0 {
				insertPoll(poll, tx)
				setPollID(pollName, member, tx)
			}

		} else {
			setPollID(pollName, member, tx)
		}

		rPhone := strings.TrimSpace(line[13])
		rName := strings.TrimSpace(line[11])
		exist, err = tx.Where("phone_no = ?", rPhone).Exists("recruiters")
		if err != nil {
			return errors.WithStack(err)
		}
		r := &models.Recruiter{}
		if !exist {
			r.PhoneNo = rPhone
			r.Name = rName
			if i != 0 {
				insertRecruiter(r, tx)
				setRecruiterID(rPhone, member, tx)
			}

		} else {
			setRecruiterID(rPhone, member, tx)
		}

		if i != 0 {
			insertMember(member, tx)
		}
		i++
	}

	os.Remove(fileName)

	return c.Render(201, r.JSON("data processing complete"))
}

func setPollID(pollName string, member *models.Member, tx *pop.Connection) {
	polls := []models.Poll{}
	err := tx.Where("name = ?", pollName).All(&polls)
	if err != nil {
		fmt.Print(err)
	} else {
		member.PollID = polls[0].ID
	}
}

func setRecruiterID(p string, member *models.Member, tx *pop.Connection) {
	rs := []models.Recruiter{}
	err := tx.Where("phone_no = ?", p).All(&rs)
	if err != nil {
		fmt.Print(err)
	} else {
		member.RecruiterID = rs[0].ID
	}
}

func insertPoll(poll *models.Poll, tx *pop.Connection) {
	verrs, err := tx.ValidateAndSave(poll)
	poll.ID = uuid.UUID{}
	if err != nil {
		fmt.Print(verrs)
	}
}

func insertRecruiter(r *models.Recruiter, tx *pop.Connection) {
	verrs, err := tx.ValidateAndSave(r)
	r.ID = uuid.UUID{}
	if err != nil {
		fmt.Print(verrs)
	}
}

// insertMember creates new member row in the DB
func insertMember(member *models.Member, tx *pop.Connection) {
	verrs, err := tx.ValidateAndSave(member)
	member.ID = uuid.UUID{}

	if err != nil {
		fmt.Print(verrs)
	}
}

// Edit default implementation. Returns a 404
func (v MembersResource) Edit(c buffalo.Context) error {
	return c.Error(404, errors.New("not available"))
}

// MembersUpdate changes a Member in the DB. This function is mapped to
// the path PUT /members/{member_id}
func MembersUpdate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Member
	member := &models.Member{}
	if err := tx.Find(member, c.Param("id")); err != nil {
		return c.Error(404, err)
	}

	preUpdateVoted := member.Voted

	// Bind Member to the html form elements
	if err := c.Bind(member); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(member)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Render errors as JSON
		return c.Render(400, r.JSON(verrs))
	}

	postUpdateVoted := member.Voted
	if preUpdateVoted == false && postUpdateVoted == true && member.RecruiterPhone != "" {
		SendSms(
			"+1"+member.RecruiterPhone,
			os.Getenv("TWILIO_NO"),
			member.FirstName+" "+member.LastName+" just voted",
		)
	}

	return c.Render(200, r.JSON(member))
}

// SendSms sends out sms
func SendSms(to string, from string, message string) error {
	twilio := gotwilio.NewTwilioClient(os.Getenv("TWILIO_AC_SID"), os.Getenv("TWILIO_AUTH_TOKEN"))
	if _, _, err := twilio.SendSMS(from, to, message, "", ""); err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}

// Destroy deletes a Member from the DB. This function is mapped
// to the path DELETE /members/{member_id}
func (v MembersResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Member
	member := &models.Member{}

	// To find the Member the parameter member_id is used.
	if err := tx.Find(member, c.Param("member_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(member); err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(member))
}

// MembersSearch performs search applying filters from the values in the query parameters
func MembersSearch(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Member
	// member := &models.Member{}
	members := &models.Members{}

	// if c.Param("address") != "" {
	// 	sql := "SELECT * FROM (SELECT *, concat(unit_number::text, street_number::text, street_name::text) AS address FROM public.members) AS temp where temp.address LIKE '%$1%'"
	// 	// sql = sql + ";"
	// 	// args := []string{c.Param("address")}
	// 	fmt.Println("Execute RAW query >>>>>>>>>>>>>>>>>")

	// 	if err := tx.RawQuery(sql, "Jack").All(members); err != nil {
	// 		fmt.Println("Error in query:", err)
	// 	}
	// 	fmt.Println("query result is:", members)
	// }

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	if c.Param("poll_id") != "" {
		q = q.Where("poll_id = ?", c.Param("poll_id"))
	}

	if c.Param("voted") != "" {
		q = q.Where("voted = ?", c.Param("voted"))
	}

	if c.Param("voter_id") != "" {
		q = q.Where("voter_id = ?", c.Param("voter_id"))
	}

	if c.Param("recruiter_id") != "" {
		q = q.Where("recruiter_id = ?", c.Param("recruiter_id"))
	}

	if err := q.All(members); err != nil {
		return c.Error(404, err)
	}

	result := struct {
		models.Members     `json:"members"`
		Page               int `json:"page"`
		PerPage            int `json:"perPage"`
		Offset             int `json:"offset"`
		TotalEntriesSize   int `json:"totalEntriesSize"`
		CurrentEntriesSize int `json:"currentEntriesSize"`
		TotalPages         int `json:"totalPages"`
	}{
		*members,
		q.Paginator.Page,
		q.Paginator.PerPage,
		q.Paginator.Offset,
		q.Paginator.TotalEntriesSize,
		q.Paginator.CurrentEntriesSize,
		q.Paginator.TotalPages,
	}

	return c.Render(200, r.JSON(result))
}

// RecruitersMembersSearch performs search applying filters from the values in the query parameters
func RecruitersMembersSearch(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	members := &models.Members{}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params()).Where("recruiter_id = ?", c.Param("id"))

	if c.Param("voted") != "" {
		q = q.Where("voted = ?", c.Param("voted"))
	}

	if c.Param("voter_id") != "" {
		q = q.Where("voter_id = ?", c.Param("voter_id"))
	}

	if err := q.All(members); err != nil {
		return c.Error(404, err)
	}

	result := struct {
		models.Members     `json:"members"`
		Page               int `json:"page"`
		PerPage            int `json:"perPage"`
		Offset             int `json:"offset"`
		TotalEntriesSize   int `json:"totalEntriesSize"`
		CurrentEntriesSize int `json:"currentEntriesSize"`
		TotalPages         int `json:"totalPages"`
	}{
		*members,
		q.Paginator.Page,
		q.Paginator.PerPage,
		q.Paginator.Offset,
		q.Paginator.TotalEntriesSize,
		q.Paginator.CurrentEntriesSize,
		q.Paginator.TotalPages,
	}

	return c.Render(200, r.JSON(result))
}
