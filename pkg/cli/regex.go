package cli

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (c *cli) UpdateView(rawRe string) {
	compiledRe, err := regexp.Compile(rawRe)
	c.UpdateInfoView(rawRe, compiledRe, err)
	c.HandleFilter(compiledRe)
}

func (c *cli) UpdateInfoView(rawRe string, compiledRe *regexp.Regexp, err error) {
	if rawRe == "" {
		c.infoView.SetText("Enter a regex").SetTextColor(tcell.ColorViolet)
	} else if err != nil {
		c.ShowError(err)
	} else {
		c.infoView.
			SetText(fmt.Sprintf("Compiled: %+v (%d captures)", compiledRe, compiledRe.NumSubexp())).
			SetTextColor(tcell.ColorTeal)
	}
}

func (c *cli) ShowError(err error) {
	c.infoView.SetText(fmt.Sprintf("%v", err)).SetTextColor(tcell.ColorRed)
}

func (c *cli) HandleFilter(re *regexp.Regexp) {
	fv := c.focusedFileView()

	// populate the text view with fields highlighted
	processedText := ""
	highlightids := []string{}
	lines := strings.Split(tview.Escape(fv.rawText), "\n")
	matchingCaptures := map[int][]string{}
	for lineNo, rawline := range lines {
		if re == nil {
			processedText += rawline + "\n"
			continue
		}
		capMatches := re.FindAllStringSubmatchIndex(rawline, -1)
		// TODO: bytesbuffer would be snappier
		line := ""

		var prevHighlight *HighlightID
		var currentHighlight *HighlightID
		for n := 0; n < len(rawline); n++ {
			prevHighlight = currentHighlight
			currentHighlight = nil
			for matchID, match := range capMatches {
				for i := 0; i < len(match)/2; i++ {
					if n >= match[i] && n < match[i+1] {
						currentHighlight = NewHighlightID(lineNo, matchID, i)
						break
					}
				}
				if currentHighlight != nil {
					break
				}
			}

			if prevHighlight != currentHighlight {
				if currentHighlight != nil {
					color := `blue` // normal highlights
					if currentHighlight.IsCapture() {
						color = `red`
					}
					line += fmt.Sprintf(`["%s"][%s]`, currentHighlight, color)
					highlightids = append(highlightids, currentHighlight.String())
				}
				if currentHighlight == nil {
					line += `[""][white]`
				}
			}
			line += string(rawline[n])
		}

		processedText += line + "\n"

		// now store the captures for viewing on the right panel
		captures := []string{}
		for _, match := range capMatches {
			if len(match) <= 2 {
				continue
			}
			for i := 0; i < len(match)/2; i++ {
				if i == 0 {
					continue
				}
				captures = append(captures, rawline[match[2*i]:match[2*i+1]])
			}
		}
		matchingCaptures[lineNo] = captures
	}
	fv.textView.Highlight(highlightids...)
	fv.textView.SetText(processedText)

	// for the fields view, for the currently selected lines, show the matches in a list
	linesWithCaptures := make([]int, len(matchingCaptures))
	i := 0
	for k := range matchingCaptures {
		linesWithCaptures[i] = k
		i++
	}
	sort.Ints(linesWithCaptures)
	txt := ""
	for lineNo := range linesWithCaptures {
		captures := matchingCaptures[lineNo]
		if len(captures) == 0 {
			continue
		}
		x := make([]string, len(captures))
		for i, f := range captures {
			x[i] = fmt.Sprintf("  match %d: %s", i, f)
		}
		txt += fmt.Sprintf("Line %d:\n%s\n", lineNo, strings.Join(x, "\n"))
	}
	if txt == "" {
		txt = "No captures"
	}
	fv.fieldView.SetText(txt).ScrollToBeginning()
}
