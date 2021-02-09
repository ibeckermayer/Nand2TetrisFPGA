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
	if err := ce.jt.Advance(); err != nil {
		return err
	}
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
			return SyntaxError(err)
		}
		return SyntaxError(fmt.Errorf("expected keyword \"class\""))
	}

	// Found class, write keyword
	ce.marshaljt() // <keyword> class </keyword>

	if err := ce.jt.Advance(); err != nil {
		return err
	}

	// Check that next token is an identifier
	if _, err := ce.jt.Identifier(); err != nil {
		return err
	}

	ce.marshaljt() // <identifier> ClassName </identifier>

	if err := ce.jt.Advance(); err != nil {
		return err
	}

	// Check that next token is "{"
	if err := ce.compileSymbol("{"); err != nil {
		return SyntaxError(err)
	}

	if err := ce.jt.Advance(); err != nil {
		return err
	}

	// Loop through and compile all of the classVarDecs
	for kw, err := ce.jt.KeyWord(); kw == "static" || kw == "field"; kw, err = ce.jt.KeyWord() {
		if err != nil {
			return SyntaxError(err)
		}

		if err := ce.compileClassVarDec(); err != nil {
			return SyntaxError(err)
		}

		if err := ce.jt.Advance(); err != nil {
			return err
		}
	}

	// Loop through and compile all of the subroutines.
	// The previous loop will have called Advance() and then hit a non static/field
	for kw, err := ce.jt.KeyWord(); kw == "constructor" || kw == "function" || kw == "method"; kw, err = ce.jt.KeyWord() {
		if err != nil {
			return SyntaxError(err)
		}

		if err := ce.compileSubroutine(); err != nil {
			return SyntaxError(err)
		}

		if err := ce.jt.Advance(); err != nil {
			return err
		}

	}

	// Check that next token is "}"
	// The previous loop should have called Advance() for this symbol
	if err := ce.compileSymbol("}"); err != nil {
		return SyntaxError(err)
	}

	return nil
}

// 'void | 'int' | 'char' | 'boolean' | className
func (ce *CompilationEngine) compileVoidOrType() error {
	if kw, err := ce.jt.KeyWord(); kw == "void" && err == nil {
		ce.marshaljt() // <keyword> void </keyword>
		return nil
	}
	err := ce.compileType()
	if err != nil {
		// see errString in compileType()
		return SyntaxError(fmt.Errorf("%v or \"void\"", err.Error()))
	}
	return err
}

// 'int' | 'char' | 'boolean' | className
func (ce *CompilationEngine) compileType() error {
	var errString string = "expected a type: %v className or %v \"int\" or \"char\" or \"boolean\""

	switch ce.jt.TokenType() {
	case keyWord:
		kw, _ := ce.jt.KeyWord()
		if !(kw == "int" || kw == "char" || kw == "boolean") {
			return SyntaxError(fmt.Errorf(errString, identifier, keyWord))
		}
		// Found "int" or "char" or "boolean"
		ce.marshaljt() // <keyword> * </keyword>
		return nil
	case identifier:
		ce.marshaljt() // <identifier> className </identifier>
		return nil
	default:
		return SyntaxError(fmt.Errorf(errString, identifier, keyWord))
	}
}

// ('static' | 'field') type varName (',' varName)* ';'
func (ce *CompilationEngine) compileClassVarDec() error {
	ce.openXMLTag("classVarDec")        // <classVarDec>
	defer ce.closeXMLTag("classVarDec") //</classVarDec>

	// Confirm that first token is "static" or "field"
	if kw, err := ce.jt.KeyWord(); !(kw == "static" || kw == "field") {
		if err != nil {
			return SyntaxError(err)
		}
		return SyntaxError(fmt.Errorf("expected %v \"static\" or \"field\"", keyWord))
	}

	// found "static" or "field"
	ce.marshaljt() // <keyword> * </keyword>

	if err := ce.jt.Advance(); err != nil {
		return err
	}

	// try to compile a type
	if err := ce.compileType(); err != nil {
		return err
	}

	// Check for varName
	if err := ce.jt.Advance(); err != nil {
		return err
	}
	_, err := ce.jt.Identifier()
	if err != nil {
		return SyntaxError(err)
	}
	ce.marshaljt() // <identifier> varName </identifier>

	// Check for a comma separated list of more varNames
	if err := ce.jt.Advance(); err != nil {
		return err
	}
	for sym, err := ce.jt.Symbol(); sym == ","; sym, err = ce.jt.Symbol() {
		if err != nil {
			return SyntaxError(err)
		}
		ce.marshaljt() // <keyword> , </keyword>

		// Check for varName
		if err := ce.jt.Advance(); err != nil {
			return err
		}
		_, err := ce.jt.Identifier()
		if err != nil {
			return SyntaxError(err)
		}
		ce.marshaljt() // <identifier> varName </identifier>
		if err := ce.jt.Advance(); err != nil {
			return err
		}
	}

	// Should wind up at a ";"
	if err = ce.compileSymbol(";"); err != nil {
		return SyntaxError(err)
	}

	return nil
}

