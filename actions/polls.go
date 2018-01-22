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

// PollsList gets all Polls. This function is mapped to the path
// GET /polls
func PollsList(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	polls := &models.Polls{}

	// Retrieve all Polls from the DB
	if err := tx.All(polls); err != nil {
		return errors.WithStack(err)
	}

	type PollMemberStats struct {
		Total     int `json:"total"`
		Voted     int `json:"voted"`
		Supporter int `json:"supporter"`
	}

	type PollWithMembers struct {
		models.Poll
		MemberStatistic PollMemberStats `json:"memberStatistic"`
	}

	type PollsWithMembers []PollWithMembers

	result := PollsWithMembers{}

	// retrieve members of each pole
	for _, p := range *polls {
		members := &models.Members{}
		s := PollMemberStats{}
		q := tx.BelongsTo(&p)

		// get total members in the poll
		cnt, err := q.Count(members)
		if err != nil {
			return errors.WithStack(err)
		}
		s.Total = cnt

		// get total voted members in the poll
		cnt, err = q.Where("voted = ?", true).Count(members)
		if err != nil {
			return errors.WithStack(err)
		}
		s.Voted = cnt

		if err := q.All(members); err != nil {
			return errors.WithStack(err)
		}

		tmp := PollWithMembers{p, s}

		result = append(result, tmp)
	}

	return c.Render(200, r.JSON(result))
}
