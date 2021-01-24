package compiler

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"
)

// JackTokenizer is public interface for incrementing and reading the state of the tokenizer state machine
type JackTokenizer interface {
	Advance() error
	HasMoreTokens() bool
	TokenType() string
	Symbol() (string, error)
	IntVal() (int, error)
	StringVal() (string, error)
	KeyWord() (string, error)
	Identifier() (string, error)
	MarshalXML(e *xml.Encoder, start xml.StartElement) error // Used for testing
}

type jackTokenizer struct {
	filename   string // The input file name
	stream     string // The entire file as a string
	streamlen  int
	i          int // Index of the character in the stream that is currently being analyzed
	tokenType  string
	symbol     string // Becomes accessible when tokenType == "SYMBOL"
	intVal     int    // Becomes accessible when tokenType == "INT_CONST"
	stringVal  string // Becomes accessible when tokenType == "STRING_CONST"
	keyWord    string // Becomes accessible when tokenType == "KEYWORD"
	identifier string // Becomes accessible when tokenType == "IDENTIFIER"
}

// NewJackTokenizer creates a tokenizer
func NewJackTokenizer(filename string) (JackTokenizer, error) {
	data, err := ioutil.ReadFile(filename)
	text := string(data)

	return &jackTokenizer{
		filename:  filename,
		stream:    text,
		streamlen: len(text),
		i:         0,
	}, err
}

// Helper function to check if a string is a Jack language KeyWord
func isKeyWord(potentialKeyWord string) bool {
	switch potentialKeyWord {
	case
		"class",
		"method",
		"function",
		"constructor",
		"int",
		"boolean",
		"char",
		"void",
		"var",
		"static",
		"field",
		"let",
		"do",
		"if",
		"else",
		"while",
		"return",
		"true",
		"false",
		"null",
		"this":
		return true
	}
	return false
}

// Each [non-recursive] call to jt.Advance() eats the next token in the jt.stream and and updates jt's
// internal state to reflect the token it just ate: It updates tokenType and sets whichever
// of jt.intVal, jt.stringVal, jt.keyWord, or jt.identifier corresponds.
func (jt *jackTokenizer) Advance() error {
	// Needed to halt execution in case we reach EOF on a recursive call. Otherwise, callers should be
	// checking that jt.HasMoreTokens() before calling jt.Advance()
	if !(jt.i < jt.streamlen) {
		return nil
	}

	// Skip all whitespace characters
	if unicode.IsSpace(rune(jt.curChar())) {
		jt.i++
		jt.Advance()
		return nil
	}

	// Skip line comments
	if jt.curChar() == '/' && jt.curChar(1) == '/' {
		jt.i = jt.i + 2
		// Advance until either a newline or EOF
		for jt.HasMoreTokens() && jt.curChar() != '\n' {
			jt.i++
		}
		jt.i++ // Eat the '\n'
		jt.Advance()
		return nil
	}

	// Skip block comments
	if jt.curChar() == '/' && jt.curChar(1) == '*' {
		jt.i = jt.i + 2
		for jt.HasMoreTokens() && !(jt.curChar() == '*' && jt.curChar(1) == '/') {
			jt.i++
		}
		if !(jt.curChar() == '*' && jt.curChar(1) == '/') {
			// Hit EOF before block comment closed
			return errors.New("Encountered block comment open characters '/*' but didn't find subsequent block comment close characters '*/'")
		}
		jt.i = jt.i + 2 // Eat the closing "*/"
		jt.Advance()
		return nil
	}

	// Comments and whitespace have been skipped, now determine what type of lexical element we're analyzing
	if strings.Contains("{}()[].,;+-*/&|<>=~", string(jt.curChar())) {
		// If we're at a single-character symbol
		jt.tokenType = "SYMBOL"
		jt.symbol = string(jt.curChar())
		jt.i++ // Eat the symbol
	} else if unicode.IsDigit(rune(jt.curChar())) {
		// If we're at an integer constant
		jt.tokenType = "INT_CONST"
		j := jt.i // Save our current index
		jt.i++    // Eat the current char

		for unicode.IsDigit(rune(jt.curChar())) {
			// Eat the subsequent numeric values to capture the full constant
			jt.i++
		}
		// Attempt to convert the captured constant to an integer
		val, err := strconv.Atoi(jt.stream[j:jt.i])
		if err != nil {
			return fmt.Errorf("Compiler bug: attempted to parse %v as an integer", jt.stream[j:jt.i])
		}
		if val > 32767 {
			return fmt.Errorf("Integer constants must be a decimal number in the range of 0 .. 32767. Got %v", val)
		}
		jt.intVal = val
	} else if jt.curChar() == '"' {
		// If we're at a string constant
		jt.tokenType = "STRING_CONST"
		jt.i++    // Eat the '"'
		j := jt.i // Save our current index
		for jt.HasMoreTokens() && jt.curChar() != '"' {
			// Eat the rest of the string constant
			if jt.curChar() == '\n' {
				return errors.New("String constants cannot contain newline characters")
			}
			jt.i++
		}
		if !jt.HasMoreTokens() {
			return errors.New("Hit EOF before string constant terminated")
		}
		if jt.i == j {
			// TODO: Maybe we should support empty strings?
			return errors.New("Encountered unsupported empty string constant")
		}
		jt.stringVal = jt.stream[j:jt.i]
		jt.i++ // Eat the closing '"'
	} else {
		// We're at either an identifier or a keyword or an invalid character
		if !(unicode.IsLetter(rune(jt.curChar())) || jt.curChar() == '_') {
			// If we hit an invalid first character for an indentifier or keyword, return error
			return fmt.Errorf("Encountered invalid character: %v", jt.curChar())
		}
		j := jt.i // Save current index
		for jt.HasMoreTokens() && (unicode.IsLetter(rune(jt.curChar())) || unicode.IsDigit(rune(jt.curChar())) || jt.curChar() == '_') {
			// Eat the rest of the characters in this identifier or keyword
			jt.i++
		}

		if isKeyWord(jt.stream[j:jt.i]) {
			// If this is a keyword
			jt.tokenType = "KEYWORD"
			jt.keyWord = jt.stream[j:jt.i]
		} else {
			// Else this is an identifier
			jt.tokenType = "IDENTIFIER"
			jt.identifier = jt.stream[j:jt.i]
		}
	}

	return nil
}

