package compiler

import (
	"errors"
	"fmt"
	"log"
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
	jt *JackTokenizer // A tokenizer set up to tokenize the file we want to compile
	st *SymbolTable   // The symbol table
	cw *CodeWriter    // The code writer
}

// NewCompilationEngine takes in a path to a jack file and returns an initialized CompilationEngine
// ready to compile it. Since we assume one jack class per file, we will assume one compilation engine
// per class.
func NewCompilationEngine(jackFilePath string) (*CompilationEngine, error) {
	ce := &CompilationEngine{}

	// Initialize the ce's corresponding JackTokenizer
	jt, err := NewJackTokenizer(jackFilePath)
	if err != nil {
		return nil, err
	}
	ce.jt = jt

	// Create the code writer
	cw, err := NewCodeWriter(jackFilePath)
	if err != nil {
		return nil, err
	}
	ce.cw = cw

	// Create the symbol table
	st := NewSymbolTable()
	ce.st = st

	return ce, nil
}

// Run runs the compiler on ce.jackFilePath
func (ce *CompilationEngine) Run() error {
	// Close the output file after run
	defer ce.cw.Close()

	// Advance to eat the first token and call compileClass, which will recursively compile the entire file
	if err := ce.advance(); err != nil {
		return err
	}

	return ce.compileClass()
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
	// Check that first token is "class"
	if err := ce.checkForKeyword("class"); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return err
	}

	// Check that next token is an identifier
	className, err := ce.jt.Identifier()
	if err != nil {
		return err
	}

	// set className for function naming
	ce.cw.className = className

	if err := ce.advance(); err != nil {
		return err
	}

	// Check that next token is "{"
	if err := ce.checkForSymbol("{"); err != nil {
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
	if err := ce.checkForSymbol("}"); err != nil {
		return SyntaxError(err)
	}

	return nil
}

// 'void | 'int' | 'char' | 'boolean' | className
func (ce *CompilationEngine) getVoidOrType() (string, error) {
	if kw, err := ce.jt.KeyWord(); kw == "void" {
		return kw, err
	}
	return ce.getType()
}

// 'int' | 'char' | 'boolean' | className
// check for a type and return the string
func (ce *CompilationEngine) getType() (string, error) {
	var errString string = "expected a type: %v className or %v \"int\" or \"char\" or \"boolean\""

	switch ce.jt.TokenType() {
	case keyWord:
		kw, _ := ce.jt.KeyWord()
		if !(kw == "int" || kw == "char" || kw == "boolean") {
			return "", SyntaxError(fmt.Errorf(errString, identifier, keyWord))
		}
		return kw, nil
	case identifier:
		id, _ := ce.jt.Identifier()
		return id, nil
	default:
		return "", SyntaxError(fmt.Errorf(errString, identifier, keyWord))
	}
}

// ('static' | 'field') type varName (',' varName)* ';'
func (ce *CompilationEngine) compileClassVarDec() error {
	// extract "static" or "field"
	kind, _ := ce.jt.KeyWord()

	// get type
	if err := ce.advance(); err != nil {
		return err
	}
	type_, err := ce.getType()
	if err != nil {
		return err
	}

	// Check for varName
	if err := ce.advance(); err != nil {
		return err
	}
	name, err := ce.getVarName()
	if err != nil {
		return SyntaxError(err)
	}

	// Define a new entry in the symbol table
	ce.st.Define(name, type_, Kind(kind))

	// Check for a comma separated list of more varNames
	if err := ce.advance(); err != nil {
		return err
	}
	for sym, err := ce.jt.Symbol(); sym == ","; sym, err = ce.jt.Symbol() {
		if err != nil {
			return SyntaxError(err)
		}

		// Get name
		if err := ce.advance(); err != nil {
			return err
		}
		name, err := ce.getVarName()
		if err != nil {
			return SyntaxError(err)
		}

		ce.st.Define(name, type_, Kind(kind))

		if err := ce.advance(); err != nil {
			return err
		}
	}

	// Should wind up at a ";"
	if err := ce.checkForSymbol(";"); err != nil {
		return SyntaxError(err)
	}

	return nil
}

