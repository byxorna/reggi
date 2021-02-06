package ui

import (
	"fmt"
	"regexp"

	"github.com/byxorna/regtest/pkg/regex"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	m.err = nil
	switch msg := msg.(type) {
	case error:
		m.err = msg
	case regexp.Regexp:
		cmds = append(cmds, func() tea.Msg {
			m.UpdateContent(msg)
			return nil
		})
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			switch m.focus {
			case focusPager:

				needSync := false
				switch msg.String() {
				case `q`:
					return m, tea.Quit
				case `i`, `a`, `A`, `I`, `o`, `O`:
					cmd := m.SetFocus(focusInput)
					cmds = append(cmds, cmd)
				case "home", "g":
					m.viewport.GotoTop()
					needSync = true
				case "end", "G":
					m.viewport.GotoBottom()
					needSync = true
				case "ctrl+f":
					m.viewport.HalfViewDown()
					needSync = true
				case "ctrl+b":
					m.viewport.HalfViewUp()
					needSync = true
				case "H":
					m.pageDots.PrevPage()
				case "L":
					m.pageDots.NextPage()
				}

				m.pageDots, cmd = m.pageDots.Update(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
				m.updateViewportContents()

				m.viewport, cmd = m.viewport.Update(msg)
				cmds = append(cmds, cmd)
				if needSync && m.viewport.HighPerformanceRendering {
					cmds = append(cmds, viewport.Sync(m.viewport))
				}

			case focusInput:
				switch msg.Type {
				case tea.KeyCtrlC:
					return m, tea.Quit
				case tea.KeyEsc:
					cmd := m.SetFocus(focusPager)
					cmds = append(cmds, cmd)
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

	// TODO: move this compilation into an async Cmd
	cmds = append(cmds, func() tea.Msg {
		return m.CompileInput()
	})

	return m, tea.Batch(cmds...)
}

func (m *Model) CompileInput() tea.Msg {
	currentValue := m.textInput.Value()
	if m.previousInput == currentValue {
		// skip compilation if no value change
		return nil
	}
	m.previousInput = currentValue
	flags := ""
	if m.multiline {
		flags += "m"
	} else {
		flags += "s"
	}
	if m.caseInsensitive {
		flags += "i"
	}
	if len(flags) > 0 {
		flags = "(?" + flags + ")"
	}
	//m.regex, m.err = regexp.Compile(fmt.Sprintf(`%s%s`, flags, currentValue))
	re, err := regexp.Compile(fmt.Sprintf(`%s%s`, flags, currentValue))
	if err != nil {
		return err
	}
	return re
}

func (m *Model) UpdateContent(re regexp.Regexp) tea.Msg {
	regex.ProcessText(re, m.focusedFile().contents)
	return nil
}
