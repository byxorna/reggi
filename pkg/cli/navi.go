package cli

import (
	tcell "github.com/gdamore/tcell/v2"
)

var ()

type focus int
type focusdir int

const (
	FocusInput focus = iota
	FocusText
	FocusCaptures

	FocusDirectionLeft focusdir = iota
	FocusDirectionRight
	FocusDirectionUp
	FocusDirectionDown
)

func (c *cli) HandleInputCapture() {
	windowNaviLatch := false
	c.Application.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			// handle navigation between window elements with ctrl-w + directionals
			switch {
			case event.Key() == tcell.KeyCtrlW:
				windowNaviLatch = true
				return nil
			case windowNaviLatch:
				windowNaviLatch = false
				switch event.Rune() {
				case 'h', tcell.RuneLArrow:
					c.SetFocus(FocusDirectionLeft)
					return nil
				case 'j', tcell.RuneDArrow:
					c.SetFocus(FocusDirectionDown)
					return nil
				case 'k', tcell.RuneUArrow:
					c.SetFocus(FocusDirectionUp)
					return nil
				case 'l', tcell.RuneRArrow:
					c.SetFocus(FocusDirectionLeft)
					return nil
				}
			}

			return event
		})
}

func (c *cli) SetFocus(direction focusdir) {
	fv := c.focusedFileView()
	switch c.focus {
	case FocusInput:
		switch direction {
		case FocusDirectionDown, FocusDirectionUp:
			c.Application.SetFocus(fv.textView)
			c.focus = FocusText
		}
	case FocusCaptures:
		switch direction {
		case FocusDirectionDown, FocusDirectionUp:
			c.Application.SetFocus(c.inputView)
			c.focus = FocusInput
		case FocusDirectionLeft, FocusDirectionRight:
			c.Application.SetFocus(fv.textView)
			c.focus = FocusText
		}
	case FocusText:
		switch direction {
		case FocusDirectionDown, FocusDirectionUp:
			c.Application.SetFocus(c.inputView)
			c.focus = FocusInput
		case FocusDirectionLeft, FocusDirectionRight:
			c.Application.SetFocus(fv.fieldView)
			c.focus = FocusCaptures
		}
	}
}

func (c *cli) focusedFileView() *fileView {
	focusedFile, _ := c.pages.GetFrontPage()
	fv := c.fileViews[focusedFile]
	return fv
}
