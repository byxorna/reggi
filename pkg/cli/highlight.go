package cli

import (
	"fmt"

	tcell "github.com/gdamore/tcell/v2"
)

var (
	matchColor = tcell.ColorSilver
	// TODO: use https://medialab.github.io/iwanthue/ to roll better colors
	captureColors = []tcell.Color{
		tcell.ColorBlue,
		tcell.ColorMaroon,
		tcell.ColorGreen,
		tcell.ColorYellow,
		tcell.ColorNavy,
		tcell.ColorPurple,
		tcell.ColorTeal,
		tcell.ColorRed,
		tcell.ColorLime,
		tcell.ColorFuchsia,
		tcell.ColorAqua,
	}
)

const (
	NoCapture = -1
)

func CaptureColor(hl HighlightID) tcell.Color {
	if hl.Capture == NoCapture {
		return matchColor
	}
	return captureColors[hl.Capture%len(captureColors)]
}

// HighlightID is a region identifier in tview that helps
// identify the line, match, and submatch
type HighlightID struct {
	LineNum  int
	MatchNum int
	Capture  int
}

func (h *HighlightID) String() string {
	if h.Capture != NoCapture {
		return fmt.Sprintf("%d:%d:%d", h.LineNum, h.MatchNum, h.Capture)
	} else {
		return fmt.Sprintf("%d:%d:-", h.LineNum, h.MatchNum)
	}
}

func NewHighlightID(linenum int, matchnum int, submatch int) *HighlightID {
	h := HighlightID{
		LineNum:  linenum,
		MatchNum: matchnum,
		Capture:  NoCapture,
	}
	if submatch > 0 {
		h.Capture = submatch - 1
	}
	return &h
}
