package ui

import (
	"bufio"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/paginator"
	input "github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	headerHeight               = 5
	footerHeight               = 3
	useHighPerformanceRenderer = true
)

type Model struct {
	ready bool

	textInput      input.Model
	paginationView paginator.Model
	viewport       viewport.Model
	err            error

	focusedTab int
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
	textInput.CharLimit = 156
	textInput.Width = 50
	textInput.Focus()

	paginationView := paginator.NewModel()
	paginationView.TotalPages = len(inputFiles)

	//vp := viewport.Model{
	//	YOffset: 0,
	//	//Height:                   10,
	//	HighPerformanceRendering: false,
	//}
	//vp.SetContent("testing\nhello\n")

	return &Model{
		textInput:      textInput,
		paginationView: paginationView,
		focusedTab:     0,
		inputFiles:     inputFiles,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return input.Blink
}

func (m *Model) focusedContents() string {
	return m.inputFiles[m.focusedTab].contents
}

func (m Model) View() string {
	return fmt.Sprintf(
		"Loaded %d files: %d %v\n\n%s\n\n%s\n%s\n%s",
		len(m.inputFiles),
		len(m.focusedContents()),
		m.inputFiles,
		m.textInput.View(),
		"(ctrl+c to quit)",
		m.viewport.View(),
		m.paginationView.View(),
	) + "\n"
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			return m, tea.Quit
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
			m.viewport.SetContent(m.inputFiles[m.focusedTab].contents)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMargins
		}

		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	if useHighPerformanceRenderer {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
