package ui

import (
	"fmt"
	"regexp"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}
	var lines []string // lines that change in viewport
	viewportUpdated := false
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			switch m.focus {
			case focusPager:
				switch msg.String() {
				case `q`:
					return m, tea.Quit
				case `i`, `a`, `A`, `I`, `o`, `O`:
					cmd = m.SetFocus(focusInput)
				case "home", "g":
					viewportUpdated = true
					lines = m.viewport.GotoTop()
					if m.viewport.HighPerformanceRendering {
						cmds = append(cmds, viewport.ViewUp(m.viewport, lines))
					}
				case "end", "G":
					viewportUpdated = true
					lines = m.viewport.GotoBottom()
					if m.viewport.HighPerformanceRendering {
						cmds = append(cmds, viewport.ViewDown(m.viewport, lines))
					}
				case "ctrl+f":
					viewportUpdated = true
					lines = m.viewport.HalfViewDown()
					if m.viewport.HighPerformanceRendering {
						cmds = append(cmds, viewport.ViewDown(m.viewport, lines))
					}
				case "ctrl+b":
					viewportUpdated = true
					lines = m.viewport.HalfViewUp()
					if m.viewport.HighPerformanceRendering {
						cmds = append(cmds, viewport.ViewUp(m.viewport, lines))
					}
				case "down", "j":
					viewportUpdated = true
					lines = m.viewport.LineDown(1)
					if m.viewport.HighPerformanceRendering {
						cmds = append(cmds, viewport.ViewDown(m.viewport, lines))
					}
				case "up", "k":
					viewportUpdated = true
					lines = m.viewport.LineUp(1)
					if m.viewport.HighPerformanceRendering {
						cmds = append(cmds, viewport.ViewUp(m.viewport, lines))
					}
				case "H":
					viewportUpdated = true
					m.pageDots.PrevPage()
					m.pageDots, cmd = m.pageDots.Update(msg)
					cmds = append(cmds, cmd)
					m.SetInfo("Tab " + indigoFg(m.focusedFile().source))
				case "L":
					viewportUpdated = true
					m.pageDots.NextPage()
					m.pageDots, cmd = m.pageDots.Update(msg)
					cmds = append(cmds, cmd)
					m.SetInfo("Tab " + indigoFg(m.focusedFile().source))
				}

			case focusInput:
				switch msg.Type {
				case tea.KeyCtrlI:
					m.caseInsensitive = !m.caseInsensitive
					en := "enabled"
					if !m.caseInsensitive {
						en = "disabled"
					}
					m.SetInfo("Case insensitive matching " + yellowFg(en))
					m.UpdatePrompt()
				case tea.KeyCtrlL:
					m.multiline = !m.multiline
					en := greenFg("single line")
					if m.multiline {
						en = redFg("multiline")
					}
					m.SetInfo("Matching set to " + en)
					m.UpdatePrompt()

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
		if !m.ready {
			// Since this program is using the full size of the viewport we need
			// to wait until we've received the window dimensions before we
			// can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			viewportUpdated = m.initializeViewport(msg.Width, msg.Height)
		} else {
			viewportUpdated = m.resizeViewport(msg.Width, msg.Height)
		}

		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// this handling needs to by synchronous because it modifies application state
	// and we cannot defer this to a Cmd
	shouldUpdate := m.HandleInput()

	if shouldUpdate || viewportUpdated {
		m.updateViewportContents()
		if m.viewport.HighPerformanceRendering {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}
	m.HandleUpdateTime()

	return m, tea.Batch(cmds...)
}

// HandleInput processes input text, compiles regex, and possibly
// mutates the content in the view if matches are found (for highlighting)
// This should run as fast as possible with no changes, to avoid busy looping.
// This should only compile regex/apply matches if input has changed
func (m *Model) HandleInput() (shouldUpdate bool) {
	currentValue := m.textInput.Value()
	if m.previousInput == currentValue {
		// skip compilation if no value change
		return false
	}
	m.previousInput = currentValue
	// dont process empty regex
	if currentValue == "" {
		m.re = nil
		m.err = nil
		return true
	}
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
	m.re, m.err = regexp.Compile(fmt.Sprintf(`%s%s`, flags, currentValue))
	return true
}

func (m *Model) UpdatePrompt() {
	m.textInput.Prompt = getPrompt(m.focus == focusInput, m.multiline, m.caseInsensitive)
}

func (m *Model) HandleUpdateTime() {
	if time.Since(m.updateTime) > infoClearDuration {
		m.info = ""
	}
}
