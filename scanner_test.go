package variable_json

import (
	"reflect"
	"testing"
)

type token struct {
	typ int
	val interface{}
}

func TestLex(t *testing.T) {
	testcases := []struct {
		input  string
		output []token
	}{{
		input:  "",
		output: nil,
	}, {
		input:  " \t\n",
		output: nil,
	}, {
		input: `"abcd"`,
		output: []token{{
			typ: String,
			val: `abcd`,
		}},
	}, {
		input: `"abcd\"\\\/\b\f\n\r\tef"`,
		output: []token{{
			typ: String,
			val: "abcd\"\\/\b\f\n\r\tef",
		}},
	}, {
		input: "1234.56e-20",
		output: []token{{
			typ: Number,
			val: float64(1234.56e-20),
		}},
	}, {
		input: "true false null",
		output: []token{{
			typ: Literal,
			val: true,
		}, {
			typ: Literal,
			val: false,
		}, {
			typ: Literal,
			val: nil,
		}},
	}, {
		input: `{ "a" "\"上xx学校\"" 1 true`,
		output: []token{{
			typ: '{',
		}, {
			typ: String,
			val: "a",
		}, {
			typ: String,
			val: "\"上xx学校\"",
		}, {
			typ: Number,
			val: float64(1),
		}, {
			typ: Literal,
			val: true,
		}},
	}, {
		input: `"aa\`,
		output: []token{{
			typ: LexError,
		}},
	}, {
		input: `"aa`,
		output: []token{{
			typ: LexError,
		}},
	}, {
		input: "12ee",
		output: []token{{
			typ: LexError,
		}},
	}, {
		input: "invalid",
		output: []token{{
			typ: LexError,
		}},
	}}

	for _, tc := range testcases {
		var got []token
		l := newLex([]byte(tc.input))
		var lval yySymType
		for {
			typ := l.Lex(&lval)
			if typ == 0 {
				break
			}
			got = append(got, token{typ: typ, val: lval.val})
		}
		if !reflect.DeepEqual(got, tc.output) {
			t.Errorf("%v: %v, want %v", tc.input, got, tc.output)
		}
	}
}