// ('constructor' | 'function' | 'method') ('void' | type) subroutineName '(' parameterList ')' subroutineBody
func (ce *CompilationEngine) compileSubroutine() error {
	// Clear the subroutine table
	ce.st.StartSubroutine()

	// Calling function checked that current token is "constructor" or "function" or "method"
	if kw, err := ce.jt.KeyWord(); !(kw == "constructor" || kw == "function" || kw == "method") {
		if err != nil {
			return SyntaxError(err)
		}
		return SyntaxError(fmt.Errorf("expected %v \"constructor\" or \"function\" or \"method\"", keyWord))
	}

	if err := ce.advance(); err != nil {
		return err
	}

	// Get a ('void' | type)
	_, err := ce.getVoidOrType()
	if err != nil {
		return err
	}

	// Get the subroutineName
	if err := ce.advance(); err != nil {
		return err
	}
	subroutineName, err := ce.jt.Identifier()
	if err != nil {
		return SyntaxError(err)
	}

	// Eat the subroutineName
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}

	if err := ce.compileParameterList(); err != nil {
		return SyntaxError(err)
	}

	// subroutineBody: '{' varDec* statements '}'

	// Eat what should be a "{"
	if err = ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err = ce.checkForSymbol("{"); err != nil {
		return SyntaxError(err)
	}
	// Eat the "{"
	if err = ce.advance(); err != nil {
		return SyntaxError(err)
	}

	// varDec*
	if err := ce.compileVarDecs(); err != nil {
		return SyntaxError(err)
	}

	// Now that all the vars have been declared, we know how to declare the vm code function
	ce.cw.WriteFunction(subroutineName, ce.st.VarCount(KIND_VAR))

	// Now that function has been declared, write its body
	if err = ce.compileStatements(); err != nil {
		return SyntaxError(err)
	}

	if err = ce.checkForSymbol("}"); err != nil {
		return SyntaxError(err)
	}

	return nil
}

// varDec*
// compileVarDecs compiles all the variable declarations, so that when it returns a call to
// ce.st.VarCount(KIND_VAR) will give you the number of local variables for the current subroutine
func (ce *CompilationEngine) compileVarDecs() error {
	var err error
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
	return nil
}

// checks that token is identifier and returns it, or a error
func (ce *CompilationEngine) getVarName() (string, error) {
	id, err := ce.jt.Identifier()
	if err != nil {
		return "", SyntaxError(fmt.Errorf("Expected an %v for the varName", identifier))
	}
	return id, nil

}

// checkForIdentifier checks that the current token is an identifier
func (ce *CompilationEngine) checkForIdentifier() error {
	// Next should be a varName
	if ce.jt.TokenType() != identifier {
		return SyntaxError(fmt.Errorf("Expected an %v", identifier))
	}
	return nil
}

// checkForSymbol checks that the passed symbol is currently being parsed
func (ce *CompilationEngine) checkForSymbol(sym string) error {
	if s, err := ce.jt.Symbol(); s != sym {
		if err != nil {
			return SyntaxError(err)
		}
		return SyntaxError(fmt.Errorf("expected the %v \"%v\"", symbol, sym))
	}

	return nil
}

func (ce *CompilationEngine) checkForKeyword(kw string) error {
	if k, err := ce.jt.KeyWord(); k != kw {
		if err != nil {
			return SyntaxError(err)
		}
		return SyntaxError(fmt.Errorf("expected the %v \"%v\"", keyWord, kw))
	}

	return nil
}

