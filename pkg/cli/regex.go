package cli

import (
	"fmt"
	"regexp"
	"strings"

	tcell "github.com/gdamore/tcell/v2"
)

func (c *cli) UpdateView(txt string) {
	if txt == "" {
		c.infoView.SetText("Enter a regex").SetTextColor(tcell.ColorViolet)
		c.HandleFilter(nil, txt)
		return
	}
	re, err := regexp.Compile(txt)
	if err != nil {
		c.infoView.SetText(fmt.Sprintf("%v", err)).
			SetTextColor(tcell.ColorRed)
	} else {
		c.infoView.SetText(fmt.Sprintf("%+v", re)).SetTextColor(tcell.ColorTeal)
	}
	c.HandleFilter(re, txt)
	c.Application.Draw()
}

func (c *cli) HandleFilter(re *regexp.Regexp, input string) {
	// populate the text view with fields highlighted
	processedText := ""
	highlightids := []string{}
	lines := strings.Split(c.rawText, "\n")
	matchingCaptures := make([]map[int]string, len(lines)) // capture fields
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

		//matchingCaptures[lineNo] = captures
		processedText += line + "\n"
	}
	c.textView.Highlight(highlightids...)
	c.textView.SetText(processedText)

	// for the fields view, for the currently selected lines, show the matches in a list
	txt := ""
	for lineNo, captures := range matchingCaptures {
		if len(captures) == 0 {
			continue
		}
		x := make([]string, len(captures))
		for i, f := range captures {
			x[i] = fmt.Sprintf(" => %d: %s", i, f)
		}
		txt += fmt.Sprintf("%d:\n%s\n", lineNo, strings.Join(x, "\n"))
	}
	c.fieldView.SetText(txt).ScrollToBeginning()
}
