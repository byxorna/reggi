package cli

import (
	"github.com/rivo/tview"
)

type CLI interface {
	Run() error
}

type cli struct {
	*tview.Application
	box *tview.Box
}

func New() CLI {
	c := cli{}
	c.box = tview.NewBox()
	c.box.SetBorder(true).SetTitle("Hello, world!")
	c.Application = tview.NewApplication()
	c.Application.SetRoot(c.box, true)
	return &c
}
