package cli

import (
	"fmt"
	"regexp"

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
	//processedText := ""
	//for lineNo, rawline := range strings.Split(c.rawText, "\n") {
	//	matches := re.FindAllStringIndex(rawline, -1)
	//	offset := 0
	//	line := "* "
	//	for matchID, match := range matches {
	//		regionID := fmt.Sprintf("%d:%d", lineNo, matchID)
	//		line += fmt.Sprintf(`%s["%s"]%s[""]`,
	//			rawline[offset:offset+match[0]],
	//			regionID,
	//			rawline[offset+match[0]:offset+match[1]])
	//	}
	//	if len(matches) == 0 {
	//		line = rawline
	//	}
	//	processedText += "\n" + line
	//}
	//c.textView.SetText(processedText)
	c.textView.SetText(c.rawText)
}
