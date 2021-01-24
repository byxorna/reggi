package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/byxorna/regtest/pkg/input"
	"github.com/byxorna/regtest/pkg/version"
	"github.com/rivo/tview"
)

const (
	InputCompileDelay = 300 * time.Millisecond
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
	treeView  *tview.TreeView

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
		SetScrollable(false).SetDynamicColors(true)

	c.textView = tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			c.Application.Draw()
		})

	c.treeView = tview.NewTreeView()

	c.inputView = inputView()

	c.layout.AddItem(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(c.inputView, 3, 1, true).
			AddItem(c.infoView, 1, 1, false).
			AddItem(tview.NewGrid().
				SetRows(0).
				SetColumns(-4, -1).
				AddItem(c.textView, 0, 0, 1, 4, 10, 10, false), 0, 5, false),
		0, 1, false).
		SetFullScreen(true)
	c.Application.SetRoot(c.layout, false).SetFocus(c.inputView)

	c.loadFile(files)
	c.layout.SetTitle(c.windowTitle()).
		SetBorder(true)

	go input.Debounce(InputCompileDelay, c.inputChan, func(txt string) {
		re, err := regexp.Compile(txt)
		c.UpdateInfo(txt, re, err)
		c.Application.Draw()
	})
	// debounce keystrokes and aggregate evnts to compile regex after delay
	c.inputView.SetChangedFunc(func(txt string) {
		c.inputChan <- txt
	})
	return &c
}

func inputView() *tview.InputField {
	f := tview.NewInputField()
	f.SetBorder(true).SetTitle("Regex")
	return f
}

func (c *cli) loadFile(files []string) error {
	for _, f := range files {
		c.files = append(c.files, f)
	}
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
	c.textView.SetText(c.rawText)
	return nil
}

func (c *cli) windowTitle() string {
	return fmt.Sprintf("[%d:%d] Regview: %s (%s)", c.fileidx, len(c.files), c.activeFile.Name(), version.Version)
}
