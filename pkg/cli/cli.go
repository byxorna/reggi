package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/byxorna/regtest/pkg/input"
	"github.com/byxorna/regtest/pkg/version"
	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	InputCompileDelay = 300 * time.Millisecond
)

var (
	defaultRegex = regexp.MustCompile(`Enter a regex`)
)

type CLI interface {
	Run() error
}

type cli struct {
	*tview.Application
	layout    *tview.Flex
	infoView  *tview.TextView
	inputView *tview.InputField
	pages     *tview.Pages

	inputChan chan string

	fileViews map[string]*fileView
}

// houses the view for a single file's contents and matches
type fileView struct {
	textView  *tview.TextView
	fieldView *tview.TextView
	rawText   string
}

func New(files []string) CLI {
	c := cli{
		inputChan: make(chan string),
		fileViews: map[string]*fileView{},
	}
	c.Application = tview.NewApplication()
	c.layout = tview.NewFlex()
	c.infoView = tview.NewTextView().
		SetScrollable(false)
	c.infoView.SetBorderPadding(0, 0, 1, 1).
		SetBorder(true)

	c.inputView = inputView()
	c.inputView.SetBorder(true)

	c.pages = tview.NewPages()
	for _, f := range files {
		err := c.OpenFile(f)
		if err != nil {
			c.ShowError(err)
		}
	}
	// TODO handle no file here
	c.pages.ShowPage(files[0])

	c.layout.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(c.inputView, 0, 1, true).
			AddItem(c.infoView, 0, 1, false), 3, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(c.pages, 0, 1, false), 0, 5, false), 0, 1, false).SetFullScreen(true)
	c.Application.SetRoot(c.layout, false).SetFocus(c.inputView)

	c.layout.SetTitle(c.windowTitle()).SetBorder(true)

	go input.Debounce(InputCompileDelay, c.inputChan, func(txt string) {
		c.UpdateView(txt)
	})
	// debounce keystrokes and aggregate evnts to compile regex after delay
	c.inputView.SetChangedFunc(func(txt string) {
		c.inputChan <- txt
	})
	return &c
}

func (c *cli) OpenFile(f string) error {
	fh, err := os.Open(f)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(fh)
	if err != nil {
		return err
	}

	textView := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true).
		SetRegions(true)
	//SetChangedFunc(func() {
	//		c.Application.Draw()
	//	})
	fieldView := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(false).
		SetWrap(false).
		SetRegions(false)
	fieldView.SetBorder(true).SetTitle("Captures")

	fv := fileView{
		textView:  textView,
		fieldView: fieldView,
		rawText:   string(data),
	}
	p := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(fv.textView, 0, 4, false).
		AddItem(fv.fieldView, 30, 1, false)

	c.fileViews[f] = &fv

	c.pages.AddPage(f, p, true, false)
	return nil
}

func (c *cli) Run() error {
	c.HandleFilter(defaultRegex)
	return c.Application.Run()
}

func inputView() *tview.InputField {
	f := tview.NewInputField()
	f.SetLabel("r/")
	f.SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite).
		SetTitle("Regex")
	f.SetBorder(true).SetBorderPadding(0, 0, 1, 1)
	return f
}

func (c *cli) windowTitle() string {
	focusedFile, _ := c.pages.GetFrontPage()
	numFiles := c.pages.GetPageCount()
	return fmt.Sprintf("[%d] Regtest: %s (%s)", numFiles, focusedFile, version.Version)
}
