package ui

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/charmbracelet/bubbles/paginator"
	input "github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	headerHeight               = 5 // TODO: this needs to be dynamic or it screws up redraw of the pager
	footerHeight               = 1
	useHighPerformanceRenderer = false

	focusedPrompt = fuchsiaFg("> ")
	blurredPrompt = midGrayFg("> ")
)

type focusType int

const (
	focusInput focusType = iota
	focusPager
)

type Model struct {
	ready bool
	focus focusType
	page  int

	textInput input.Model
	pageDots  paginator.Model
	viewport  viewport.Model

	regex           *regexp.Regexp
	err             error
	multiline       bool
	caseInsensitive bool

	inputFiles []*inputFile
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
	textInput.Placeholder = "enter a regex"
	textInput.Prompt = focusedPrompt
	textInput.CharLimit = 156
	textInput.Width = 50
	textInput.Focus()

	pageDots := paginator.NewModel()
	pageDots.TotalPages = len(inputFiles)
	pageDots.Type = paginator.Dots

	return &Model{
		textInput:  textInput,
		pageDots:   pageDots,
		inputFiles: inputFiles,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return input.Blink
}

func (m Model) SetFocus(f focusType) (Model, tea.Cmd) {
	m.focus = f
	switch m.focus {
	case focusInput:
		m.textInput.Focus()
		m.textInput.Prompt = focusedPrompt
		return m, input.Blink
	default:
		m.textInput.Blur()
		m.textInput.Prompt = blurredPrompt
		return m, nil
	}
}

func (m *Model) focusedFile() *inputFile {
	return m.inputFiles[m.pageDots.Page]
}

func (m *Model) updateViewportContents() {
	if m.page != m.pageDots.Page {
		m.viewport.SetContent(m.focusedFile().contents)
		m.viewport.YOffset = 0
		m.viewport.YPosition = 0
		m.page = m.pageDots.Page
	}
}
