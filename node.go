package variable_json

import (
	"bytes"
	"fmt"
	"strconv"
)

type JsonNodeType int64

const (
	JsonNodeObject   JsonNodeType = 1
	JsonNodeArray    JsonNodeType = 2
	JsonNodeNumber   JsonNodeType = 3
	JsonNodeString   JsonNodeType = 4
	JsonNodeLiteral  JsonNodeType = 5
	JsonNodeVariable JsonNodeType = 6
)

type JsonNode struct {
	Type JsonNodeType
	Val  interface{}
}

func newNumber(val interface{}) *JsonNode {
	return &JsonNode{
		Type: JsonNodeNumber,
		Val:  val,
	}
}

func newArray(val interface{}) *JsonNode {
	return &JsonNode{
		Type: JsonNodeArray,
		Val:  val,
	}
}

func newObject(val interface{}) *JsonNode {
	return &JsonNode{
		Type: JsonNodeObject,
		Val:  val,
	}
}

func newString(val interface{}) *JsonNode {
	return &JsonNode{
		Type: JsonNodeString,
		Val:  val,
	}
}

func newLiteral(val interface{}) *JsonNode {
	return &JsonNode{
		Type: JsonNodeLiteral,
		Val:  val,
	}
}

func newVariable(val interface{}) *JsonNode {
	return &JsonNode{
		Type: JsonNodeVariable,
		Val:  val,
	}
}

func (node *JsonNode) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}

	switch node.Type {
	case JsonNodeObject:
		buf.WriteByte(byte('{'))

		object := node.Val.(map[string]interface{})

		idx := 0
		for key, val := range object {
			data, err := val.(*JsonNode).MarshalJSON()
			if err != nil {
				return nil, err
			}

			buf.WriteByte(byte('"'))
			buf.WriteString(key)
			buf.WriteByte(byte('"'))
			buf.WriteByte(byte(':'))
			buf.Write(data)
			if idx < len(object)-1 {
				buf.WriteByte(byte(','))
			}
			idx++
		}

		buf.WriteByte(byte('}'))
	case JsonNodeArray:
		buf.WriteByte(byte('['))

		array := node.Val.([]interface{})

		for idx, val := range array {
			data, err := val.(*JsonNode).MarshalJSON()
			if err != nil {
				return nil, err
			}

			buf.Write(data)
			if idx < len(node.Val.([]interface{}))-1 {
				buf.WriteByte(byte(','))
			}
			idx++
		}

		buf.WriteByte(byte(']'))
	case JsonNodeLiteral:
		if node.Val == nil {
			buf.WriteString("null")
		} else {
			buf.WriteString(fmt.Sprintf(`%+v`, node.Val))
		}
	case JsonNodeNumber:
		buf.WriteString(fmt.Sprintf(`%+v`, node.Val))
	case JsonNodeString:
		buf.WriteString(strconv.Quote(node.Val.(string)))
	case JsonNodeVariable:
		buf.WriteString(fmt.Sprintf(`$%+v`, node.Val))
	default:
		return nil, fmt.Errorf("unknow type %d and val %+v", node.Type, node.Val)
	}

	return buf.Bytes(), nil
}

func maybeSetResult(l yyLexer, v interface{}) {
	l.(*lex).result = v
}