// curChar returns the current character being analyzed by the jt if no arguments are given.
// Can optionally be called with an offset argument to look ahead, i.e. jt.curChar(1) == jt.stream[jt.i + 1]
func (jt *jackTokenizer) curChar(offset ...int) byte {
	if len(offset) == 0 {
		return jt.stream[jt.i]
	}
	return jt.stream[jt.i+offset[0]]
}

func (jt *jackTokenizer) HasMoreTokens() bool {
	return jt.i < jt.streamlen
}

func (jt *jackTokenizer) TokenType() string {
	return jt.tokenType
}

func (jt *jackTokenizer) Symbol() (string, error) {
	if jt.tokenType != "SYMBOL" {
		return "", fmt.Errorf("Attempted to access jt.symbol but jt.tokenType was not \"SYMBOL\" (it was \"%v\")", jt.tokenType)
	}
	return jt.symbol, nil
}

func (jt *jackTokenizer) IntVal() (int, error) {
	if jt.tokenType != "INT_CONST" {
		return 0, fmt.Errorf("Attempted to access jt.intVal but jt.tokenType was not \"INT_CONST\" (it was \"%v\")", jt.tokenType)
	}
	return jt.intVal, nil

}

func (jt *jackTokenizer) StringVal() (string, error) {
	if jt.tokenType != "STRING_CONST" {
		return "", fmt.Errorf("Attempted to access jt.stringVal but jt.tokenType was not \"STRING_CONST\" (it was \"%v\")", jt.tokenType)
	}
	return jt.stringVal, nil

}

func (jt *jackTokenizer) KeyWord() (string, error) {
	if jt.tokenType != "KEYWORD" {
		return "", fmt.Errorf("Attempted to access jt.keyWord but jt.tokenType was not \"KEYWORD\" (it was \"%v\")", jt.tokenType)
	}
	return jt.keyWord, nil

}

func (jt *jackTokenizer) Identifier() (string, error) {
	if jt.tokenType != "IDENTIFIER" {
		return "", fmt.Errorf("Attempted to access jt.identifier but jt.tokenType was not \"IDENTIFIER\" (it was \"%v\")", jt.tokenType)
	}
	return jt.identifier, nil

}

// https://golang.org/pkg/encoding/xml/#Marshaler
func (jt *jackTokenizer) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	var err error
	var elemName string
	var data string
	var charData string

	switch jt.TokenType() {
	case "SYMBOL":
		elemName = "symbol"
		data, err = jt.Symbol()
		if err != nil {
			return err
		}
	case "INT_CONST":
		elemName = "integerConstant"
		intData, err := jt.IntVal()
		if err != nil {
			return err
		}
		data = strconv.Itoa(intData)
	case "STRING_CONST":
		elemName = "stringConstant"
		data, err = jt.StringVal()
		if err != nil {
			return err
		}
	case "KEYWORD":
		elemName = "keyword"
		data, err = jt.KeyWord()
		if err != nil {
			return err
		}
	case "IDENTIFIER":
		elemName = "identifier"
		data, err = jt.Identifier()
		if err != nil {
			return err
		}
	}

	charData = fmt.Sprintf(" %v ", data)

	err = e.EncodeToken(
		xml.StartElement{Name: xml.Name{Space: "", Local: elemName}, Attr: []xml.Attr{}})
	err = e.EncodeToken(xml.CharData(charData))
	err = e.EncodeToken(xml.EndElement{Name: xml.Name{Space: "", Local: elemName}})

	return err
}
