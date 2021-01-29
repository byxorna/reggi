package cli

import (
	//tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// houses the view for a single file's contents and matches
type fileView struct {
	*tview.Flex
	textView  *tview.TextView
	fieldView *tview.TextView
	rawText   string
	fileName  string
}

func NewFileView(fName string, content string) *fileView {
	textView := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true).
		SetRegions(true)

		/*
			textView.SetInputCapture(
				func(event *tcell.EventKey) *tcell.EventKey {
					r, c := textView.GetScrollOffset()
					fmt.Printf("%d|%d", r, c)
					return event
				})
		*/

	fieldView := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(false).
		SetWrap(false).
		SetRegions(false)
	fieldView.SetBorder(true).SetTitle("Captures")

	container := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(textView, 0, 4, false).
		AddItem(fieldView, 30, 1, false)

	fv := fileView{
		Flex:      container,
		textView:  textView,
		fieldView: fieldView,
		rawText:   content,
		fileName:  fName,
	}
	return &fv
}

func (fv *fileView) FileName() string {
	return fv.fileName
}
