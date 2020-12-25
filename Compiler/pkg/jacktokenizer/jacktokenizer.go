package jacktokenizer

import (
	"io/ioutil"
)

// JackTokenizer is public interface for incrementing and reading the state of the tokenizer state machine
type JackTokenizer interface {
	Advance()
	HasMoreTokens() bool
	TokenType() (string, error)
	Symbol() (string, error)
	IntVal() (int, error)
	StringVal() (string, error)
	KeyWord() (string, error)
	Identifier() (string, error)
}

type jackTokenizer struct {
	filename   string
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

func (j *jackTokenizer) Advance() {
	panic("not implemented") // TODO: Implement
}

func (j *jackTokenizer) HasMoreTokens() bool {
	panic("not implemented") // TODO: Implement
}

func (j *jackTokenizer) TokenType() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (j *jackTokenizer) Symbol() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (j *jackTokenizer) IntVal() (int, error) {
	panic("not implemented") // TODO: Implement
}

func (j *jackTokenizer) StringVal() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (j *jackTokenizer) KeyWord() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (j *jackTokenizer) Identifier() (string, error) {
	panic("not implemented") // TODO: Implement
}
