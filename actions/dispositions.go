package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	uuid "github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/wung-s/gotv/models"
)

// DispositionsCreate default implementation.
func DispositionsCreate(c buffalo.Context) error {

	disposition := &models.Disposition{}

	if err := c.Bind(disposition); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	q := tx.Where("id = ?", c.Param("id"))

	exist, err := q.Exists("members")
	if err != nil {
		return errors.WithStack(err)
	}

	if !exist {
		return c.Render(500, r.JSON("Member not found"))
	}

	id := uuid.UUID{}
	if id, err = uuid.FromString(c.Param("id")); err != nil {
		return errors.WithStack(err)
	}

	disposition.MemberID = id

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(disposition)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Render errors as JSON
		return c.Render(400, r.JSON(verrs))
	}
	return c.Render(200, r.JSON(disposition))
}
