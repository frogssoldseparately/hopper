package main

import (
	"slices"

	"github.com/dlclark/regexp2/v2"
)

type Expression struct {
	id         int
	components []*Token
	line       int
	column     int
}

const (
	fillerId     int = iota // 0
	patternId               // 1
	namedScopeId            // 2
	endblockId              // 3
)

var expressionMatchers = map[int]*regexp2.Regexp{
	fillerId:     regexp2.MustCompile(`^\x05\x02\x06`, regexOptions),
	patternId:    regexp2.MustCompile(`^\x05\x03\x06(\x04\x06)*\x08*`, regexOptions),
	namedScopeId: regexp2.MustCompile(`^\x05\x00`, regexOptions),
	endblockId:   regexp2.MustCompile(`^\x01`, regexOptions),
}

func filterTokens(tokens *[]*Token, filters []int) {
	replacement := []*Token{}
	for _, token := range *tokens {
		if !slices.Contains(filters, token.id) {
			replacement = append(replacement, token)
		}
	}
	*tokens = replacement
}

func (parser *Parser) makeExpression(id int, components []*Token) *Expression {
	ignore := false
	empty := false
	var filters []int
	line := components[0].line
	column := components[0].line

	switch id {
	case fillerId:
		filters = []int{equalsId}
	case patternId:
		filters = []int{colonId, pipeId}
	default:
		// do nothing
	}
	if ignore {
		return nil
	}
	if empty {
		components = []*Token{}
	}
	if len(components) > 0 && len(filters) > 0 {
		filterTokens(&components, filters)
	}
	return &Expression{id, components, line, column}
}

func (parser *Parser) readExpression(raw *[]rune, tokens *[]*Token) (int, *Expression) {
	id := -1
	matchLength := -1
	for key := 0; key <= endblockId; key++ {
		regex := expressionMatchers[key]
		if match, _ := regex.FindRunesMatch(*raw); match != nil {
			id = key
			matchLength = match.RuneLength
			break
		}
	}
	if id == -1 {
		return matchLength, nil
	}
	return matchLength, parser.makeExpression(id, (*tokens)[0:matchLength])
}
