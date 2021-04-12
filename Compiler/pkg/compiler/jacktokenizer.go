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

// TokenType is a string representing a type of lexical element that a token might be
type TokenType string

const (
	symbol     = TokenType("SYMBOL")
	intConst   = TokenType("INT_CONST")
	strConst   = TokenType("STRING_CONST")
	keyWord    = TokenType("KEYWORD")
	identifier = TokenType("IDENTIFIER")
)

// InvalidAccessError is returned when the caller attempted to access on type of lexical element when the tokenizer was at another
type InvalidAccessError struct {
	wasAttempted TokenType // Type the caller attempted to access
	wasValid     TokenType // Type the tokenizer was able to access
	wasValidVal  string    // The value of the valid token
}

func (e *InvalidAccessError) Error() string {
	return fmt.Sprintf("expected a token of type `%v` but found a token \"%v\" of type `%v` instead", e.wasAttempted, e.wasValidVal, e.wasValid)
}

func (jt *JackTokenizer) getValidVal() string {
	switch jt.tokenType {
	case symbol:
		return jt.symbol
	case intConst:
		return strconv.Itoa(int(jt.intVal))
	case strConst:
		return jt.stringVal
	case keyWord:
		return jt.keyWord
	case identifier:
		return jt.identifier
	}
	return ""
}

func (jt *JackTokenizer) newInvalidAccessError(wasAttempted TokenType) *InvalidAccessError {
	return &InvalidAccessError{
		wasAttempted: wasAttempted,
		wasValid:     jt.tokenType,
		wasValidVal:  jt.getValidVal(),
	}
}

// JackTokenizer is walks through a Jack program, setting its own internal state to reflect the lexical state of the program
type JackTokenizer struct {
	filename   string // The input file name
	stream     string // The entire file as a string
	streamlen  uint
	i          uint // Index of the character in the stream that is currently being analyzed
	tokenType  TokenType
	symbol     string // Becomes accessible when tokenType == "SYMBOL"
	intVal     uint   // Becomes accessible when tokenType == "INT_CONST"
	stringVal  string // Becomes accessible when tokenType == "STRING_CONST"
	keyWord    string // Becomes accessible when tokenType == "KEYWORD"
	identifier string // Becomes accessible when tokenType == "IDENTIFIER"
}

// NewJackTokenizer creates a tokenizer
func NewJackTokenizer(filename string) (*JackTokenizer, error) {
	data, err := ioutil.ReadFile(filename)
	text := string(data)

	return &JackTokenizer{
		filename:  filename,
		stream:    text,
		streamlen: uint(len(text)),
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

// Peek peaks ahead to the next character in the input stream.
// Because our input stream index comes to rest after the last character of our current
// token, we just call jt.curChar() with no arguments
func (jt *JackTokenizer) Peek() (byte, error) {
	if !(jt.i < jt.streamlen) {
		return 0, fmt.Errorf("unexpect EOF!")
	}
	return jt.curChar(), nil
}

// Advance -- each [non-recursive] call to jt.Advance() eats the next token in the jt.stream and and updates jt's
// internal state to reflect the token it just ate: It updates tokenType and sets whichever
// of jt.intVal, jt.stringVal, jt.keyWord, or jt.identifier corresponds.
func (jt *JackTokenizer) Advance() error {
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
		jt.tokenType = symbol
		jt.symbol = string(jt.curChar())
		jt.i++ // Eat the symbol
	} else if unicode.IsDigit(rune(jt.curChar())) {
		// If we're at an integer constant
		jt.tokenType = intConst
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
		jt.intVal = uint(val)
	} else if jt.curChar() == '"' {
		// If we're at a string constant
		jt.tokenType = strConst
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
			jt.tokenType = keyWord
			jt.keyWord = jt.stream[j:jt.i]
		} else {
			// Else this is an identifier
			jt.tokenType = identifier
			jt.identifier = jt.stream[j:jt.i]
		}
	}

	return nil
}

// curChar returns the current character being analyzed by the jt if no arguments are given.
// Can optionally be called with an offset argument to look ahead, i.e. jt.curChar(1) == jt.stream[jt.i + 1]
func (jt *JackTokenizer) curChar(offset ...uint) byte {
	if len(offset) == 0 {
		return jt.stream[jt.i]
	}
	return jt.stream[jt.i+offset[0]]
}

// HasMoreTokens returns true if the tokenizer has more tokens to scan
func (jt *JackTokenizer) HasMoreTokens() bool {
	return jt.i < jt.streamlen
}

// TokenType returns the TokenType of the current token
func (jt *JackTokenizer) TokenType() TokenType {
	return jt.tokenType
}

// Symbol returns the raw token if TokenType is symbol
func (jt *JackTokenizer) Symbol() (string, error) {
	if jt.tokenType != symbol {
		return "", jt.newInvalidAccessError(symbol)
	}
	return jt.symbol, nil
}

// IntVal returns the raw token if TokenType is intConst
func (jt *JackTokenizer) IntVal() (uint, error) {
	if jt.tokenType != intConst {
		return 0, jt.newInvalidAccessError(intConst)
	}
	return jt.intVal, nil

}

// StringVal returns the raw token if TokenType is strConst
func (jt *JackTokenizer) StringVal() (string, error) {
	if jt.tokenType != strConst {
		return "", jt.newInvalidAccessError(strConst)
	}
	return jt.stringVal, nil

}

// KeyWord returns the raw token if TokenType is keyWord
func (jt *JackTokenizer) KeyWord() (string, error) {
	if jt.tokenType != keyWord {
		return "", jt.newInvalidAccessError(keyWord)
	}
	return jt.keyWord, nil

}

// Identifier returns the raw token if TokenType is identifier
func (jt *JackTokenizer) Identifier() (string, error) {
	if jt.tokenType != identifier {
		return "", jt.newInvalidAccessError(identifier)
	}
	return jt.identifier, nil

}

// MarshalXML makes JackTokenizer a Marshaler (see https://golang.org/pkg/encoding/xml/#Marshaler)
func (jt *JackTokenizer) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	var err error
	var elemName string
	var data string
	var charData string

	switch jt.TokenType() {
	case symbol:
		elemName = "symbol"
		data, err = jt.Symbol()
		if err != nil {
			return err
		}
	case intConst:
		elemName = "integerConstant"
		intData, err := jt.IntVal()
		if err != nil {
			return err
		}
		data = strconv.Itoa(int(intData))
	case strConst:
		elemName = "stringConstant"
		data, err = jt.StringVal()
		if err != nil {
			return err
		}
	case keyWord:
		elemName = "keyword"
		data, err = jt.KeyWord()
		if err != nil {
			return err
		}
	case identifier:
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
