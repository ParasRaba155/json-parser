package main

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

var (
	trueByte  = []byte("true")
	falseByte = []byte("false")
	nullByte  = []byte("null")
)

// tokenType represents the different JSON tokens
//
//go:generate stringer -type=tokenType
type tokenType int

const (
	INVALID tokenType = iota

	// curly braces for object
	LEFT_CURLY_BRACES  tokenType = 1
	RIGHT_CURLY_BRACES tokenType = 2

	// square brackets for arrays
	LEFT_SQUARE_BRACKET  tokenType = 3
	RIGHT_SQUARE_BRACKET tokenType = 4

	// separator tokens
	COLON tokenType = 5
	COMMA tokenType = 7

	// primitive types tokens
	STRING       tokenType = 8
	FLOAT_NUMBER tokenType = 9
	INT_NUMBER   tokenType = 10
	BOOLEAN      tokenType = 11
	NULL         tokenType = 12

	// special End of file token
	EOF tokenType = 13
)

// Token containing the value and type of the token, and current pos in the
// input
type Token struct {
	Value string    // Value of the token
	Type  tokenType // The type of the token
	Pos   int       // Position of the token
}

// Lexer will read the input and breaks it into tokens
// It will shift from left to right, keeping track of characters
// and move its pos accordingly
type Lexer struct {
	input []byte
	pos   int
}

// nextChar will read the next character from the input, return it
// will return 0 if we have shifted through all the chars in input
// will shift the position to the right of the current char
func (l *Lexer) nextChar() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	return ch
}

// peekChar will read the next character from the input, return it
// will return 0 if we have shifted through all the chars in input
// it will not move the cursor position
func (l *Lexer) peekChar() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := l.input[l.pos]
	return ch
}

// nextToken will read the chars from the input, and create appropriate
// tokens from the input
func (l *Lexer) nextToken() Token {
	for {
		ch := l.nextChar()
		switch ch {
		case '{':
			return Token{Type: LEFT_CURLY_BRACES, Pos: l.pos - 1}
		case '}':
			return Token{Type: RIGHT_CURLY_BRACES, Pos: l.pos - 1}
		case '[':
			return Token{Type: LEFT_SQUARE_BRACKET, Pos: l.pos - 1}
		case ']':
			return Token{Type: RIGHT_SQUARE_BRACKET, Pos: l.pos - 1}
		case ':':
			return Token{Type: COLON, Pos: l.pos - 1}
		case ',':
			return Token{Type: COMMA, Pos: l.pos - 1}
		case '"':
			return l.readString()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			return l.readNumber()
		case 't', 'f':
			return l.readBoolean()
		case 'n':
			return l.readNull()
		case 0:
			// 0 byte is represented as EOF
			return Token{Type: EOF, Pos: l.pos}
		default:
			if unicode.IsSpace(rune(ch)) {
				continue
			}
			return Token{Type: INVALID, Pos: l.pos, Value: fmt.Sprintf("Unexpected char: %c", rune(ch))}
		}
	}
}

// readString will try to read the string from the current position
// Should be called only when the current char is found to be the \" char
//
// NOTE: Will change the position of the input pointer, as it uses 'nextChar`
func (l *Lexer) readString() Token {
	start := l.pos - 1
	// READ till the end of string or till we encounter EOF
	for {
		ch := l.nextChar()
		if ch == '"' || ch == 0 {
			break
		}
	}
	// if the last char is not "\"" then the string is unterminated, handle it
	if l.input[l.pos-1] != '"' {
		return Token{Type: INVALID, Pos: start, Value: "Unterminated string"}
	}
	return Token{Type: STRING, Value: string(l.input[start:l.pos]), Pos: start}
}

func (l *Lexer) readNumber() Token {
	start := l.pos - 1
	numType := INT_NUMBER
	// read till the end of number
	for {
		ch := l.peekChar()

		// change the number type
		if ch == '.' {
			numType = FLOAT_NUMBER
		}

		// check for the end of the line or end of file or end of object
		if ch == ',' || ch == '}' || ch == 0 || ch == '\n' {
			break
		}
		l.nextChar()
	}

	numStr := string(l.input[start:l.pos])

	// try to parse the number into float, if unsuccessful that means
	// there is some error
	_, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return Token{Type: INVALID, Pos: start, Value: "Invalid number"}
	}
	return Token{Type: numType, Value: numStr, Pos: start}
}

func (l *Lexer) readBoolean() Token {
	start := l.pos - 1
	// read till the end
	for {
		ch := l.peekChar()
		// check for the end of the line or end of file or end of object
		if ch == ',' || ch == '}' || ch == 0 || ch == '\n' {
			break
		}
		l.nextChar()
	}
	boolByte := l.input[start:l.pos]
	// slog.Info("FUCK",
	// 	slog.String("boolByte", string(boolByte)),
	// 	slog.Bool("is true", bytes.Equal(boolByte, trueByte)),
	// 	slog.Bool("is false", bytes.Equal(boolByte, falseByte)),
	// )
	if !bytes.Equal(boolByte, trueByte) && !bytes.Equal(boolByte, falseByte) {
		return Token{Type: INVALID, Pos: start, Value: "Expected boolean"}
	}
	return Token{Type: BOOLEAN, Value: string(boolByte), Pos: start}
}

func (l *Lexer) readNull() Token {
	start := l.pos - 1
	// read till the end
	for {
		ch := l.peekChar()
		// check for the end of the line or end of file or end of object
		if ch == ',' || ch == '}' || ch == 0 || ch == '\n' {
			break
		}
		l.nextChar()
	}
	found := l.input[start:l.pos]
	if !bytes.Equal(found, nullByte) {
		return Token{Type: INVALID, Pos: start, Value: "Expected null"}
	}
	return Token{Type: NULL, Value: string(found), Pos: start}
}
