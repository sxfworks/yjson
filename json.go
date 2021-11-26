package variable_json

import "fmt"

func Parse(input []byte) (*JsonNode, error) {
	if len(input) <= 0 {
		return nil, fmt.Errorf("input invalid")
	}

	l := newLex(input)
	_ = yyParse(l)

	if l.err != nil {
		return nil, l.err
	}

	return l.result.(*JsonNode), nil
}
