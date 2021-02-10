package regex

import (
	"fmt"
	"regexp"
	"strings"
)

type Capture struct {
	Index          int
	ByteIndexStart int
	ByteIndexEnd   int
	Name           string
	Extract        string
}

type LineMatches struct {
	LineNum int
	// RawText is the full text that was matched against
	RawText string
	// Coptures match the capture ID (number, or perl-style named capture)
	// to the text capturing it
	ExpressionMatch Capture
	Submatches      []Capture
}

// from https://golang.org/pkg/regexp/#Regexp.FindAllStringSubmatch
// If 'Submatch' is present, the return value is a slice identifying the successive submatches of the expression. Submatches are matches of parenthesized subexpressions (also known as capturing groups) within the regular expression, numbered from left to right in order of opening parenthesis. Submatch 0 is the match of the entire expression, submatch 1 the match of the first parenthesized subexpression, and so on.
// If 'Index' is present, matches and submatches are identified by byte index pairs within the input string: result[2*n:2*n+1] identifies the indexes of the nth submatch. The pair for n==0 identifies the match of the entire expression. If 'Index' is not present, the match is identified by the text of the match/submatch. If an index is negative or text is nil, it means that subexpression did not match any string in the input. For 'String' versions an empty string means either no match or an empty match.

func ExtractMatches(re *regexp.Regexp, multiline bool, input string) []LineMatches {
	results := []LineMatches{}
	if re == nil {
		return results
	}
	var inputsToMatch []string
	if multiline {
		inputsToMatch = []string{input}
	} else {
		inputsToMatch = strings.SplitAfter(input, "\n")
	}
	for lineNum, line := range inputsToMatch {
		// TODO handle multiline
		bytesLine := []byte(line)
		m := re.FindSubmatchIndex(bytesLine)

		if m == nil || len(m) == 0 {
			// nil means no matches
			continue
		}

		lm := LineMatches{LineNum: lineNum}
		for n := 0; n < len(m)/2; n++ {
			if n == 0 {
				expression := Capture{
					Name:           fmt.Sprintf("%d", n),
					ByteIndexStart: m[2*n],
					ByteIndexEnd:   m[2*n+1],
					Extract:        string(bytesLine[m[2*n]:m[2*n+1]]),
				}
				lm.ExpressionMatch = expression
			} else {
				submatch := Capture{
					Name:           fmt.Sprintf("%d", n-1),
					ByteIndexStart: m[2*n],
					ByteIndexEnd:   m[2*n+1],
					Extract:        string(bytesLine[m[2*n]:m[2*n+1]]),
				}
				lm.Submatches = append(lm.Submatches, submatch)
			}
		}

		// TODO handle captures with capMatches := re.FindAllStringSubmatch(line, -1)
		results = append(results, lm)
	}
	return results
}
