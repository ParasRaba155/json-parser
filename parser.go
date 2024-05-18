package main

import "fmt"

// JSONObj represents a valid json object in the Go world
type JSONObj struct {
	pairs []KeyValue
}

// KeyValue is the key and value of each field of json object
type KeyValue struct {
	Key   string
	Value jsonVal
}

// jsonVal interface must be satisfied by the each primitive value of key value pair of
// json object
type jsonVal interface {
	Value() any
}

// jsonString is the representation of string in go
type jsonString string

var _ jsonVal = jsonString("") // compile time check for interface impl

// Value to implement the jsonVal interface
func (s jsonString) Value() any {
	return s
}

// Parser for json inputs in byte
//
// NOTE: Use the `NewParser` to construct the Parser, do not use it directly
type Parser struct {
	lexer     *Lexer
	currToken Token
}

// parseError custom error for messaging the json parse errors throughout the parser
type parseError struct {
	Message string
	Pos     int
}

func (e *parseError) Error() string {
	return fmt.Sprintf("JSON parse error at position %d: %s", e.Pos, e.Message)
}

func newJSONParseError(msg string, pos int) error {
	return &parseError{Message: msg, Pos: pos}
}

// NewParser the constructor for the Parser,which initializes the Parser
func NewParser(input []byte) *Parser {
	lex := Lexer{input: input}
	return &Parser{lexer: &lex, currToken: lex.nextToken()}
}

// nextToken the helper function to get the next token from the lexer
// and it sets the currToken to the next token
func (p *Parser) nextToken() {
	p.currToken = p.lexer.nextToken()
}

// getPos the helper function to get the current token's position
func (p *Parser) getPos() int {
	return p.currToken.Pos
}

// Parse will parse the input with parsing rule on the token obtained from the
// lexer
//
// It will return the JSONObj if successfully parsed,otherwise will throw error
// of type jsonParseError
func (p *Parser) Parse() (JSONObj, error) {
	obj := JSONObj{}
	if p.currToken.Type != LEFT_CURLY_BRACES {
		return obj, newJSONParseError("Expected '{' at the start of the json", p.getPos())
	}
	p.nextToken()

	for p.currToken.Type != RIGHT_CURLY_BRACES {
		// try and parse the key
		if p.currToken.Type != STRING {
			return obj, newJSONParseError("Expected string for key", p.currToken.Pos)
		}
		key := p.currToken.Value[1 : len(p.currToken.Value)-1]
		p.nextToken()

		// After parsing the key, we must have a colon
		if p.currToken.Type != COLON {
			return obj, newJSONParseError("Expected ':'", p.currToken.Pos)
		}
		p.nextToken()

		// try and parse the value corresponding to current key
		value, err := p.parseValue()
		if err != nil {
			return obj, err
		}
		p.nextToken()

		obj.pairs = append(obj.pairs, KeyValue{Key: key, Value: value})

		// TODO: Deal with trailing comma
		if p.currToken.Type == COMMA {
			p.nextToken()
			continue
		}

		if p.currToken.Type != RIGHT_CURLY_BRACES {
			return obj, newJSONParseError("Expected } or ,", p.currToken.Pos)
		}
	}
	return obj, nil
}

// parseValue from the current token
func (p *Parser) parseValue() (jsonVal, error) {
	switch p.currToken.Type {
	case STRING:
		value := jsonString(p.currToken.Value[1 : len(p.currToken.Value)-1])
		return value, nil
	default:
		return nil, newJSONParseError("Expected string value", p.currToken.Pos)
	}
}
