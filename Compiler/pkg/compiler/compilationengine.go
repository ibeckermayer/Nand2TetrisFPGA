package compiler

import (
	"encoding/xml"
	"fmt"
	"os"
)

// SyntaxError is a syntax error
type SyntaxError error

// CompilationEngine effects the actual compilation output.
// Gets its input from a JackTokenizer and emits its parsed structure into an output file/stream.
type CompilationEngine struct {
	JackFilePath string         // The name of the .jack input file to be compiled.
	jt           *JackTokenizer // A tokenizer set up to tokenize the file we want to compile
	outputFile   *os.File       // The output file
	xmlEnc       *xml.Encoder   // The xml encoder for testing
}

// Run runs the compiler on ce.JackFilePath
func (ce *CompilationEngine) Run() error {
	// Initialize the ce's corresponding JackTokenizer
	jt, err := NewJackTokenizer(ce.JackFilePath)
	if err != nil {
		return err
	}
	ce.jt = jt

	// Create the output file
	outputFile, err := os.Create(fmt.Sprintf("%v_out.xml", ce.JackFilePath[0:len(ce.JackFilePath)-len(".jack")]))
	if err != nil {
		return err
	}
	ce.outputFile = outputFile
	defer ce.outputFile.Close()

	// Create xml encoder
	ce.xmlEnc = xml.NewEncoder(outputFile)
	ce.xmlEnc.Indent("", "  ")
	defer ce.xmlEnc.Flush()

	// Advance to eat the first token and call compileClass, which will recursively compile the entire file
	ce.jt.Advance()
	return ce.compileClass()
}

// Shorthand for opening an xml tag
func (ce *CompilationEngine) openXMLTag(tag string) error {
	return ce.xmlEnc.EncodeToken(xml.StartElement{Name: xml.Name{Space: "", Local: tag}, Attr: []xml.Attr{}})
}

// Shorthand for writing raw data in an XML tag, prefixed and postfixed with a space character
func (ce *CompilationEngine) writeXMLData(data string) error {
	return ce.xmlEnc.EncodeToken(xml.CharData(fmt.Sprintf(" %v ", data)))
}

// Shorthand for closing an xml tag
func (ce *CompilationEngine) closeXMLTag(tag string) error {
	return ce.xmlEnc.EncodeToken(xml.EndElement{Name: xml.Name{Space: "", Local: tag}})
}

// Shorthand for calling jt.MarshalXML to take advantage of already existing xml encoding logic
func (ce *CompilationEngine) marshaljt() error {
	return ce.jt.MarshalXML(ce.xmlEnc, xml.StartElement{Name: xml.Name{Space: "not used", Local: "not used"}, Attr: []xml.Attr{}})
}

// Assumes the first token has already been consumed (should be 'class').
// 'class' className '{' classVarDec* subroutineDec* '}'
func (ce *CompilationEngine) compileClass() error {
	ce.openXMLTag("class")        // <class>
	defer ce.closeXMLTag("class") // defer closing </class>

	// Check that first token is "class"
	if kw, err := ce.jt.KeyWord(); kw != "class" {
		if err != nil {
			return err
		}
		return SyntaxError(fmt.Errorf("expected keyword \"class\""))
	}

	// Found class, write keyword
	ce.marshaljt() // <keyword> class </keyword>

	ce.jt.Advance()

	// Check that next token is an identifier
	if _, err := ce.jt.Identifier(); err != nil {
		return err
	}

	ce.marshaljt() // <identifier> ClassName </identifier>

	ce.jt.Advance()

	// Check that next token is "{"
	if sym, err := ce.jt.Symbol(); sym != "{" {
		if err != nil {
			return err
		}
		return SyntaxError(fmt.Errorf("expected the symbol \"{\""))
	}

	ce.marshaljt() // <symbol> { </symbol>

	// Loop through and compile all of the classVarDecs
	if err := ce.jt.Advance(); err != nil {
		return err
	}
	for kw, err := ce.jt.KeyWord(); kw == "static" || kw == "field"; kw, err = ce.jt.KeyWord() {
		if err != nil {
			return err
		}

		ce.compileClassVarDec()

		if err := ce.jt.Advance(); err != nil {
			return err
		}
	}

	// Loop through and compile all of the subroutines.
	// The previous loop will have called Advance() and then hit a non static/field
	for kw, err := ce.jt.KeyWord(); kw == "constructor" || kw == "function" || kw == "method"; kw, err = ce.jt.KeyWord() {
		if err != nil {
			return err
		}

		ce.compileSubroutine()

		if err := ce.jt.Advance(); err != nil {
			return err
		}

	}

	// Check that next token is "}"
	// The previous loop should have called Advance() for this symbol
	if sym, err := ce.jt.Symbol(); sym != "}" {
		if err != nil {
			return err
		}
		return SyntaxError(fmt.Errorf("expected the symbol \"}\" but got \"%v\" instead", sym))
	}

	ce.marshaljt() // <symbol> } </symbol>

	return nil
}

