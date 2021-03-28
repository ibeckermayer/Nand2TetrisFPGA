package compiler

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
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
	jackFilePath string         // The name of the .jack input file to be compiled.
	jt           *JackTokenizer // A tokenizer set up to tokenize the file we want to compile
	outputFile   *os.File       // The output file
	xmlEnc       *xml.Encoder   // The xml encoder for testing
}

// NewCompilationEngine takes in a path to a jack file and returns an initialized CompilationEngine
// ready to compile it.
func NewCompilationEngine(jackFilePath string) (*CompilationEngine, error) {
	ce := &CompilationEngine{
		jackFilePath: jackFilePath,
	}

	// Initialize the ce's corresponding JackTokenizer
	jt, err := NewJackTokenizer(ce.jackFilePath)
	if err != nil {
		return nil, err
	}
	ce.jt = jt

	// Create the output file
	outputFile, err := os.Create(fmt.Sprintf("%v_out.xml", ce.jackFilePath[0:len(ce.jackFilePath)-len(".jack")]))
	if err != nil {
		return nil, err
	}
	ce.outputFile = outputFile

	// Create xml encoder
	ce.xmlEnc = xml.NewEncoder(outputFile)
	ce.xmlEnc.Indent("", "  ")

	return ce, nil
}

// Run runs the compiler on ce.jackFilePath
func (ce *CompilationEngine) Run() error {
	defer ce.outputFile.Close()
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

// '{' varDec* statements '}'
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

	// varDec*
	for kw, _ := ce.jt.KeyWord(); kw == "var"; kw, _ = ce.jt.KeyWord() {
		if err = ce.compileVarDec(); err != nil {
			return SyntaxError(err)
		}
		// compiled a varDec, advance and try again or move on
		if err = ce.advance(); err != nil {
			return SyntaxError(err)
		}
	}

	if err = ce.compileStatements(); err != nil {
		return SyntaxError(err)
	}

	if err = ce.compileSymbol("}"); err != nil {
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

func (ce *CompilationEngine) compileIdentifier() error {
	// Next should be a varName
	if ce.jt.TokenType() != identifier {
		return SyntaxError(fmt.Errorf("Expected an %v", identifier))
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

	// compile var
	if err := ce.compileKeyword("var"); err != nil {
		return SyntaxError(err)
	}

	// Advance and compile type
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileType(); err != nil {
		return err
	}

	// Advance and compile varName
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileVarName(); err != nil {
		return SyntaxError(err)
	}

	// Now advance again and see if there's a (',' varName)*
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	for sym, _ := ce.jt.Symbol(); sym == ","; sym, _ = ce.jt.Symbol() {
		// (',' varName)*
		if err := ce.compileSymbol(","); err != nil {
			return SyntaxError(err)
		}

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

	return nil
}

func isSatementKeyword(kw string) bool {
	return kw == "let" || kw == "if" || kw == "while" || kw == "do" || kw == "return"
}

// statement*
// statement: letStatement | ifStatement | whileStatement | doStatement | returnStatement
// Whichever function calls this function should have advanced us to the first word of the first statement.
// This function exits after a loop hits a non-statement-starting keyword, so the caller should expect to
// already be at the next token.
func (ce *CompilationEngine) compileStatements() error {
	ce.openXMLTag("statements")
	defer ce.closeXMLTag("statements")

	for kw, err := ce.jt.KeyWord(); isSatementKeyword(kw); kw, err = ce.jt.KeyWord() {
		if err != nil {
			return SyntaxError(err)
		}

		// Each of the following compileX() functions should eat their final character before returning
		// in order to simplify switch case logic. There is no need to call advance() after compileX(),
		// since we can expect every compileX() below will have done so already.
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
	}

	return nil
}

// subroutineName '(' expressionList ')' |
// (className | varName) '.' subroutineName '(' expressionList ')'
func (ce *CompilationEngine) compileSubroutineCall() error {
	if err := ce.compileIdentifier(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	var err error
	var sym string
	sym, err = ce.jt.Symbol() // either a "." or a "("
	if err != nil {
		return SyntaxError(fmt.Errorf("expected a \".\" or \"(\""))
	}
	if sym == "." {
		// ".", so compile and eat it and compile its subsequent subroutineName
		if err := ce.compileSymbol("."); err != nil {
			return SyntaxError(err)
		}
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
		// subroutineName
		if err := ce.compileIdentifier(); err != nil {
			return SyntaxError(err)
		}
		// advance and set external scope's sym variable to the next sym, which ought to be a "("
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
		sym, err = ce.jt.Symbol() // should be a "(" (checked outside this if statement)
		if err != nil {
			return SyntaxError(fmt.Errorf("expected a \".\" or \"(\""))
		}
	}

	if err := ce.compileSymbol("("); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileExpressionList(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileSymbol(")"); err != nil {
		return SyntaxError(err)
	}

	return nil
}

// 'do' subroutineCall ';'
// Eats it's own final character, so caller needn't immediately call advance() upon this function returning
func (ce *CompilationEngine) compileDo() error {
	ce.openXMLTag("doStatement")
	defer ce.closeXMLTag("doStatement")

	if err := ce.compileKeyword("do"); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileSubroutineCall(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileSymbol(";"); err != nil {
		return SyntaxError(err)
	}
	// eat own final character
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	return nil
}

// 'let' varName ('[' expression ']')? '=' expression ';'
// Eats it's own final character, so caller needn't immediately call advance() upon this function returning
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
		if err := ce.compileSymbol("["); err != nil {
			return SyntaxError(err)
		}
		// eat the '[' and compile expression
		if err = ce.advance(); err != nil {
			return SyntaxError(err)
		}
		if err = ce.compileExpression(); err != nil {
			return SyntaxError(err)
		}
		if err := ce.compileSymbol("]"); err != nil {
			return SyntaxError(err)
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

	if err = ce.compileSymbol(";"); err != nil {
		return SyntaxError(err)
	}

	// Eat own final character
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	return nil
}

// 'while '(' expression ')' '{' statements '}'
// eats its own final token before returning, so calling function can expect to already be on the next token
func (ce *CompilationEngine) compileWhile() error {
	ce.openXMLTag("whileStatement")
	defer ce.closeXMLTag("whileStatement")

	if err := ce.compileKeyword("while"); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileSymbol("("); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileExpression(); err != nil {
		return SyntaxError(err)
	}

	if err := ce.compileSymbol(")"); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileSymbol("{"); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileStatements(); err != nil {
		return SyntaxError(err)
	}

	if err := ce.compileSymbol("}"); err != nil {
		return SyntaxError(err)
	}
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}

	return nil
}

// 'return' expression? ';'
// Eats it's own final character, so caller needn't immediately call advance() upon this function returning
func (ce *CompilationEngine) compileReturn() error {
	ce.openXMLTag("returnStatement")
	defer ce.closeXMLTag("returnStatement")

	if err := ce.compileKeyword("return"); err != nil {
		return SyntaxError(err)
	}

	// advance and check for ';'
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}

	// if not ';', should be an expression
	if sym, _ := ce.jt.Symbol(); sym != ";" {
		if err := ce.compileExpression(); err != nil {
			return SyntaxError(err)
		}
	}

	// now compile the ';'
	if err := ce.compileSymbol(";"); err != nil {
		return SyntaxError(err)
	}

	// eat own final character
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}

	return nil
}

// 'if' '(' expression ')' '{' statements '}'
// ('else' '{' statements '}')?
// Eats it's own final character, so caller needn't immediately call advance() upon this function returning
func (ce *CompilationEngine) compileIf() error {
	// this chunk of code is used in regular if and if-else, so abstracted into an internal
	// function to avoid having it written twice
	compileStatementsSubsection := func() error {
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
		if err := ce.compileSymbol("{"); err != nil {
			return SyntaxError(err)
		}

		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
		if err := ce.compileStatements(); err != nil {
			return SyntaxError(err)
		}

		// compileStatements loops us to next token so no need to call advance()
		if err := ce.compileSymbol("}"); err != nil {
			return SyntaxError(err)
		}
		return nil
	}

	ce.openXMLTag("ifStatement")
	defer ce.closeXMLTag("ifStatement")

	if err := ce.compileKeyword("if"); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileSymbol("("); err != nil {
		return SyntaxError(err)
	}
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}

	if err := ce.compileExpression(); err != nil {
		return SyntaxError(err)
	}

	// compileExpression loops us to next token so no need to call advance()
	if err := ce.compileSymbol(")"); err != nil {
		return SyntaxError(err)
	}

	if err := compileStatementsSubsection(); err != nil {
		return SyntaxError(err)
	}

	// advance and check for an else statement
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if ce.jt.TokenType() == keyWord {
		if kw, _ := ce.jt.KeyWord(); kw == "else" {
			if err := ce.compileKeyword("else"); err != nil {
				return SyntaxError(err)
			}

			if err := compileStatementsSubsection(); err != nil {
				return SyntaxError(err)
			}

			// advance after the else is compiled, so that in both the if and if-else cases
			// the function returns with the next token in the chamber (iow this function's api
			// should be the same in both the if and if-else cases)
			if err := ce.advance(); err != nil {
				return SyntaxError(err)
			}
		}
	}

	return nil
}

// isOp returns true if the passed symbol is a valid op, else returns false
func isOp(sym string) bool {
	return (sym == "+" || sym == "-" || sym == "*" || sym == "/" || sym == "&" || sym == "|" || sym == "<" || sym == ">" || sym == "=")
}

// term (op term)*
// Loops until a non (op term) is found, so caller should expect to be at the next token when this function returns.
func (ce *CompilationEngine) compileExpression() error {
	ce.openXMLTag("expression")
	defer ce.closeXMLTag("expression")

	if err := ce.compileTerm(); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}

	for sym, _ := ce.jt.Symbol(); isOp(sym); sym, _ = ce.jt.Symbol() {
		if err := ce.compileSymbol(sym); err != nil {
			return SyntaxError(err)
		}
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
		if err := ce.compileTerm(); err != nil {
			return SyntaxError(err)
		}
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
	}

	return nil
}

// integerConstant | stringConstant | keywordConstant |
// varName | varName '[' expression ']' | subroutineCall |
// '(' expression ')' | unaryOp term
func (ce *CompilationEngine) compileTerm() error {
	ce.openXMLTag("term")
	defer ce.closeXMLTag("term")

	if ce.jt.TokenType() == intConst {
		ce.openXMLTag("integerConstant")
		ic, err := ce.jt.IntVal()
		if err != nil {
			return SyntaxError(err)
		}
		ce.writeXMLData(strconv.Itoa(ic))
		ce.closeXMLTag("integerConstant")
	} else if ce.jt.TokenType() == strConst {
		ce.openXMLTag("stringConstant")
		sc, err := ce.jt.StringVal()
		if err != nil {
			return SyntaxError(err)
		}
		ce.writeXMLData(sc)
		ce.closeXMLTag("stringConstant")
	} else if ce.jt.TokenType() == keyWord {
		kw, err := ce.jt.KeyWord()
		if err != nil {
			return SyntaxError(err)
		}
		if !(kw == "true" || kw == "false" || kw == "null" || kw == "this") {
			return SyntaxError(fmt.Errorf("term keyWord must be one of \"true\", \"false\", \"null\", or \"this\""))
		}
		return ce.marshaljt()
	} else if ce.jt.TokenType() == symbol {
		// '(' expression ')' | unaryOp term
		sym, err := ce.jt.Symbol()
		if err != nil {
			return SyntaxError(err)
		}
		if sym == "(" {
			if err := ce.compileSymbol("("); err != nil {
				return SyntaxError(err)
			}
			if err := ce.advance(); err != nil {
				return SyntaxError(err)
			}
			if err := ce.compileExpression(); err != nil {
				return SyntaxError(err)
			}
			if err = ce.compileSymbol(")"); err != nil {
				return SyntaxError(err)
			}
		} else if sym == "-" {
			if err := ce.compileSymbol("-"); err != nil {
				return SyntaxError(err)
			}
			if err := ce.advance(); err != nil {
				return SyntaxError(err)
			}
			if err := ce.compileTerm(); err != nil {
				return SyntaxError(err)
			}
		} else if sym == "~" {
			if err := ce.compileSymbol("~"); err != nil {
				return SyntaxError(err)
			}
			if err := ce.advance(); err != nil {
				return SyntaxError(err)
			}
			if err := ce.compileTerm(); err != nil {
				return SyntaxError(err)
			}

		} else {
			return SyntaxError(fmt.Errorf("invalid symbol in term, symbol must be one of \"(\" or \"-\" or \"~\""))
		}
	} else if ce.jt.TokenType() == identifier {
		// TODO: need to implement look ahead checks to determine
		// varName | varName '[' expression ']' | subroutineCall
		var err error
		var peeked byte
		peeked, err = ce.jt.Peek()
		if err != nil {
			return SyntaxError(err)
		}
		if peeked == '.' || peeked == '(' {
			return ce.compileSubroutineCall()
		} else if peeked == '[' {
			if err := ce.compileVarName(); err != nil {
				return SyntaxError(err)
			}
			if err := ce.advance(); err != nil {
				return SyntaxError(err)
			}
			if err := ce.compileSymbol("["); err != nil {
				return SyntaxError(err)
			}
			if err := ce.advance(); err != nil {
				return SyntaxError(err)
			}
			if err := ce.compileExpression(); err != nil {
				return SyntaxError(err)
			}
			if err := ce.compileSymbol("]"); err != nil {
				return SyntaxError(err)
			}
		} else {
			return ce.compileVarName()
		}
	} else {
		return SyntaxError(fmt.Errorf("unknown error"))
	}

	return nil
}

// (expression (',' expression)* )?
// caller should expect to be at the next token when this fucntion returns
func (ce *CompilationEngine) compileExpressionList() error {
	ce.openXMLTag("expressionList")
	defer ce.closeXMLTag("expressionList")

	// Advance and check if we are at a closing parenthesis
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	sym, _ := ce.jt.Symbol()
	if ce.jt.TokenType() == symbol && sym == ")" {
		// The next token was a closing parenthesis; simply return, and the caller is expected
		// to compile the closing paren
		return nil
	}

	// compile expression
	if err := ce.compileExpression(); err != nil {
		return SyntaxError(err)
	}
	for sym, _ := ce.jt.Symbol(); sym == ","; sym, _ = ce.jt.Symbol() {
		ce.compileSymbol(sym) // compile ","
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
		if err := ce.compileExpression(); err != nil {
			return SyntaxError(err)
		}
	}

	return nil
}