// ('constructor' | 'function' | 'method') ('void' | type) subroutineName '(' parameterList ')' subroutineBody
func (ce *CompilationEngine) compileSubroutine() error {
	ce.openXMLTag("subroutineDec")        // <subroutineDec>
	defer ce.closeXMLTag("subroutineDec") //</subroutineDec>

	// Confirm that current token is "constructor" or "function" or "method"
	if kw, err := ce.jt.KeyWord(); !(kw == "constructor" || kw == "function" || kw == "method") {
		if err != nil {
			return SyntaxError(err)
		}
		return SyntaxError(fmt.Errorf("expected %v \"constructor\" or \"function\" or \"method\"", keyWord))
	}

	// found "constructor" or "function" or "method"
	ce.marshaljt() // <keyword> * </keyword>

	if err := ce.jt.Advance(); err != nil {
		return err
	}

	// Compile a ('void' | type)
	if err := ce.compileVoidOrType(); err != nil {
		return err
	}

	if err := ce.jt.Advance(); err != nil {
		return err
	}

	// check for subroutineName
	_, err := ce.jt.Identifier()
	if err != nil {
		return SyntaxError(err)
	}

	ce.marshaljt() // <identifier> subroutineName </identifier>

	// Eat the subroutineName
	if err := ce.jt.Advance(); err != nil {
		return SyntaxError(err)
	}

	if err := ce.compileParameterList(); err != nil {
		return SyntaxError(err)
	}

	return nil
}

func (ce *CompilationEngine) compileSymbol(sym string) error {
	if s, err := ce.jt.Symbol(); s != sym {
		if err != nil {
			return SyntaxError(err)
		}
		return SyntaxError(fmt.Errorf("expected the %v \"%v\"", symbol, sym))
	}

	ce.marshaljt() // <symbol> sym </symbol>
	return nil
}

// '(' ((type varName)(',' type varName)*)? ')'
func (ce *CompilationEngine) compileParameterList() error {
	if err := ce.compileSymbol("("); err != nil {
		return SyntaxError(err)
	}
	if err := ce.jt.Advance(); err != nil {
		return err
	}

	ce.openXMLTag("parameterList") // <parameterList>

	// While we have yet to hit the closing ")"
	for sym, _ := ce.jt.Symbol(); sym != ")"; sym, _ = ce.jt.Symbol() {
		// First token should be a type
		err := ce.compileType()
		if err != nil {
			return err
		}

		// Eat the type token
		if err := ce.jt.Advance(); err != nil {
			return err
		}

		// Next should be a varName
		if ce.jt.TokenType() != identifier {
			return SyntaxError(fmt.Errorf("Expected an %v for the varName", identifier))
		}
		ce.marshaljt() // <identifier> varName </identifier>

		// Eat the varName token
		if err := ce.jt.Advance(); err != nil {
			return err
		}

		// Now we should be at either a "," or the closing ")"
		sym, err := ce.jt.Symbol()
		if err != nil || !(sym == "," || sym == ")") {
			return SyntaxError(fmt.Errorf("Expected either a \",\" or a \")\""))
		}
		if sym == "," {
			ce.compileSymbol(sym)
			if err := ce.jt.Advance(); err != nil {
				return err
			}
		}
		// Else we were at a ")", let the loop break
	}

	ce.closeXMLTag("parameterList") // </parameterList>
	// Now we should be at the closing ")"
	ce.compileSymbol(")")
	if err := ce.jt.Advance(); err != nil {
		return err
	}

	return nil
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
