package ui

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/byxorna/regtest/pkg/version"
	"github.com/charmbracelet/bubbles/paginator"
	input "github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	runewidth "github.com/mattn/go-runewidth"
	te "github.com/muesli/termenv"
)

var (
	headerHeight               = 4 // TODO: this needs to be dynamic or it screws up redraw of the pager
	footerHeight               = 1
	useHighPerformanceRenderer = false

	color         = te.ColorProfile().Color
	focusedPrompt = te.String("> ").Foreground(color("205")).String()
	blurredPrompt = "> "
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

	textInput      input.Model
	paginationView paginator.Model
	viewport       viewport.Model

	regex *regexp.Regexp
	err   error

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

	paginationView := paginator.NewModel()
	paginationView.TotalPages = len(inputFiles)
	paginationView.Type = paginator.Dots

	return &Model{
		textInput:      textInput,
		paginationView: paginationView,
		inputFiles:     inputFiles,
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
	return m.inputFiles[m.paginationView.Page]
}

func (m *Model) formatLineSpread(left, right string) string {
	space := m.viewport.Width - (runewidth.StringWidth(left) + runewidth.StringWidth(right))
	if space < 1 {
		space = 1
	}
	return fmt.Sprintf(`%s%s%s`, left, strings.Repeat(" ", space), right)
}

func (m Model) View() string {
	errStr := ""
	if m.err != nil {
		errStr = te.String(m.err.Error()).Foreground(color("197")).String()
	}
	statusLeft := fmt.Sprintf("Loaded %d files: %d", len(m.inputFiles), len(m.focusedFile().contents))
	statusRight := fmt.Sprintf(`%s %s`, version.Version, version.Commit[0:7])
	statusLine := m.formatLineSpread(statusLeft, statusRight)
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s",
		statusLine,
		m.textInput.View(),
		errStr,
		m.viewport.View(),
		m.formatLineSpread(
			fmt.Sprintf(`%s`, m.focusedFile().source),
			fmt.Sprintf(`%d%% %s`, int(m.viewport.ScrollPercent()*100), m.paginationView.View()),
		),
	) + "\n"
}

func (m *Model) updateViewportContents() {
	if m.page != m.paginationView.Page {
		m.viewport.SetContent(m.focusedFile().contents)
		m.viewport.YOffset = 0
		m.viewport.YPosition = 0
		m.page = m.paginationView.Page
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			switch m.focus {
			case focusPager:

				sync := false
				switch msg.String() {
				case `q`:
					return m, tea.Quit
				case `i`, `a`, `A`, `I`, `o`, `O`:
					return m.SetFocus(focusInput)
				case "home", "g":
					m.viewport.GotoTop()
					sync = true
				case "end", "G":
					m.viewport.GotoBottom()
					sync = true
				case "ctrl+f":
					m.viewport.HalfViewDown()
					sync = true
				case "ctrl+b":
					m.viewport.HalfViewUp()
					sync = true
				}

				m.paginationView, cmd = m.paginationView.Update(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
				m.updateViewportContents()

				m.viewport, cmd = m.viewport.Update(msg)
				cmds = append(cmds, cmd)
				if sync && m.viewport.HighPerformanceRendering {
					cmds = append(cmds, viewport.Sync(m.viewport))
				}

			case focusInput:
				switch msg.Type {
				case tea.KeyCtrlC:
					return m, tea.Quit
				case tea.KeyEsc:
					return m.SetFocus(focusPager)
				}
				m.textInput, cmd = m.textInput.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	case tea.WindowSizeMsg:
		// https://github.com/charmbracelet/bubbletea/blob/master/examples/pager/main.go#L95
		// We've reveived terminal dimensions, either for the first time or
		// after a resize
		verticalMargins := headerHeight + footerHeight
		if !m.ready {
			// Since this program is using the full size of the viewport we need
			// to wait until we've received the window dimensions before we
			// can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.Model{Width: msg.Width, Height: msg.Height - verticalMargins}
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.focusedFile().contents)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMargins
		}

		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	m.regex, m.err = regexp.Compile(m.textInput.Value())

	return m, tea.Batch(cmds...)
}
