package ui

import (
	"regexp"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

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
