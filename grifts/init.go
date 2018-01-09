package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/wung-s/gotv/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
