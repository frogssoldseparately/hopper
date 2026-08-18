package main

import (
	"strings"

	"github.com/dlclark/regexp2/v2"
)

type Token struct {
	id     int
	value  string
	line   int
	column int
}

const regexOptions = regexp2.Singleline

// Token IDs
const (
	leftBraceId  int = iota // 0
	rightBraceId            // 1
	equalsId                // 2
	colonId                 // 3
	pipeId                  // 4
	identId                 // 5
	stringId                // 6
	commentId               // 7
	injectId                // 8
	newlineId               // 9
	whitespaceId            // 10
)

var tokenMatchers = map[int]*regexp2.Regexp{
	leftBraceId:  regexp2.MustCompile(`^((\{))`, regexOptions),
	rightBraceId: regexp2.MustCompile(`^((\}))`, regexOptions),
	equalsId:     regexp2.MustCompile(`^((\=))`, regexOptions),
	colonId:      regexp2.MustCompile(`^((\:))`, regexOptions),
	pipeId:       regexp2.MustCompile(`^((\|))`, regexOptions),
	identId:      regexp2.MustCompile(`^(([a-zA-Z]+))`, regexOptions),
	stringId:     regexp2.MustCompile(`^((\x60[^\x60]*?\x60))`, regexOptions),
	commentId:    regexp2.MustCompile(`^((\/\/.*?(?=\n)))`, regexOptions),
	injectId:     regexp2.MustCompile(`^((>.*?(?=\n)))`, regexOptions),
	newlineId:    regexp2.MustCompile(`^((\n))`, regexOptions),
	whitespaceId: regexp2.MustCompile(`^((\s+))`, regexOptions),
}

func (parser *Parser) makeToken(id int, components []string) *Token {
	ignore := false
	empty := false
	value := ""
	trimChars := ""
	joinChars := ""
	leftTrim := 0
	rightTrim := 0
	bothTrim := 0

	switch id {
	case stringId:
		bothTrim = 1
	case commentId:
		fallthrough
	case whitespaceId:
		ignore = true
	case injectId:
		trimChars = "> "
	case newlineId:
		ignore = true
		parser.incrementLineNumber()
	default:
		// do nothing
	}

	if ignore {
		return nil
	}
	if !empty && len(value) == 0 {
		value = strings.Trim(strings.Join(components, joinChars), trimChars)
		if bothTrim != 0 {
			leftTrim = bothTrim
			rightTrim = bothTrim
		}
		if leftTrim+rightTrim+bothTrim > 0 {
			value = value[leftTrim : len(value)-rightTrim]
		}
	}
	return &Token{id, value, parser.lineNumber, parser.columnNumber}
}

func (parser *Parser) readToken(raw *string) (int, *Token) {
	id := -1
	matchLength := -1
	components := []string{}
	for key := 0; key <= whitespaceId; key++ {
		regex := tokenMatchers[key]
		if match, _ := regex.FindStringMatch(*raw); match != nil {
			id = key
			matchLength = match.RuneLength
			for _, component := range match.Captures {
				components = append(components, component.String())
			}
			break
		}
	}
	if id == -1 {
		return matchLength, nil
	}
	return matchLength, parser.makeToken(id, components)
}
