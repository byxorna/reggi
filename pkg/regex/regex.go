package regex

import (
	"regexp"
	"strings"
)

type Capture struct {
	Index   int
	Name    string
	Extract string
}

type LineMatches struct {
	LineNum int
	// RawText is the full text that was matched against
	RawText string
	// Coptures match the capture ID (number, or perl-style named capture)
	// to the text capturing it
	Captures []Capture
	Matches  []string
}

func ProcessText(re *regexp.Regexp, input string) []LineMatches {
	results := []LineMatches{}
	for lineNum, line := range strings.Split(input, "\n") {
		// TODO handle multiline
		m := re.FindAllString(line, -1)
		if m == nil || len(m) == 0 {
			// nil means no matches
			continue
		}
		match := LineMatches{
			LineNum: lineNum,
			RawText: line,
			Matches: m,
		}

		// TODO handle captures with capMatches := re.FindAllStringSubmatch(line, -1)
		results = append(results, match)
	}
	return results
}

/*
	// TODO handle multiline stuff
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

}
*/
