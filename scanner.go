package variable_json

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type lex struct {
	input  []byte
	pos    int
	result interface{}
	err    error
}

func newLex(input []byte) *lex {
	return &lex{
		input: input,
	}
}

func (l *lex) Lex(lval *yySymType) int {
	return l.scanNormal(lval)
}

func (l *lex) scanNormal(lval *yySymType) int {
	for b := l.next(); b != 0; b = l.next() {
		switch {
		case unicode.IsSpace(rune(b)):
			continue
		case b == '"':
			return l.scanString(lval)
		case unicode.IsDigit(rune(b)) || b == '+' || b == '-':
			l.backup()
			return l.scanNum(lval)
		case unicode.IsLetter(rune(b)):
			l.backup()
			return l.scanLiteral(lval)
		case b == '$':
			return l.scanVariable(lval)
		default:
			return int(b)
		}
	}
	return 0
}

func isVariableChar(b rune) bool {
	return ('a' <= b && b <= 'z') ||
		('A' <= b && b <= 'Z') ||
		('0' <= b && b <= '9') ||
		b == '_' ||
		b == '-'
}

func isValidVariable(val string) bool {
	if len(val) <= 0 {
		return false
	}

	b := val[0]

	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z')
}

func (l *lex) scanVariable(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.next()
		switch {
		case isVariableChar(rune(b)):
			buf.WriteByte(b)
		default:
			l.backup()

			if !isValidVariable(buf.String()) {
				return LexError
			}

			lval.val = buf.String()
			return Variable
		}
	}
}

var escape = map[byte]byte{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

func (l *lex) scanString(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for b := l.next(); b != 0; b = l.next() {
		switch b {
		case '\\':
			b1 := l.next()
			if b1 == 'u' {
				buf.WriteByte(b)
				buf.WriteByte(b1)
				if !l.writeUnicode(buf) {
					return LexError
				}
			} else {
				b2 := escape[b1]
				if b2 == 0 {
					return LexError
				}
				buf.WriteByte(b2)
			}
		case '"':
			lval.val = buf.String()
			return String
		default:
			buf.WriteByte(b)
		}
	}
	return LexError
}

func (l *lex) writeUnicode(buf *bytes.Buffer) bool {
	for i := 0; i < 4; i++ {
		b := l.next()
		if ('0' <= b && b <= '9') || ('a' <= b &&b <= 'z') {
			buf.WriteByte(b)
		} else {
			return false
		}
	}

	return true
}

func (l *lex) scanNum(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.next()
		switch {
		case unicode.IsDigit(rune(b)):
			buf.WriteByte(b)
		case strings.IndexByte(".+-eE", b) != -1:
			buf.WriteByte(b)
		default:
			l.backup()
			val, err := strconv.ParseFloat(buf.String(), 64)
			if err != nil {
				return LexError
			}
			lval.val = val
			return Number
		}
	}
}

var literal = map[string]interface{}{
	"true":  true,
	"false": false,
	"null":  nil,
}

func (l *lex) scanLiteral(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		b := l.next()
		switch {
		case unicode.IsLetter(rune(b)):
			buf.WriteByte(b)
		default:
			l.backup()
			val, ok := literal[buf.String()]
			if !ok {
				return LexError
			}
			lval.val = val
			return Literal
		}
	}
}

func (l *lex) backup() {
	if l.pos == -1 {
		return
	}
	l.pos--
}

func (l *lex) next() byte {
	if l.pos >= len(l.input) || l.pos == -1 {
		l.pos = -1
		return 0
	}
	l.pos++
	return l.input[l.pos-1]
}

func (l *lex) Error(s string) {
	start := l.pos - 100
	end := l.pos + 100

	if start < 0 {
		start = 0
	}

	if end > len(l.input) {
		end = len(l.input)
	}

	l.err = errors.New(string(l.input[start:end]) + " - " + s)
}
