package ui

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/byxorna/regtest/pkg/regex"
	"github.com/charmbracelet/bubbles/paginator"
	input "github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	headerHeight               = 3 // TODO: this needs to be dynamic or it screws up redraw of the pager
	footerHeight               = 2
	useHighPerformanceRenderer = false // TODO: this doesnt work so hot right now

	infoClearDuration = 3 * time.Second
)

const (
	prompt = "> "
)

type focusType int

const (
	focusInput focusType = iota
	focusPager
)

type Model struct {
	ready           bool
	focus           focusType
	page            int
	multiline       bool
	caseInsensitive bool

	previousInput string
	textInput     input.Model
	pageDots      paginator.Model
	viewport      viewport.Model

	re         *regexp.Regexp
	err        error
	info       string
	updateTime time.Time

	inputFiles []*inputFile

	lineMatches []regex.LineMatches
}

func New(files []string) (*Model, error) {
	inputFiles := []*inputFile{}
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "Reading from stdin...\n")
		f, err := NewInputFile("/dev/stdin", bufio.NewReader(os.Stdin))
		if err != nil {
			return nil, err
		}
		inputFiles = append(inputFiles, f)
	} else {
		for _, src := range files {
			reader, err := os.Open(src)
			if err != nil {
				return nil, err
			}
			f, err := NewInputFile(src, reader)
			if err != nil {
				return nil, err
			}
			inputFiles = append(inputFiles, f)
		}
	}

	textInput := input.NewModel()
	textInput.Placeholder = "enter a regex (https://golang.org/pkg/regexp/syntax/)"
	textInput.CharLimit = 156
	textInput.Width = 50
	textInput.Prompt = getPrompt(true, false, false)
	textInput.Focus()

	pageDots := paginator.NewModel()
	pageDots.TotalPages = len(inputFiles)
	pageDots.Type = paginator.Dots

	return &Model{
		textInput:  textInput,
		pageDots:   pageDots,
		inputFiles: inputFiles,
		updateTime: time.Now(),
	}, nil
}

func (m Model) Init() tea.Cmd {
	return m.SetFocus(focusInput)
}

func getPrompt(focused, multiline, insensitive bool) string {
	// prefix prompt with our indicators for mode
	modes := ""
	if multiline {
		modes += redFg("m")
		modes += darkGrayFg("s")
	} else {
		modes += darkGrayFg("m")
		modes += greenFg("s")
	}
	if insensitive {
		modes += yellowFg("i")
	} else {
		modes += darkGrayFg("i")
	}
	localPrompt := fmt.Sprintf(" %4s ", modes)

	if focused {
		return localPrompt + fuchsiaFg(prompt)
	}
	return localPrompt + midGrayFg(prompt)
}

func (m *Model) SetFocus(f focusType) tea.Cmd {
	m.focus = f
	m.textInput.Prompt = getPrompt(m.focus == focusInput, m.multiline, m.caseInsensitive)

	switch m.focus {
	case focusInput:
		m.textInput.Focus()
		m.SetInfo("Focus: " + fuchsiaFg("input"))
		return input.Blink
	default:
		m.textInput.Blur()
		m.SetInfo("Focus: " + greenFg("pager"))
		return nil
	}
}

func (m *Model) focusedFile() *inputFile {
	return m.inputFiles[m.pageDots.Page]
}

func (m *Model) getHighlightedFileContents() string {
	c := m.inputFiles[m.pageDots.Page].contents
	m.lineMatches = regex.ExtractMatches(m.re, c)
	if len(m.lineMatches) == 0 {
		return c
	}

	// highlight text and return that
	highlightedText := ""
	lmIdx := 0
	for lineNum, line := range strings.Split(c, "\n") {
		if len(m.lineMatches) > lmIdx && m.lineMatches[lmIdx].LineNum == lineNum {
			lm := m.lineMatches[lmIdx]
			hl := ""
			hl += fmt.Sprintf("%+v\n", lm)
			rem := line
			for _, m := range lm.Matches {
				splits := strings.SplitN(rem, m, 2)
				switch len(splits) {
				case 0:
					hl += rem
				case 1:
					hl += m
					rem = ""
				default:
					hl += splits[0] + highlightStyle(m)
					rem = splits[1]
				}

			}
			highlightedText += hl + rem + "\n"
		} else {
			// take line as is - no matches
			highlightedText += line + "\n"
		}
	}
	return highlightedText
}

func (m *Model) updateViewportContents() {
	m.viewport.SetContent(m.getHighlightedFileContents())
	if m.page != m.pageDots.Page {
		// TODO: stash scroll offsets for focused file so we can restore when paging back
		m.viewport.YOffset = 0
		m.viewport.YPosition = 0
		m.page = m.pageDots.Page
	}
}

func (m *Model) SetInfo(info string) {
	m.updateTime = time.Now()
	m.info = info
}