// '(' ((type varName)(',' type varName)*)? ')'
func (ce *CompilationEngine) compileParameterList() error {
	if err := ce.checkForSymbol("("); err != nil {
		return SyntaxError(err)
	}
	if err := ce.advance(); err != nil {
		return err
	}

	// While we have yet to hit the closing ")"
	for sym, _ := ce.jt.Symbol(); sym != ")"; sym, _ = ce.jt.Symbol() {
		// First token should be a type
		type_, err := ce.getType()
		if err != nil {
			return err
		}

		// Next should be a varName
		if err := ce.advance(); err != nil {
			return err
		}
		name, err := ce.getVarName()
		if err != nil {
			return SyntaxError(err)
		}

		// Eat the varName token
		if err := ce.advance(); err != nil {
			return err
		}

		// Define the parameter in the symbol table
		ce.st.Define(name, type_, KIND_ARG)

		// Now we should be at either a "," or the closing ")"
		if err := ce.checkForSymbol(","); err == nil {
			// if we're at the ",", advance and run through another round of the loop
			if err := ce.advance(); err != nil {
				return SyntaxError(err)
			}
			// but first check that we didn't bump immediately into a ")", which would be invalid
			if err := ce.checkForSymbol(")"); err == nil {
				return SyntaxError(fmt.Errorf("invalid syntax \",)\""))
			}
		}
		// Else we were at a ")", let the loop break and check for that
	}

	// Now we should be at the closing ")"
	if err := ce.checkForSymbol(")"); err != nil {
		return SyntaxError(err)
	}

	return nil
}

