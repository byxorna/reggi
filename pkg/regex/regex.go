package regex

import (
	"fmt"
	"regexp"
)

type Capture struct {
	Index          int
	ByteIndexStart int
	ByteIndexEnd   int
	Name           string
	Extract        string
}

type LineMatches struct {
	//LineNum int
	// RawText is the full text that was matched against
	RawText string
	// Coptures match the capture ID (number, or perl-style named capture)
	// to the text capturing it
	Expressions []Capture
	Submatches  []Capture
}

// from https://golang.org/pkg/regexp/#Regexp.FindAllStringSubmatch
// If 'Submatch' is present, the return value is a slice identifying the successive submatches of the expression. Submatches are matches of parenthesized subexpressions (also known as capturing groups) within the regular expression, numbered from left to right in order of opening parenthesis. Submatch 0 is the match of the entire expression, submatch 1 the match of the first parenthesized subexpression, and so on.
// If 'Index' is present, matches and submatches are identified by byte index pairs within the input string: result[2*n:2*n+1] identifies the indexes of the nth submatch. The pair for n==0 identifies the match of the entire expression. If 'Index' is not present, the match is identified by the text of the match/submatch. If an index is negative or text is nil, it means that subexpression did not match any string in the input. For 'String' versions an empty string means either no match or an empty match.

func ExtractMatches(re *regexp.Regexp, matchall bool, input string) *LineMatches {
	if re == nil {
		return nil
	}
	var lm *LineMatches
	bytesLine := []byte(input)
	var m [][]int
	if matchall {
		m = re.FindAllSubmatchIndex(bytesLine, -1)
	} else {
		singleMatch := re.FindSubmatchIndex(bytesLine)
		m = [][]int{singleMatch}
	}

	if m == nil || len(m) == 0 {
		return nil
	}

	lm = &LineMatches{}
	for _, match := range m {
		for n := 0; n < len(match)/2; n++ {
			if n == 0 {
				expression := Capture{
					Name:           fmt.Sprintf("%d", n),
					ByteIndexStart: match[2*n],
					ByteIndexEnd:   match[2*n+1],
					Extract:        string(bytesLine[match[2*n]:match[2*n+1]]),
				}
				lm.Expressions = append(lm.Expressions, expression)
			} else {
				if match[2*n] == -1 || match[2*n+1] == -1 {
					// skip repeating captures with negative indicies
					continue
				}
				submatch := Capture{
					Name:           fmt.Sprintf("%d", n-1),
					ByteIndexStart: match[2*n],
					ByteIndexEnd:   match[2*n+1],
					Extract:        string(bytesLine[match[2*n]:match[2*n+1]]),
				}
				lm.Submatches = append(lm.Submatches, submatch)
			}
		}
	}
	return lm
}
