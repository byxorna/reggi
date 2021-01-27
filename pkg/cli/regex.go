package cli

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (c *cli) UpdateView(regexInput string) {
	if regexInput == "" {
		c.infoView.SetText("Enter a regex").SetTextColor(tcell.ColorViolet)
	}

	re, err := regexp.Compile(regexInput)
	if err != nil {
		c.ShowError(err)
	} else {
		c.infoView.SetText(fmt.Sprintf("%+v", re)).SetTextColor(tcell.ColorTeal)
	}
	c.HandleFilter(re)
	c.Application.Draw()
}

func (c *cli) ShowError(err error) {
	c.infoView.SetText(fmt.Sprintf("%v", err)).SetTextColor(tcell.ColorRed)
}

func (c *cli) HandleFilter(re *regexp.Regexp) {
	focusedFile, _ := c.pages.GetFrontPage()
	fv := c.fileViews[focusedFile]

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

		var prevHighlight string
		var currentHighlight string
		for n := 0; n < len(rawline); n++ {
			prevHighlight = currentHighlight
			currentHighlight = ""
			for matchID, match := range capMatches {
				for i := 0; i < len(match)/2; i++ {
					if n >= match[i] && n < match[i+1] {
						currentHighlight = NewHighlightID(lineNo, matchID, i).String()
						break
					}
				}
				if currentHighlight != "" {
					break
				}
			}

			if prevHighlight != currentHighlight {
				if currentHighlight != "" {
					line += fmt.Sprintf(`["%s"][blue]`, currentHighlight)
					highlightids = append(highlightids, currentHighlight)
				}
				if currentHighlight == "" {
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
			x[i] = fmt.Sprintf(" => %d: %s", i, f)
		}
		txt += fmt.Sprintf("%d:\n%s\n", lineNo, strings.Join(x, "\n"))
	}
	fv.fieldView.SetText(txt).ScrollToBeginning()
}
