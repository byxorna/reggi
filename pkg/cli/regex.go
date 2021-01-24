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
	highlights := []string{}
	lines := strings.Split(c.rawText, "\n")
	matchingFields := make([][]string, len(lines))
	for lineNo, rawline := range lines {
		if re == nil {
			processedText += fmt.Sprintf("%d: %s\n", lineNo, rawline)
			continue
		}
		matches := re.FindAllStringIndex(rawline, -1)
		offset := 0
		line := ""
		fields := make([]string, len(matches))
		for matchID, match := range matches {
			regionID := fmt.Sprintf("%d:%d", lineNo, matchID)
			highlights = append(highlights, regionID)
			line += fmt.Sprintf(`%s["%s"]%s[""]`,
				rawline[offset:match[0]],
				regionID,
				rawline[match[0]:match[1]])
			offset = match[1]

			fields[matchID] = rawline[match[0]:match[1]]
		}
		matchingFields[lineNo] = fields
		if len(matches) == 0 {
			line = rawline
		} else {
			line += rawline[offset:len(rawline)]
		}
		processedText += fmt.Sprintf("%d: %s\n", lineNo, line)
	}
	c.textView.Highlight(highlights...)
	c.textView.SetText(processedText)

	// for the fields view, for the currently selected lines, show the matches in a list
	txt := ""
	for lineNo, fields := range matchingFields {
		x := make([]string, len(fields))
		for i, f := range fields {
			x[i] = fmt.Sprintf(" => %d: %s", i, f)
		}
		if len(fields) > 0 {
			txt += fmt.Sprintf("%d:\n%s\n", lineNo, strings.Join(x, "\n"))
		}
	}
	c.fieldView.SetText(txt).ScrollToBeginning()
}
