package compiler

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
)

//SyntaxError logs the function name, file, and line number
func SyntaxError(err error) error {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, fn, line, _ := runtime.Caller(1)

		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
		return err
	}
	return nil
}

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
	if err := ce.advance(); err != nil {
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

// advance throws an error if there are no more tokens, else it returns ce.jt.Advance()
func (ce *CompilationEngine) advance() error {
	if !ce.jt.HasMoreTokens() {
		return errors.New("ran out of tokens")
	}
	return ce.jt.Advance()
}

// Assumes the first token has already been consumed (should be 'class').
// 'class' className '{' classVarDec* subroutineDec* '}'
func (ce *CompilationEngine) compileClass() error {
	ce.openXMLTag("class")        // <class>
	defer ce.closeXMLTag("class") // defer closing </class>

	// Check that first token is "class"
	if err := ce.compileKeyword("class"); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return err
	}

	// Check that next token is an identifier
	if _, err := ce.jt.Identifier(); err != nil {
		return err
	}

	ce.marshaljt() // <identifier> ClassName </identifier>

	if err := ce.advance(); err != nil {
		return err
	}

	// Check that next token is "{"
	if err := ce.compileSymbol("{"); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
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

		if err := ce.advance(); err != nil {
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

		if err := ce.advance(); err != nil {
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

	if err := ce.advance(); err != nil {
		return err
	}

	// try to compile a type
	if err := ce.compileType(); err != nil {
		return err
	}

	// Check for varName
	if err := ce.advance(); err != nil {
		return err
	}
	if err := ce.compileVarName(); err != nil {
		return SyntaxError(err)
	}

	// Check for a comma separated list of more varNames
	if err := ce.advance(); err != nil {
		return err
	}
	for sym, err := ce.jt.Symbol(); sym == ","; sym, err = ce.jt.Symbol() {
		if err != nil {
			return SyntaxError(err)
		}
		ce.marshaljt() // <symbol> , </symbol>

		// Check for varName
		if err := ce.advance(); err != nil {
			return err
		}
		if err := ce.compileVarName(); err != nil {
			return SyntaxError(err)
		}

		if err := ce.advance(); err != nil {
			return err
		}
	}

	// Should wind up at a ";"
	if err := ce.compileSymbol(";"); err != nil {
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

	if err := ce.advance(); err != nil {
		return err
	}

	// Compile a ('void' | type)
	if err := ce.compileVoidOrType(); err != nil {
		return err
	}

	if err := ce.advance(); err != nil {
		return err
	}

	// check for subroutineName
	_, err := ce.jt.Identifier()
	if err != nil {
		return SyntaxError(err)
	}

	ce.marshaljt() // <identifier> subroutineName </identifier>

	// Eat the subroutineName
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}

	if err := ce.compileParameterList(); err != nil {
		return SyntaxError(err)
	}

	if err := ce.compileSubroutineBody(); err != nil {
		return SyntaxError(err)
	}

	return nil
}

// '{' varDec* statements'}'
func (ce *CompilationEngine) compileSubroutineBody() error {
	ce.openXMLTag("subroutineBody")        // <subroutineBody>
	defer ce.closeXMLTag("subroutineBody") //</subroutineBody>

	var err error

	// Eat what should be a "{"
	if err = ce.advance(); err != nil {
		return SyntaxError(err)
	}

	if err = ce.compileSymbol("{"); err != nil {
		return SyntaxError(err)
	}

	// Eat the "{"
	if err = ce.advance(); err != nil {
		return SyntaxError(err)
	}

	if err = ce.compileVarDec(); err != nil {
		return SyntaxError(err)
	}

	if err = ce.compileStatements(); err != nil {
		return SyntaxError(err)
	}

	return nil
}

func (ce *CompilationEngine) compileVarName() error {
	// Next should be a varName
	if ce.jt.TokenType() != identifier {
		return SyntaxError(fmt.Errorf("Expected an %v for the varName", identifier))
	}
	return ce.marshaljt() // <identifier> varName </identifier>
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

func (ce *CompilationEngine) compileKeyword(kw string) error {
	if k, err := ce.jt.KeyWord(); k != kw {
		if err != nil {
			return SyntaxError(err)
		}
		return SyntaxError(fmt.Errorf("expected the %v \"%v\"", keyWord, kw))
	}

	ce.marshaljt() // <symbol> kw </symbol>
	return nil
}

// '(' ((type varName)(',' type varName)*)? ')'
func (ce *CompilationEngine) compileParameterList() error {
	if err := ce.compileSymbol("("); err != nil {
		return SyntaxError(err)
	}
	if err := ce.advance(); err != nil {
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
		if err := ce.advance(); err != nil {
			return err
		}

		// Next should be a varName
		if err := ce.compileVarName(); err != nil {
			return SyntaxError(err)
		}

		// Eat the varName token
		if err := ce.advance(); err != nil {
			return err
		}

		// Now we should be at either a "," or the closing ")"
		sym, err := ce.jt.Symbol()
		if err != nil || !(sym == "," || sym == ")") {
			return SyntaxError(fmt.Errorf("Expected either a \",\" or a \")\""))
		}
		if sym == "," {
			ce.compileSymbol(sym)
			if err := ce.advance(); err != nil {
				return err
			}
		}
		// Else we were at a ")", let the loop break
	}

	ce.closeXMLTag("parameterList") // </parameterList>
	// Now we should be at the closing ")"
	ce.compileSymbol(")")
	return nil
}

// 'var' type varName (',' varName)* ';'
func (ce *CompilationEngine) compileVarDec() error {
	ce.openXMLTag("varDec")
	defer ce.closeXMLTag("varDec")

	// Check for 'var'
	if kw, err := ce.jt.KeyWord(); kw != "var" || err != nil {
		return SyntaxError(err)
	}

	// 'var' type varName (',' varName)* ';'
	for kw, _ := ce.jt.KeyWord(); kw == "var"; kw, _ = ce.jt.KeyWord() {
		ce.marshaljt() // <keyword> var </keyword>
		// Advance and compile type
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
		if err := ce.compileType(); err != nil {
			return err
		}
		// Advance and compile varName (',' varName)*
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
		// Next should be a varName
		if err := ce.compileVarName(); err != nil {
			return SyntaxError(err)
		}

		// eat the varName
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}

		for sym, _ := ce.jt.Symbol(); sym == ","; sym, _ = ce.jt.Symbol() {
			ce.marshaljt() // <keyword> , </keyword>
			// Advance and check for "varName"
			if err := ce.advance(); err != nil {
				return SyntaxError(err)
			}
			if err := ce.compileVarName(); err != nil {
				return SyntaxError(err)
			}

			// Advance to either the next "," and repeat the loop, or break and check for ";"
			if err := ce.advance(); err != nil {
				return SyntaxError(err)
			}
		}

		// Check for ";"
		if err := ce.compileSymbol(";"); err != nil {
			return SyntaxError(err)
		}
		// Advance, either to another "var" in which case the loop will repeat, or move on to the next type of element
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
	}

	return nil
}

func isSatementKeyword(kw string) bool {
	return kw == "let" || kw == "if" || kw == "while" || kw == "do" || kw == "return"
}

// statement*
// statement: letStatement | ifStatement | whileStatement | doStatement | returnStatement
func (ce *CompilationEngine) compileStatements() error {
	ce.openXMLTag("statements")
	defer ce.closeXMLTag("statements")

	// compileVarDec is called before this function, and will have advanced us to the first word of the first statement

	for kw, err := ce.jt.KeyWord(); isSatementKeyword(kw); kw, err = ce.jt.KeyWord() {
		if err != nil {
			return SyntaxError(err)
		}

		switch kw {
		case "let":
			if err := ce.compileLet(); err != nil {
				return SyntaxError(err)
			}
		case "if":
			if err := ce.compileIf(); err != nil {
				return SyntaxError(err)
			}
		case "while":
			if err := ce.compileWhile(); err != nil {
				return SyntaxError(err)
			}
		case "do":
			if err := ce.compileDo(); err != nil {
				return SyntaxError(err)
			}
		case "return":
			if err := ce.compileReturn(); err != nil {
				return SyntaxError(err)
			}
		default:
			return SyntaxError(fmt.Errorf("unexpected error in compileStatements, should be impossible"))
		}

		// Advance in order to check for another statement keyword on the next round of the loop
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
	}

	return nil
}

func (ce *CompilationEngine) compileDo() error {
	panic("not implemented") // TODO: Implement
}

// 'let' varName ('[' expression ']')? '=' expression ';'
func (ce *CompilationEngine) compileLet() error {
	ce.openXMLTag("letStatement")
	defer ce.closeXMLTag("letStatement")
	var err error

	if err = ce.compileKeyword("let"); err != nil {
		return SyntaxError(err)
	}

	// advance and compile varName
	if err = ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err = ce.compileVarName(); err != nil {
		return SyntaxError(err)
	}

	// advance and check for '['
	if err = ce.advance(); err != nil {
		return SyntaxError(err)
	}
	var sym string
	if sym, err = ce.jt.Symbol(); err != nil {
		return SyntaxError(err)
	}
	if sym == "[" {
		// eat the '[' and compile expression
		if err = ce.advance(); err != nil {
			return SyntaxError(err)
		}
		if err = ce.compileExpression(); err != nil {
			return SyntaxError(err)
		}
		// TODO: may need to call ce.advance() here depending on nature of compileExpression()
		// compile closing ']'
		if sym, err = ce.jt.Symbol(); err != nil {
			return SyntaxError(err)
		}
		if sym != "]" {
			return SyntaxError(fmt.Errorf("expected closing \"]\""))
		}
		if err = ce.advance(); err != nil {
			return SyntaxError(err)
		}
	}

	// demand, compile, and eat '='
	if err = ce.compileSymbol("="); err != nil {
		return SyntaxError(err)
	}
	if err = ce.advance(); err != nil {
		return SyntaxError(err)
	}

	// compile expression
	if err = ce.compileExpression(); err != nil {
		return SyntaxError(err)
	}

	// TODO: may need to call ce.advance() here depending on nature of compileExpression()
	if err = ce.compileSymbol(";"); err != nil {
		return SyntaxError(err)
	}
	return nil
}

// 'while '(' expression ')' '{' statements '}'
func (ce *CompilationEngine) compileWhile() error {
	panic("not implemented") // TODO: Implement
}

// 'return' expression? ';'
func (ce *CompilationEngine) compileReturn() error {
	panic("not implemented") // TODO: Implement
}

// 'if' '{' expression '}' '{' statements '}'
// ('else' '{' statements '}')?
func (ce *CompilationEngine) compileIf() error {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) compileExpression() error {
	ce.openXMLTag("expression")
	defer ce.closeXMLTag("expression")

	return nil
}

func (ce *CompilationEngine) compileTerm() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) compileExpressionList() {
	panic("not implemented") // TODO: Implement
}
