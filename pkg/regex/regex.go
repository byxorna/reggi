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

// from https://golang.org/pkg/regexp/#Regexp.FindAllStringSubmatch
// If 'Submatch' is present, the return value is a slice identifying the successive submatches of the expression. Submatches are matches of parenthesized subexpressions (also known as capturing groups) within the regular expression, numbered from left to right in order of opening parenthesis. Submatch 0 is the match of the entire expression, submatch 1 the match of the first parenthesized subexpression, and so on.
// If 'Index' is present, matches and submatches are identified by byte index pairs within the input string: result[2*n:2*n+1] identifies the indexes of the nth submatch. The pair for n==0 identifies the match of the entire expression. If 'Index' is not present, the match is identified by the text of the match/submatch. If an index is negative or text is nil, it means that subexpression did not match any string in the input. For 'String' versions an empty string means either no match or an empty match.

func ExtractMatches(re *regexp.Regexp, input string) []LineMatches {
	results := []LineMatches{}
	if re == nil {
		return results
	}
	for lineNum, line := range strings.Split(input, "\n") {
		// TODO handle multiline
		m := re.FindAllString(line, -1)
		if m == nil || len(m) == 0 {
			// nil means no matches
			continue
		}
		match := LineMatches{
			LineNum: lineNum,
			//RawText: line,
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