// ('static' | 'field') type varName (',' varName)* ';'
func (ce *CompilationEngine) compileClassVarDec() error {
	ce.openXMLTag("classVarDec")        // <classVarDec>
	defer ce.closeXMLTag("classVarDec") //</classVarDec>

	// Confirm that first token is "static" or "field"
	if kw, err := ce.jt.KeyWord(); !(kw == "static" || kw == "field") {
		if err != nil {
			return err
		}
		return SyntaxError(fmt.Errorf("expected %v \"static\" or \"field\"", keyWord))
	}

	// found "static" or "field"
	ce.marshaljt() // <keyword> * </keyword>

	ce.jt.Advance()

	// Check for a type: 'int' | 'char' | 'boolean' | className
	switch ce.jt.TokenType() {
	case keyWord:
		kw, _ := ce.jt.KeyWord()
		if !(kw == "int" || kw == "char" || kw == "boolean") {
			return SyntaxError(fmt.Errorf("expected a type (%v \"int\" or \"char\" or \"boolean\" or %v className)", keyWord, identifier))
		}
		// Found "int" or "char" or "boolean"
		ce.marshaljt() // <keyword> * </keyword>
	case identifier:
		ce.marshaljt() // <identifier> className </identifier>
	default:
		return SyntaxError(fmt.Errorf("expected a type (%v \"int\" or \"char\" or \"boolean\" or %v className)", keyWord, identifier))
	}

	// Check for varName
	ce.jt.Advance()
	_, err := ce.jt.Identifier()
	if err != nil {
		return err
	}
	ce.marshaljt() // <identifier> varName </identifier>

	// Check for a comma separated list of more varNames
	ce.jt.Advance()
	for sym, err := ce.jt.Symbol(); sym == ","; sym, err = ce.jt.Symbol() {
		if err != nil {
			return err
		}
		ce.marshaljt() // <keyword> , </keyword>

		// Check for varName
		ce.jt.Advance()
		_, err := ce.jt.Identifier()
		if err != nil {
			return err
		}
		ce.marshaljt() // <identifier> varName </identifier>
		ce.jt.Advance()
	}

	// Should wind up at a ";"
	if sym, err := ce.jt.Symbol(); sym != ";" {
		if err != nil {
			return err
		}
		return SyntaxError(fmt.Errorf("expected the %v \";\"", symbol))
	}

	ce.marshaljt() // <symbol> ; </symbol>

	return nil
}

func (ce *CompilationEngine) compileSubroutine() {
	// panic("not implemented") // TODO: Implement
	return
}

func (ce *CompilationEngine) compileParameterList() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) compileVarDec() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) compileStatements() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) compileDo() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) compileLet() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) compileWhile() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) compileReturn() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) compileIf() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) compileExpression() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) compileTerm() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) compileExpressionList() {
	panic("not implemented") // TODO: Implement
}
