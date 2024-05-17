package main

import (
	"fmt"
	"unicode"
)

// tokenType represents the different JSON token types
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
	STRING  tokenType = 8
	NUMBER  tokenType = 9
	BOOLEAN tokenType = 10
	NULL    tokenType = 11

	// special End of file token
	EOF tokenType = 12
)

type Token struct {
	Value string    // Value of the token // TODO: INCORPORATE DIFFERENT TYPES OF TOKENS
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
