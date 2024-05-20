package main

import (
	"fmt"
	"strconv"
)

// JSONObj represents a valid json object in the Go world
type JSONObj struct {
	pairs []KeyValue
}

// Value implements jsonVal.
func (j JSONObj) Value() any {
	return j.pairs
}

var _ jsonVal = JSONObj{} // compile time check

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

// jsonFloat is the representation of floating numbers in go
type jsonFloat float64

var _ jsonVal = jsonFloat(0.0) // compile time check for interface impl

// Value implements jsonVal.
func (j jsonFloat) Value() any {
	return j
}

// jsonInt is the representation of int numbers in go
type jsonInt int

var _ jsonVal = jsonInt(0)

// Value implements jsonVal.
func (j jsonInt) Value() any {
	return j
}

// jsonBool is the representation of boolean in go
type jsonBool bool

var _ jsonVal = jsonBool(false)

// Value implements jsonVal.
func (j jsonBool) Value() any {
	return j
}

// jsonArray representation of array in go
type jsonArray []jsonVal

// Value implements jsonVal.
func (j jsonArray) Value() any {
	return j
}

var _ jsonVal = jsonArray{} // compile time check

// Parser for json inputs in byte
//
// NOTE: Use the `NewParser` to construct the Parser, do not use it directly
type Parser struct {
	lexer            *Lexer
	currToken        Token
	wasPreviousComma bool // To keep track if we have already found comma
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
			return obj, newJSONParseError("Expected string for key", p.getPos())
		}

		// currently to handle the trailing commas not supported, we are keeping
		// track of these flag `wasPreviousComma` which should be set to true
		// on encountering the comma and otherwise to false
		p.wasPreviousComma = false
		key := p.currToken.Value[1 : len(p.currToken.Value)-1]
		p.nextToken()

		// After parsing the key, we must have a colon
		if p.currToken.Type != COLON {
			return obj, newJSONParseError("Expected ':'", p.getPos())
		}
		p.wasPreviousComma = false
		p.nextToken()

		// try and parse the value corresponding to current key
		value, err := p.parseValue()
		if err != nil {
			return obj, err
		}

		obj.pairs = append(obj.pairs, KeyValue{Key: key, Value: value})
		p.nextToken()

		if p.currToken.Type == COMMA {
			p.nextToken()
			p.wasPreviousComma = true
			continue
		}

		if p.currToken.Type == EOF {
			p.wasPreviousComma = false
			break
		}

		if p.currToken.Type != RIGHT_CURLY_BRACES {
			return obj, newJSONParseError("Expected } or ,", p.getPos())
		}
	}
	if p.wasPreviousComma {
		return obj, newJSONParseError("trailing commas are not supported", p.getPos())
	}
	return obj, nil
}

// parseValue from the current token
func (p *Parser) parseValue() (jsonVal, error) {
	switch p.currToken.Type {
	case STRING:
		value := jsonString(p.currToken.Value[1 : len(p.currToken.Value)-1])
		return value, nil
	case INT_NUMBER:
		num, err := strconv.Atoi(p.currToken.Value)
		if err != nil {
			return nil, newJSONParseError("expected a number", p.getPos())
		}
		value := jsonInt(num)
		return value, nil
	case FLOAT_NUMBER:
		num, err := strconv.ParseFloat(p.currToken.Value, 64)
		if err != nil {
			return nil, newJSONParseError("Expected a number", p.getPos())
		}
		value := jsonFloat(num)
		return value, nil
	case BOOLEAN:
		bool, err := strconv.ParseBool(p.currToken.Value)
		value := jsonBool(bool)
		if err != nil {
			return value, newJSONParseError("Expected a boolean", p.getPos())
		}
		return value, nil
	case NULL:
		if p.currToken.Value != "null" {
			return nil, newJSONParseError("Expected a null value", p.getPos())
		}
		return nil, nil
	case LEFT_SQUARE_BRACKET:
		return p.parseArray()
	case LEFT_CURLY_BRACES:
		return p.Parse()
	default:
		return nil, newJSONParseError("Expected value", p.getPos())
	}
}

// parseArray the helper recursive function to call on encountering an array to be
// parsed
func (p *Parser) parseArray() (jsonVal, error) {
	var arr jsonArray

	if p.currToken.Type != LEFT_SQUARE_BRACKET {
		return nil, newJSONParseError("Expected '[' at the start of array", p.getPos())
	}
	p.nextToken()

	for p.currToken.Type != RIGHT_SQUARE_BRACKET {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)
		p.nextToken()

		// TODO: Handle trailing commas in the arrays
		if p.currToken.Type == COMMA {
			p.nextToken()
			continue
		}
		if p.currToken.Type == RIGHT_SQUARE_BRACKET {
			break
		}
		return nil, newJSONParseError("unterminated array", p.getPos())
	}
	return arr, nil
}