// 'var' type varName (',' varName)* ';'
func (ce *CompilationEngine) compileVarDec() error {
	// check for var
	if err := ce.checkForKeyword("var"); err != nil {
		return SyntaxError(err)
	}

	// Advance and get type
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	type_, err := ce.getType()
	if err != nil {
		return err
	}

	// Advance and get varName
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	name, err := ce.getVarName()
	if err != nil {
		return SyntaxError(err)
	}

	// Add this var to symbol table
	ce.st.Define(name, type_, KIND_VAR)

	// Now advance again and see if there's a (',' varName)*
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	for sym, _ := ce.jt.Symbol(); sym == ","; sym, _ = ce.jt.Symbol() {
		// (',' varName)*
		// Advance and check for "varName"
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
		name, err := ce.getVarName()
		if err != nil {
			return SyntaxError(err)
		}

		// Add this var to symbol table as well
		ce.st.Define(name, type_, KIND_VAR)

		// Advance to either the next "," and repeat the loop, or break and check for ";"
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
	}

	// Check for ";"
	if err := ce.checkForSymbol(";"); err != nil {
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
	if err := ce.checkForIdentifier(); err != nil {
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
		if err := ce.checkForSymbol("."); err != nil {
			return SyntaxError(err)
		}
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
		// subroutineName
		if err := ce.checkForIdentifier(); err != nil {
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

	if err := ce.checkForSymbol("("); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileExpressionList(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.checkForSymbol(")"); err != nil {
		return SyntaxError(err)
	}

	return nil
}

// 'do' subroutineCall ';'
// Eats it's own final character, so caller needn't immediately call advance() upon this function returning
func (ce *CompilationEngine) compileDo() error {
	if err := ce.checkForKeyword("do"); err != nil {
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
	if err := ce.checkForSymbol(";"); err != nil {
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
	var err error

	if err = ce.checkForKeyword("let"); err != nil {
		return SyntaxError(err)
	}

	// advance and get varName
	if err = ce.advance(); err != nil {
		return SyntaxError(err)
	}
	_, err = ce.getVarName()
	if err != nil {
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
		if err := ce.checkForSymbol("["); err != nil {
			return SyntaxError(err)
		}
		// eat the '[' and compile expression
		if err = ce.advance(); err != nil {
			return SyntaxError(err)
		}
		if err = ce.compileExpression(); err != nil {
			return SyntaxError(err)
		}
		if err := ce.checkForSymbol("]"); err != nil {
			return SyntaxError(err)
		}
		if err = ce.advance(); err != nil {
			return SyntaxError(err)
		}
	}

	// demand, compile, and eat '='
	if err = ce.checkForSymbol("="); err != nil {
		return SyntaxError(err)
	}
	if err = ce.advance(); err != nil {
		return SyntaxError(err)
	}

	// compile expression
	if err = ce.compileExpression(); err != nil {
		return SyntaxError(err)
	}

	if err = ce.checkForSymbol(";"); err != nil {
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
	if err := ce.checkForKeyword("while"); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.checkForSymbol("("); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileExpression(); err != nil {
		return SyntaxError(err)
	}

	if err := ce.checkForSymbol(")"); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.checkForSymbol("{"); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.compileStatements(); err != nil {
		return SyntaxError(err)
	}

	if err := ce.checkForSymbol("}"); err != nil {
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
	if err := ce.checkForKeyword("return"); err != nil {
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
	if err := ce.checkForSymbol(";"); err != nil {
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
		if err := ce.checkForSymbol("{"); err != nil {
			return SyntaxError(err)
		}

		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
		if err := ce.compileStatements(); err != nil {
			return SyntaxError(err)
		}

		// compileStatements loops us to next token so no need to call advance()
		if err := ce.checkForSymbol("}"); err != nil {
			return SyntaxError(err)
		}
		return nil
	}

	if err := ce.checkForKeyword("if"); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}
	if err := ce.checkForSymbol("("); err != nil {
		return SyntaxError(err)
	}
	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}

	if err := ce.compileExpression(); err != nil {
		return SyntaxError(err)
	}

	// compileExpression loops us to next token so no need to call advance()
	if err := ce.checkForSymbol(")"); err != nil {
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
			if err := ce.checkForKeyword("else"); err != nil {
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
	if err := ce.compileTerm(); err != nil {
		return SyntaxError(err)
	}

	if err := ce.advance(); err != nil {
		return SyntaxError(err)
	}

	for sym, _ := ce.jt.Symbol(); isOp(sym); sym, _ = ce.jt.Symbol() {
		if err := ce.checkForSymbol(sym); err != nil {
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
	if ce.jt.TokenType() == intConst {
		_, err := ce.jt.IntVal()
		if err != nil {
			return SyntaxError(err)
		}
	} else if ce.jt.TokenType() == strConst {
		_, err := ce.jt.StringVal()
		if err != nil {
			return SyntaxError(err)
		}
	} else if ce.jt.TokenType() == keyWord {
		kw, err := ce.jt.KeyWord()
		if err != nil {
			return SyntaxError(err)
		}
		if !(kw == "true" || kw == "false" || kw == "null" || kw == "this") {
			return SyntaxError(fmt.Errorf("term keyWord must be one of \"true\", \"false\", \"null\", or \"this\""))
		}
	} else if ce.jt.TokenType() == symbol {
		// '(' expression ')' | unaryOp term
		sym, err := ce.jt.Symbol()
		if err != nil {
			return SyntaxError(err)
		}
		if sym == "(" {
			if err := ce.advance(); err != nil {
				return SyntaxError(err)
			}
			if err := ce.compileExpression(); err != nil {
				return SyntaxError(err)
			}
			if err = ce.checkForSymbol(")"); err != nil {
				return SyntaxError(err)
			}
		} else if sym == "-" {
			if err := ce.advance(); err != nil {
				return SyntaxError(err)
			}
			if err := ce.compileTerm(); err != nil {
				return SyntaxError(err)
			}
		} else if sym == "~" {
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
			_, err = ce.getVarName()
			if err != nil {
				return SyntaxError(err)
			}
			if err := ce.advance(); err != nil {
				return SyntaxError(err)
			}
			if err := ce.checkForSymbol("["); err != nil {
				return SyntaxError(err)
			}
			if err := ce.advance(); err != nil {
				return SyntaxError(err)
			}
			if err := ce.compileExpression(); err != nil {
				return SyntaxError(err)
			}
			if err := ce.checkForSymbol("]"); err != nil {
				return SyntaxError(err)
			}
		} else {
			_, err := ce.getVarName()
			return err
		}
	} else {
		return SyntaxError(fmt.Errorf("unknown error"))
	}

	return nil
}

// (expression (',' expression)* )?
// caller should expect to be at the next token when this function returns
func (ce *CompilationEngine) compileExpressionList() error {
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
		ce.checkForSymbol(sym)
		if err := ce.advance(); err != nil {
			return SyntaxError(err)
		}
		if err := ce.compileExpression(); err != nil {
			return SyntaxError(err)
		}
	}

	return nil
}
