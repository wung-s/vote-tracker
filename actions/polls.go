package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"github.com/wung-s/gotv/models"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Poll)
// DB Table: Plural (polls)
// Resource: Plural (Polls)
// Path: Plural (/polls)
// View Template Folder: Plural (/templates/polls/)

type PollMemberStats struct {
	Total             int `json:"total"`
	Supporter         int `json:"supporter"`
	VotedSupporter    int `json:"votedSupporter"`
	VotedNonSupporter int `json:"votedNonSupporter"`
}

type PollWithMembers struct {
	models.Poll
	MemberStatistic PollMemberStats `json:"memberStatistic"`
}

type PollsWithMembers []PollWithMembers

// PollsList gets all Polls. This function is mapped to the path
// GET /polls
func PollsList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	polls := models.Polls{}

	// Retrieve all Polls from the DB
	if err := tx.All(&polls); err != nil {
		return errors.WithStack(err)
	}

	result := PollsWithMembers{}

	// retrieve members of each pole
	for _, p := range polls {
		s, err := obtainPollStatistics(tx, &p)
		if err != nil {
			return errors.WithStack(err)
		}

		tmp := PollWithMembers{p, s}
		result = append(result, tmp)
	}

	return c.Render(200, r.JSON(result))
}

// PollsShow gets the data for one User. This function is mapped to
// the path GET /polls/{id}
func PollsShow(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Poll
	poll := &models.Poll{}

	// To find the User the parameter id is used.
	if err := tx.Find(poll, c.Param("id")); err != nil {
		return c.Error(404, err)
	}

	s, err := obtainPollStatistics(tx, poll)
	if err != nil {
		return errors.WithStack(err)
	}

	result := PollWithMembers{*poll, s}

	return c.Render(200, r.JSON(result))
}

func obtainPollStatistics(tx *pop.Connection, p *models.Poll) (PollMemberStats, error) {
	members := &models.Members{}
	pms := PollMemberStats{}

	// get total members in the poll
	q := tx.BelongsTo(p)
	cnt, err := q.Count(members)
	if err != nil {
		return pms, err
	}
	pms.Total = cnt

	// get total supporter members in the poll
	q = tx.BelongsTo(p)
	cnt, err = q.Where("supporter = ?", true).Count(members)
	if err != nil {
		return pms, err
	}
	pms.Supporter = cnt

	// get total voted supporter members in the poll
	q = tx.BelongsTo(p)
	cnt, err = q.Where("voted = ? AND supporter = ?", true, true).Count(members)
	if err != nil {
		return pms, err
	}
	pms.VotedSupporter = cnt

	// get total voted non-supporter members in the poll
	q = tx.BelongsTo(p)
	cnt, err = q.Where("voted = ? AND supporter = ?", true, false).Count(members)
	if err != nil {
		return pms, err
	}
	pms.VotedNonSupporter = cnt
	return pms, nil
}
