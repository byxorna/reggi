package ui

import (
	"fmt"
	"strings"

	"github.com/byxorna/reggi/pkg/version"
	"github.com/charmbracelet/bubbles/viewport"
	runewidth "github.com/mattn/go-runewidth"
)

var (
	helpJoin  = " • "
	pagerHelp = []string{
		`i: input`,
		`h,l: tab`,
		`j,k: scroll`,
		`g: top`,
		`G: bottom`,
		`q: quit`,
	}
	inputHelp = []string{
		`esc: pager`,
		`ctrl+y: all`,
		`ctrl+l: multiline`,
		`ctrl+s: spanline`,
		`ctrl+i: case`,
		`ctrl+c: quit`,
	}
)

func (m Model) View() string {
	infoField := ""
	if m.err != nil {
		infoField = redFg(m.err.Error())
	} else if m.info != "" {
		infoField = m.info
	}
	return strings.Join([]string{
		"\n" + m.textInput.View(),
		infoField,
		m.viewport.View(),
		m.formatLineSpread(
			fmt.Sprintf(`[%d/%d] %s`, m.pageDots.Page+1, m.pageDots.TotalPages, brightGrayFg(m.focusedFile().source)), 0,
			fmt.Sprintf(`%d%% %s (%s)`, int(m.viewport.ScrollPercent()*100), m.pageDots.View(), version.Version)),
		m.helpLine(),
	}, "\n")
}

func (m *Model) formatLineSpread(left string, extraSpace int, right string) string {
	// runewidth doesnt take into account non-printing characters, so provide a hack to let callers precompute
	// style widths for proper alignment of color text
	space := m.viewport.Width - runewidth.StringWidth(left) + extraSpace - runewidth.StringWidth(right)
	if space < 1 {
		space = 1
	}
	return fmt.Sprintf(`%s%s%s`, left, resetStyle(strings.Repeat(" ", space)), right)
}

func (m *Model) helpLine() string {
	h := ""
	mode := ""
	switch m.focus {
	case focusInput:
		h = midGrayFg(strings.Join(inputHelp, helpJoin))
		mode = fuchsiaFg("Input")
	case focusPager:
		h = midGrayFg(strings.Join(pagerHelp, helpJoin))
		mode = greenFg("Pager")
	}
	return m.formatLineSpread(h, 0, mode)
}

func (m *Model) resizeViewport(width, height int) bool {
	verticalMargins := headerHeight + footerHeight
	m.viewport.Width = width
	m.viewport.Height = height - verticalMargins
	return true
}

func (m *Model) initializeViewport(width, height int) bool {
	verticalMargins := headerHeight + footerHeight
	m.viewport = viewport.Model{Width: width, Height: height - verticalMargins}
	m.viewport.YPosition = headerHeight + 1
	m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
	m.viewport.SetContent(m.getHighlightedFileContents())
	m.ready = true
	return true
}
