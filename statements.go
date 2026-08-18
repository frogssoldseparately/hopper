package hopper

import (
	"slices"

	"github.com/dlclark/regexp2/v2"
)

type Statement struct {
	id         int
	components []*Expression
	line       int
	column     int
}

const (
	scopeId int = 0
)

var statementMatchers = map[int]*regexp2.Regexp{
	scopeId: regexp2.MustCompile(`^\x02(\x00|\x01)*\x03`, regexOptions),
}

func filterExpressions(expressions *[]*Expression, filters []int) {
	replacement := []*Expression{}
	for _, expression := range *expressions {
		if !slices.Contains(filters, expression.id) {
			replacement = append(replacement, expression)
		}
	}
	*expressions = replacement
}

func (parser *Parser) makeStatement(id int, components []*Expression) *Statement {
	ignore := false
	empty := false
	var filters []int
	line := components[0].line
	column := components[0].column

	switch id {
	default:
		// do nothing
	}
	if ignore {
		return nil
	}
	if empty {
		components = []*Expression{}
	}
	if len(components) > 0 && len(filters) > 0 {
		filterExpressions(&components, filters)
	}
	return &Statement{id, components, line, column}
}

func (parser *Parser) readStatement(raw *[]rune, expressions *[]*Expression) (int, *Statement) {
	id := -1
	matchLength := -1
	for key := 0; key <= scopeId; key++ {
		regex := statementMatchers[key]
		if match, _ := regex.FindRunesMatch(*raw); match != nil {
			id = key
			matchLength = match.RuneLength
			break
		}
	}
	if id == -1 {
		return matchLength, nil
	}
	return matchLength, parser.makeStatement(id, (*expressions)[0:matchLength])
}
