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
	textView  *tview.TextView
	infoView  *tview.TextView
	inputView *tview.InputField
	fieldView *tview.TextView

	inputChan chan string
	files     []string

	activeFile *os.File
	rawText    string
	fileidx    int
}

func New(files []string) CLI {
	c := cli{
		inputChan: make(chan string),
	}
	c.Application = tview.NewApplication()
	c.layout = tview.NewFlex()
	c.infoView = tview.NewTextView().
		SetScrollable(false)
	c.infoView.SetBorderPadding(0, 0, 1, 1).
		SetBorder(true)

	c.textView = tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			c.Application.Draw()
		})

	c.fieldView = tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(false).
		SetWrap(false).
		SetRegions(false)
	c.fieldView.SetBorder(true).SetTitle("Captures")

	c.inputView = inputView()
	c.inputView.SetBorder(true)

	c.layout.AddItem(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(
				// top input bits
				tview.NewFlex().SetDirection(tview.FlexColumn).
					AddItem(c.inputView, 0, 1, true).
					AddItem(c.infoView, 0, 1, false), 3, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(c.textView, 0, 4, false).
				AddItem(c.fieldView, 30, 1, false), 0, 5, false), 0, 1, false).
		SetFullScreen(true)
	c.Application.SetRoot(c.layout, false).SetFocus(c.inputView)

	c.loadFile(files)
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

func (c *cli) Run() error {
	c.UpdateView("")
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

func (c *cli) loadFile(files []string) error {
	for _, f := range files {
		c.files = append(c.files, f)
	}
	if len(c.files) > 0 {
		activeFile, err := os.Open(c.files[c.fileidx])
		if err != nil {
			return err
		}
		c.activeFile = activeFile
		data, err := ioutil.ReadAll(c.activeFile)
		if err != nil {
			return err
		}
		c.rawText = string(data)
	}
	c.HandleFilter(defaultRegex)
	return nil
}

func (c *cli) windowTitle() string {
	fileName := "no file"
	if c.activeFile != nil {
		fileName = c.activeFile.Name()
	}
	return fmt.Sprintf("[%d:%d] Regview: %s (%s)", c.fileidx, len(c.files), fileName, version.Version)
}
