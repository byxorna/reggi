package cli

import (
	"fmt"
	"regexp"
	"strings"

	tcell "github.com/gdamore/tcell/v2"
)

func (c *cli) UpdateInfo(txt string, re *regexp.Regexp, err error) {
	if err != nil {
		c.infoView.SetText(fmt.Sprintf("%v", err)).
			SetTextColor(tcell.ColorRed)
	} else {
		c.infoView.SetText(fmt.Sprintf("%+v", re)).SetTextColor(tcell.ColorTeal)
		c.HandleFilter(re, txt)
	}
}

func (c *cli) HandleFilter(re *regexp.Regexp, input string) {
	processedText := ""
	highlights := []string{}
	for lineNo, rawline := range strings.Split(c.rawText, "\n") {
		matches := re.FindAllStringIndex(rawline, -1)
		offset := 0
		line := ""
		for matchID, match := range matches {
			regionID := fmt.Sprintf("%d:%d", lineNo, matchID)
			highlights = append(highlights, regionID)
			line += fmt.Sprintf(`%s["%s"]%s[""]`,
				rawline[offset:match[0]],
				regionID,
				rawline[match[0]:match[1]])
			offset = match[1]
		}
		if len(matches) == 0 {
			line = rawline
		} else {
			line += rawline[offset:len(rawline)]
		}
		processedText += line + "\n"
	}
	c.textView.Highlight(highlights...)
	c.textView.SetText(processedText)
}
