package compiler

import (
	"errors"
	"fmt"
	"log"
)

var DEBUG = false

// CompilationEngine effects the actual compilation output.
// Gets its input from a JackTokenizer and emits its parsed structure into an output file/stream.
type CompilationEngine struct {
	jt *JackTokenizer // A tokenizer set up to tokenize the file we want to compile
	st *SymbolTable   // The symbol table
	cw *CodeWriter    // The code writer
	// className is the class name being compiled,
	// set at compileClass by compilation engine
	className string // The name of the class being compiled
	subKind   string // The kind of subroutine being compiled
	subName   string // The name of the subroutine being compiled
	// whileId is used to ensure unique labels are used in the vm code translation
	// for each while loop encountered in the class
	whileId uint
	// ifId is similar to whileId, but used for if-else statements
	ifId uint
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

	ce.whileId = 0
	ce.ifId = 0

	return ce, nil
}

// Run runs the compiler on ce.jackFilePath
func (ce *CompilationEngine) Run() error {
	// Close the output file after run
	defer ce.cw.Close()

	// Advance to eat the first token and call compileClass, which will recursively compile the entire file
	if err := ce.advance(); err != nil {
		return TraceError(err)
	}

	if err := ce.compileClass(); err != nil {
		log.Printf("Compilation error in %v %v %v", ce.className, ce.subKind, ce.subName)
		return TraceError(err)
	}

	return nil
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
		return TraceError(err)
	}

	if err := ce.advance(); err != nil {
		return TraceError(err)
	}

	// Check that next token is an identifier
	className, err := ce.jt.Identifier()
	if err != nil {
		return TraceError(err)
	}

	// set className for function naming
	ce.className = className

	if err := ce.advance(); err != nil {
		return TraceError(err)
	}

	// Check that next token is "{"
	if err := ce.eatSymbol("{"); err != nil {
		return TraceError(err)
	}

	// Loop through and compile all of the classVarDecs
	for kw, err := ce.jt.KeyWord(); kw == "static" || kw == "field"; kw, err = ce.jt.KeyWord() {
		if err != nil {
			return TraceError(err)
		}

		if err := ce.compileClassVarDec(); err != nil {
			return TraceError(err)
		}
	}

	// Loop through and compile all of the subroutines.
	// The previous loop will have called Advance() and then hit a non static/field
	for kw, err := ce.jt.KeyWord(); kw == "constructor" || kw == "function" || kw == "method"; kw, err = ce.jt.KeyWord() {
		if err != nil {
			return TraceError(err)
		}

		if err := ce.compileSubroutine(); err != nil {
			return TraceError(err)
		}

	}

	// Check that next token is "}"
	// The previous loop should have called Advance() for this symbol
	if err := ce.checkForSymbol("}"); err != nil {
		return TraceError(err)
	}

	return nil
}

// getVoidOrType gets the current token as a void or type ('void | 'int' | 'char' | 'boolean' | className)
// and advances the tokenizer. Returns an error if the current token is not a void or type.
func (ce *CompilationEngine) getVoidOrType() (string, error) {
	if kw, err := ce.jt.KeyWord(); kw == "void" {
		if err := ce.advance(); err != nil {
			return "", err
		}
		return kw, err
	}

	return ce.getType()

}

// getType checks for a type ('int' | 'char' | 'boolean' | className) and returns the string,
// advancing the tokenizer. Returns an error if the current token is not a type.
func (ce *CompilationEngine) getType() (string, error) {
	var errString string = "expected a type: %v className or %v \"int\" or \"char\" or \"boolean\""
	var err error
	var type_ string

	switch ce.jt.TokenType() {
	case keyWord:
		type_, err = ce.jt.KeyWord()
		if err != nil {
			return "", TraceError(err)
		}
		if !(type_ == "int" || type_ == "char" || type_ == "boolean") {
			return "", TraceError(fmt.Errorf(errString, identifier, keyWord))
		}
	case identifier:
		type_, err = ce.jt.Identifier()
		if err != nil {
			return "", TraceError(err)
		}
	default:
		return "", TraceError(fmt.Errorf(errString, identifier, keyWord))
	}

	err = ce.advance()
	if err != nil {
		return "", err
	}

	return type_, err
}

// ('static' | 'field') type varName (',' varName)* ';'
func (ce *CompilationEngine) compileClassVarDec() error {
	// extract "static" or "field"
	kind, _ := ce.jt.KeyWord()

	// get type
	if err := ce.advance(); err != nil {
		return TraceError(err)
	}
	type_, err := ce.getType()
	if err != nil {
		return TraceError(err)
	}

	name, err := ce.getVarName()
	if err != nil {
		return TraceError(err)
	}

	// Define a new entry in the symbol table
	ce.st.Define(name, type_, Kind(kind))

	// Check for a comma separated list of more varNames
	for sym, err := ce.jt.Symbol(); sym == ","; sym, err = ce.jt.Symbol() {
		if err != nil {
			return TraceError(err)
		}

		// Get name
		if err := ce.advance(); err != nil {
			return TraceError(err)
		}
		name, err := ce.getVarName()
		if err != nil {
			return TraceError(err)
		}

		ce.st.Define(name, type_, Kind(kind))
	}

	// Should wind up at a ";"
	if err := ce.eatSymbol(";"); err != nil {
		return TraceError(err)
	}

	return nil
}

func (ce *CompilationEngine) getIntVal() (uint, error) {
	i, err := ce.jt.IntVal()
	if err != nil {
		return 0, TraceError(err)
	}
	if err := ce.advance(); err != nil {
		return 0, err
	}
	return i, nil
}

// getKeyWord gets the current token as a keyword and advances the tokenizer.
// Returns an error if the current token is not a keyword.
func (ce *CompilationEngine) getKeyWord() (string, error) {
	kw, err := ce.jt.KeyWord()
	if err != nil {
		return "", TraceError(err)
	}
	if err := ce.advance(); err != nil {
		return "", err
	}
	return kw, nil
}

// getIdentifier gets the current token as an identifier and advances the tokenizer.
// Returns an error if the current token is not an identifier.
func (ce *CompilationEngine) getIdentifier() (string, error) {
	id, err := ce.jt.Identifier()
	if err != nil {
		return "", TraceError(err)
	}
	if err := ce.advance(); err != nil {
		return "", err
	}
	return id, nil
}

// '{' varDec* statements '}'
//
// subKind must be one of "constructor" or "function" or "method", and subName must be the name of the subroutine
func (ce *CompilationEngine) compileSubroutineBody(subKind string, subName string) error {
	var err error

	ce.subKind = subKind
	ce.subName = subName

	if subKind != "constructor" && subKind != "function" && subKind != "method" {
		return TraceError(fmt.Errorf("expected a subroutine kind: \"constructor\" or \"function\" or \"method\""))
	}

	if err = ce.eatSymbol("{"); err != nil {
		return TraceError(err)
	}

	// varDec*
	if err := ce.compileVarDecs(); err != nil {
		return TraceError(err)
	}

	// Now that all the vars have been declared, we know how to declare the vm code function
	ce.cw.WriteFunction(ce.className+"."+subName, ce.st.VarCount(KIND_VAR))

	// If this is a constructor, it is implied that the first part of its body
	// is code to allocate memory for the object and set `this` to the address of that object
	if subKind == "constructor" {
		// All of Jack's data types are 16-bits (1 word) long, so the size of an object is simply the number of fields it has.
		// (If a field is another object, then it's simply a pointer to that object which is 16-bits long.)
		objectSize := ce.st.VarCount(KIND_FIELD)
		// Allocate memory for the object:
		// push objectSize onto the stack
		ce.cw.WritePush(SEG_CONST, objectSize)
		// call Memory.alloc with that objectSize
		ce.cw.WriteCall("Memory.alloc", 1)
		// Memory.alloc returns the address of the allocated object, which is now on top of the stack.
		// Set `this` to the address of that allocated object.
		ce.cw.WritePop(SEG_POINTER, 0)
	} else if subKind == "method" {
		// If this is a method, it is implied that the first part of its body
		// is code to set `this` to the address of the object that the method is being called on.
		// Recall that the first argument of a method is always the object that the method is being called on.
		// Therefore, we need to set `this` to the address of that object.
		ce.cw.WritePush(SEG_ARG, 0)
		ce.cw.WritePop(SEG_POINTER, 0)
	}

	// Now that subroutine has been declared, and any implicit constructor code written, write its body.
	if err = ce.compileStatements(); err != nil {
		return TraceError(err)
	}

	if err = ce.eatSymbol("}"); err != nil {
		return TraceError(err)
	}

	return nil
}

// ('constructor' | 'function' | 'method') ('void' | type) subroutineName '(' parameterList ')' subroutineBody
func (ce *CompilationEngine) compileSubroutine() error {
	// Clear the subroutine table
	ce.st.StartSubroutine()

	// subKind is one of "constructor" or "function" or "method"
	subKind, err := ce.getKeyWord()
	if err != nil {
		return TraceError(err)
	}

	// Get a ('void' | type)
	_, err = ce.getVoidOrType()
	if err != nil {
		return TraceError(err)
	}

	// Get the subroutineName
	subroutineName, err := ce.getIdentifier()
	if err != nil {
		return TraceError(err)
	}

	if err := ce.compileParameterList(subKind); err != nil {
		return TraceError(err)
	}

	if err := ce.compileSubroutineBody(subKind, subroutineName); err != nil {
		return TraceError(err)
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
			return TraceError(err)
		}
		// compiled a varDec, advance and try again or move on
	}
	return nil
}

// checks that token is identifier and returns it, or a error,
// and advances the tokenizer.
func (ce *CompilationEngine) getVarName() (string, error) {
	id, err := ce.jt.Identifier()
	if err != nil {
		return "", TraceError(fmt.Errorf("expected an %v for the varName", identifier))
	}
	if err := ce.advance(); err != nil {
		return "", err
	}
	return id, nil
}

// checkForSymbol checks that the passed symbol is currently being parsed
func (ce *CompilationEngine) checkForSymbol(sym string) error {
	if s, err := ce.jt.Symbol(); s != sym {
		if err != nil {
			return err
		}
		return fmt.Errorf("expected the %v \"%v\"", symbol, sym)
	}

	return nil
}

func (ce *CompilationEngine) checkForKeyword(kw string) error {
	if k, err := ce.jt.KeyWord(); k != kw {
		if err != nil {
			return TraceError(err)
		}
		return TraceError(fmt.Errorf("expected the %v \"%v\"", keyWord, kw))
	}

	return nil
}

// eatKeyword checks that the passed keyword is currently being parsed, and advances the tokenizer.
// Returns an error if the current token is not the passed keyword.
func (ce *CompilationEngine) eatKeyword(kw string) error {
	if err := ce.checkForKeyword(kw); err != nil {
		return TraceError(err)
	}
	if err := ce.advance(); err != nil {
		return TraceError(err)
	}
	return nil

}

// eatSymbol checks that the passed symbol is currently being parsed, and advances the tokenizer.
// Returns an error if the current token is not the passed symbol.
func (ce *CompilationEngine) eatSymbol(sym string) error {
	if err := ce.checkForSymbol(sym); err != nil {
		return err
	}
	if err := ce.advance(); err != nil {
		return err
	}
	return nil
}

// '(' ((type varName)(',' type varName)*)? ')'
//
// subKind is ('constructor' | 'function' | 'method')
func (ce *CompilationEngine) compileParameterList(subKind string) error {
	if err := ce.eatSymbol("("); err != nil {
		return TraceError(err)
	}

	if subKind == "method" {
		// The first argument of a method is always the object that the method is being called on.
		// Therefore, we need to add one to the parameter list.
		ce.st.Define("this", ce.className, KIND_ARG)
	}

	// While we have yet to hit the closing ")"
	for sym, _ := ce.jt.Symbol(); sym != ")"; sym, _ = ce.jt.Symbol() {
		// First token should be a type
		type_, err := ce.getType()
		if err != nil {
			return TraceError(err)
		}

		name, err := ce.getVarName()
		if err != nil {
			return TraceError(err)
		}

		// Define the parameter in the symbol table
		ce.st.Define(name, type_, KIND_ARG)

		// Now we should be at either a "," or the closing ")"
		if err := ce.checkForSymbol(","); err == nil {
			// if we're at the ",", advance and run through another round of the loop
			if err := ce.advance(); err != nil {
				return TraceError(err)
			}
			// but first check that we didn't bump immediately into a ")", which would be invalid
			if err := ce.checkForSymbol(")"); err == nil {
				return TraceError(fmt.Errorf("invalid syntax \",)\""))
			}
		}
		// Else we were at a ")", let the loop break and check for that
	}

	// Now we should be at the closing ")"
	if err := ce.eatSymbol(")"); err != nil {
		return TraceError(err)
	}

	return nil
}

// 'var' type varName (',' varName)* ';'
func (ce *CompilationEngine) compileVarDec() error {
	// check for var
	if err := ce.checkForKeyword("var"); err != nil {
		return TraceError(err)
	}

	// Advance and get type
	if err := ce.advance(); err != nil {
		return TraceError(err)
	}
	type_, err := ce.getType()
	if err != nil {
		return TraceError(err)
	}

	name, err := ce.getVarName()
	if err != nil {
		return TraceError(err)
	}

	// Add this var to symbol table
	ce.st.Define(name, type_, KIND_VAR)

	// Now see if there's a (',' varName)*
	for err = ce.eatSymbol(","); err == nil; err = ce.eatSymbol(",") {
		// (',' varName)*
		// check for "varName"
		name, err := ce.getVarName()
		if err != nil {
			return TraceError(err)
		}

		// Add this var to symbol table as well
		ce.st.Define(name, type_, KIND_VAR)
	}

	// Check for ";"
	if err := ce.eatSymbol(";"); err != nil {
		return TraceError(err)
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
			return TraceError(err)
		}

		// Each of the following compileX() functions should eat their final character before returning
		// in order to simplify switch case logic. There is no need to call advance() after compileX(),
		// since we can expect every compileX() below will have done so already.
		switch kw {
		case "let":
			if err := ce.compileLet(); err != nil {
				return TraceError(err)
			}
		case "if":
			if err := ce.compileIf(); err != nil {
				return TraceError(err)
			}
		case "while":
			if err := ce.compileWhile(); err != nil {
				return TraceError(err)
			}
		case "do":
			if err := ce.compileDo(); err != nil {
				return TraceError(err)
			}
		case "return":
			if err := ce.compileReturn(); err != nil {
				return TraceError(err)
			}
		default:
			return TraceError(fmt.Errorf("unexpected error in compileStatements, should be impossible"))
		}
	}

	return nil
}

// subroutineName '(' expressionList ')' |
// (className | varName) '.' subroutineName '(' expressionList ')'
//
// MUST not use temp 1
func (ce *CompilationEngine) compileSubroutineCall() error {
	var id1, id2, name string
	var numArgs uint
	var err error

	// Store subroutineName | className | varName in id1
	id1, err = ce.jt.Identifier()
	if err != nil {
		return TraceError(err)
	}
	if err := ce.advance(); err != nil {
		return TraceError(err)
	}

	var sym string
	sym, err = ce.jt.Symbol() // either a "." or a "("
	if err != nil {
		return TraceError(fmt.Errorf("expected a \".\" or \"(\""))
	}

	// If we're at a ".", then we're dealing with a className | varName for id1
	if sym == "." {
		// Eat the "."
		if err := ce.advance(); err != nil {
			return TraceError(err)
		}
		// store subroutineName in id2
		id2, err = ce.getIdentifier()
		if err != nil {
			return TraceError(err)
		}

	}

	// Now we should be at a "("
	if err = ce.eatSymbol("("); err != nil {
		return TraceError(err)
	}

	if id2 != "" {
		// (className | varName) '.' subroutineName
		// check if id1 is a varName
		_type, err := ce.st.TypeOf(id1)
		if err == nil {
			// We have a type for id1 in the symbol table, from which we can infer that it's a varName.
			// Ergo we're dealing with a method call here, meaning that we need to push the object that's
			// being called on to the stack, and then call the method on that object:
			ce.pushVarToStack(id1)
			// Increment the number of arguments we're going to use to call the method
			// by 1, since we're pushing "this" onto the stack as the first argument.
			numArgs += 1
			// And prepend the class name of the type of id1 to the method name
			// since all VM methods are "fully qualified".
			name = _type + "." + id2
		} else {
			// if not a varName, then we assume it's a function call
			// just call it as written, e.g. "Math.max()" -> "Math.max"
			name = id1 + "." + id2
		}
	} else {
		// subroutineName
		// This is a call in the form of func(), which means that it's
		// a method call of the class we're currently in. (Recall: functions
		// and constructors must be called with their full className.subroutineName()
		// syntax in Jack.)
		//
		// Since we must be within a method or constructor of this class, we can assume
		// that THIS is already set to the object that the method is being called on.
		// Therefore, we need to push that onto the stack as the first rgument of this method call.
		ce.cw.WritePush(SEG_POINTER, 0)
		// And we need to increment the number of arguments we're going to use to call the method
		// by 1, since we're pushing the object onto the stack as the first argument.
		numArgs += 1
		// Next we set the name to be the className of the class we're currently in
		// plus the subroutineName, e.g. "Square.draw", since all VM methods are "fully qualified".
		name = ce.className + "." + id1
	}

	numParams, err := ce.compileExpressionList()
	if err != nil {
		return TraceError(err)
	}

	ce.cw.WriteCall(name, numArgs+numParams)

	if err := ce.eatSymbol(")"); err != nil {
		return TraceError(err)
	}

	return nil
}

// 'do' subroutineCall ';'
// Eats it's own final character, so caller needn't immediately call advance() upon this function returning
func (ce *CompilationEngine) compileDo() error {
	if err := ce.advance(); err != nil {
		return TraceError(err)
	}
	if err := ce.compileSubroutineCall(); err != nil {
		return TraceError(err)
	}

	if err := ce.eatSymbol(";"); err != nil {
		return TraceError(err)
	}

	// after a `do funcCall()` there will be a value on top of the stack that we aren't using, so pop that to temp
	ce.cw.WritePop(SEG_TEMP, 0)
	return nil
}

// 'let' varName ('[' expression ']')? '=' expression ';'
// Eats it's own final character, so caller needn't immediately call advance() upon this function returning
func (ce *CompilationEngine) compileLet() error {
	var err error

	if err = ce.checkForKeyword("let"); err != nil {
		return TraceError(err)
	}

	// advance and get varName
	if err = ce.advance(); err != nil {
		return TraceError(err)
	}
	varName, err := ce.getVarName()
	if err != nil {
		return TraceError(err)
	}

	// check for '['
	var sym string
	if sym, err = ce.jt.Symbol(); err != nil {
		return TraceError(err)
	}
	settingArrayElement := sym == "["

	// If we're at a "[", then we must be dealing with an array.
	// In this case, we want to push the base address of the array (`varName`) onto the stack,
	// and then add the index to that base address to get the address of the element we want to access.
	if settingArrayElement {
		if err := ce.setThatForArrayAccess(varName); err != nil {
			return TraceError(err)
		}
		// push the address of the array element we're setting onto the stack
		ce.cw.WritePush(SEG_POINTER, 1)
		// pop it over to temp for later use
		ce.cw.WritePop(SEG_TEMP, 1)
	}

	// demand, compile, and eat '='
	if err = ce.eatSymbol("="); err != nil {
		return TraceError(err)
	}

	// compile expression, whose result will wind up on the top of the stack
	if err = ce.compileExpression(); err != nil {
		return TraceError(err)
	}

	if settingArrayElement {
		// push the address of the array element we're setting, which we stored in temp
		// earlier in this function, onto the stack
		ce.cw.WritePush(SEG_TEMP, 1)
		// set THAT to the address we just pushed onto the stack
		ce.cw.WritePop(SEG_POINTER, 1)
		// pop the expression result into the THAT segment
		ce.cw.WritePop(SEG_THAT, 0)
	} else {
		// pop the expression result into its corresponding variable
		kind, err := ce.st.KindOf(varName)
		if err != nil {
			return TraceError(err)
		}
		index, err := ce.st.IndexOf(varName)
		if err != nil {
			return TraceError(err)
		}

		ce.cw.WritePop(kindToSegment(kind), index)
	}

	// check and eat ";"
	if err = ce.eatSymbol(";"); err != nil {
		return TraceError(err)
	}

	return nil
}

// 'while '(' expression ')' '{' statements '}'
// eats its own final token before returning, so calling function can expect to already be on the next token
func (ce *CompilationEngine) compileWhile() error {
	// generate unique start and end labels
	uuid := fmt.Sprintf("%v", ce.whileId)
	ce.whileId++
	startLabel := "while_start_" + ce.className + "_" + uuid
	endLabel := "while_end_" + ce.className + "_" + uuid

	// write the start label
	ce.cw.WriteLabel(startLabel)

	if err := ce.advance(); err != nil {
		return TraceError(err)
	}
	if err := ce.checkForSymbol("("); err != nil {
		return TraceError(err)
	}

	// compute the condition:
	// if true then -1 will be on top of the stack,
	// if false then 0 will be on top of the stack
	if err := ce.advance(); err != nil {
		return TraceError(err)
	}
	if err := ce.compileExpression(); err != nil {
		return TraceError(err)
	}

	// Now bit-wise not whatever is on top of the stack and if-goto the endLabel
	// if-goto only jumps if value on top of stack is nonzero, so in the case where
	// the condition was true, -1 will be bit-wise not-ed to 0, and so the if-goto will
	// be ignored. In the case where the condition was false, 0 will be bit-wise not-ed
	// to -1, and the if-goto will cause a goto that breaks out of the loop.
	ce.cw.WriteArithmetic(COM_NOT)
	ce.cw.WriteIf(endLabel)

	if err := ce.checkForSymbol(")"); err != nil {
		return TraceError(err)
	}

	if err := ce.advance(); err != nil {
		return TraceError(err)
	}
	if err := ce.checkForSymbol("{"); err != nil {
		return TraceError(err)
	}

	// compile the loop's internal logic
	if err := ce.advance(); err != nil {
		return TraceError(err)
	}
	if err := ce.compileStatements(); err != nil {
		return TraceError(err)
	}

	// end of loop, jump back to the beggining
	ce.cw.WriteGoto(startLabel)

	if err := ce.checkForSymbol("}"); err != nil {
		return TraceError(err)
	}
	if err := ce.advance(); err != nil {
		return TraceError(err)
	}

	// place the end label at the bottom of the loop in order to escape it
	ce.cw.WriteLabel(endLabel)

	return nil
}

// 'return' expression? ';'
// Eats it's own final character, so caller needn't immediately call advance() upon this function returning.
// According to the book:
// - VM functions corresponding to void Jack methods and functions must return the constant 0 as their return value.
// - When translating a do sub statement where sub is a void method or function, the caller of the corresponding VM function must pop (and ignore) the returned value (which is always the constant 0).
// However I'm skeptical as to whether the first bullet is strictly necessary based on walking through the stack model and not seeing any problems with just returning whatever junk happens to be on top
// of the stack rather than returning 0. I am going to ignore this requirement for now in order to see if it surfaces some bug that identifies why its a requirement, and have also asked a question
// on the forum here: http://nand2tetris-questions-and-answers-forum.32033.n3.nabble.com/Do-void-functions-truly-need-to-return-0-td4035927.html
func (ce *CompilationEngine) compileReturn() error {
	// advance and check for ';'
	if err := ce.advance(); err != nil {
		return TraceError(err)
	}
	// if not ';', should be an expression
	if sym, _ := ce.jt.Symbol(); sym != ";" {
		if err := ce.compileExpression(); err != nil {
			return TraceError(err)
		}
	}
	// now compile the ';'
	if err := ce.checkForSymbol(";"); err != nil {
		return TraceError(err)
	}
	// eat own final character
	if err := ce.advance(); err != nil {
		return TraceError(err)
	}
	ce.cw.WriteReturn()
	return nil
}

// 'if' '(' expression ')' '{' statements '}'
// ('else' '{' statements '}')?
// Eats it's own final character, so caller needn't immediately call advance() upon this function returning
func (ce *CompilationEngine) compileIf() error {
	// generate unique start and end labels
	uuid := fmt.Sprintf("%v", ce.ifId)
	ce.ifId++
	elseLabel := "else_" + ce.className + "_" + uuid
	endLabel := "if_else_end_" + ce.className + "_" + uuid

	// this chunk of code is used in regular if and if-else, so abstracted into an internal
	// function to avoid having it written twice
	compileStatementsSubsection := func() error {
		if err := ce.eatSymbol("{"); err != nil {
			return TraceError(err)
		}

		if err := ce.compileStatements(); err != nil {
			return TraceError(err)
		}

		if err := ce.eatSymbol("}"); err != nil {
			return TraceError(err)
		}

		return nil
	}

	if err := ce.advance(); err != nil {
		return TraceError(err)
	}
	if err := ce.checkForSymbol("("); err != nil {
		return TraceError(err)
	}

	// compute condition
	if err := ce.advance(); err != nil {
		return TraceError(err)
	}
	if err := ce.compileExpression(); err != nil {
		return TraceError(err)
	}

	// negate the conditional result on the top of the stack, so that the subsequent
	// if-goto only jumps to the else label if the condition was false:
	// - in the case that conditional evaluates to false (0), it gets negated to true (-1), which causes the if-goto to execute a jump to the elseLabel
	// - in the case that the condition evaluates to true (-1), it gets negated false (0), which means the if-goto doesn't jump, and the if statement is exectuted
	//   (which subsequently skips the else statement by jumping to the end)
	ce.cw.WriteArithmetic(COM_NOT)
	ce.cw.WriteIf(elseLabel)

	// compileExpression loops us to next token so no need to call advance()
	if err := ce.eatSymbol(")"); err != nil {
		return TraceError(err)
	}

	if err := compileStatementsSubsection(); err != nil {
		return TraceError(err)
	}

	// if condition was true and statement was executed, so skip the else by jumping to the end label
	ce.cw.WriteGoto(endLabel)

	// Write the else label. The else label can be used in this logic even in the case
	// that an if statement has no corresponding else
	// (somewhat complicated story, see Figure 8.1 in the book and walk through the logic)
	ce.cw.WriteLabel(elseLabel)

	if ce.jt.TokenType() == keyWord {
		if kw, _ := ce.jt.KeyWord(); kw == "else" {
			if err := ce.eatKeyword("else"); err != nil {
				return TraceError(err)
			}

			if err := compileStatementsSubsection(); err != nil {
				return TraceError(err)
			}
		}
	}

	// write the end label
	ce.cw.WriteLabel(endLabel)

	return nil
}

// isOp returns true if the passed symbol is a valid op, else returns false
func isOp(sym string) bool {
	return (sym == "+" || sym == "-" || sym == "*" || sym == "/" || sym == "&" || sym == "|" || sym == "<" || sym == ">" || sym == "=")
}

// MUST not use temp 1
func (ce *CompilationEngine) compileOp(sym string) {
	switch sym {
	case "+":
		ce.cw.WriteArithmetic(COM_ADD)
	case "-":
		ce.cw.WriteArithmetic(COM_SUB)
	case "*":
		ce.cw.WriteCall("Math.multiply", 2)
	case "/":
		ce.cw.WriteCall("Math.divide", 2)
	case "&":
		ce.cw.WriteArithmetic(COM_AND)
	case "|":
		ce.cw.WriteArithmetic(COM_OR)
	case "<":
		ce.cw.WriteArithmetic(COM_LT)
	case ">":
		ce.cw.WriteArithmetic(COM_GT)
	case "=":
		ce.cw.WriteArithmetic(COM_EQ)
	default:
		panic("invalid operator")
	}
}

// term (op term)*
// Loops until a non (op term) is found, so caller should expect to be at the next token when this function returns.
//
// MUST not use temp 1
func (ce *CompilationEngine) compileExpression() error {
	// compiles the first term and pushes its result onto the top of the stack
	if err := ce.compileTerm(); err != nil {
		return TraceError(err)
	}

	// (op term)*
	for sym, err := ce.jt.Symbol(); isOp(sym) && err == nil; sym, err = ce.jt.Symbol() {
		// op term
		if err = ce.advance(); err != nil {
			return TraceError(err)
		}
		// first term is on top of the stack,
		// now compile the next term so that's on top of the first
		if err = ce.compileTerm(); err != nil {
			return TraceError(err)
		}
		// then apply the operation to the two terms at the top of the stack
		ce.compileOp(sym)
	}

	return nil
}

// uintegerConstant | stringConstant | keywordConstant |
// varName | varName '[' expression ']' | subroutineCall |
// '(' expression ')' | unaryOp term
//
// MUST not use temp 1
func (ce *CompilationEngine) compileTerm() error {
	if ce.jt.TokenType() == intConst {
		intVal, err := ce.getIntVal()
		if err != nil {
			return TraceError(err)
		}
		ce.cw.WritePush(SEG_CONST, intVal)
		return nil
	} else if ce.jt.TokenType() == strConst {
		strVal, err := ce.jt.StringVal()
		if err != nil {
			return TraceError(err)
		}
		// Allocate memory for the string
		ce.cw.WritePush(SEG_CONST, uint(len(strVal)))
		ce.cw.WriteCall("String.new", 1) // create the string object, returned to top of stack
		// Write the string to memory
		for _, c := range strVal {
			// push the character to the stack
			ce.cw.WritePush(SEG_CONST, uint(c))
			// call String.appendChar. this appends the next char
			// and returns the string object to the top of the stack
			ce.cw.WriteCall("String.appendChar", 2)
		}
		ce.advance()

		return nil
	} else if ce.jt.TokenType() == keyWord {
		kw, err := ce.getKeyWord()
		if err != nil {
			return TraceError(err)
		}
		switch kw {
		case "true":
			ce.cw.WritePush(SEG_CONST, 1)
			ce.cw.WriteArithmetic(COM_NEG)
		case "false":
			ce.cw.WritePush(SEG_CONST, 0)
		case "this":
			ce.cw.WritePush(SEG_POINTER, 0)
		case "null":
			ce.cw.WritePush(SEG_CONST, 0)
		default:
			return TraceError(fmt.Errorf("term keyWord must be one of \"true\", \"false\", \"null\", or \"this\""))
		}
		return nil
	} else if ce.jt.TokenType() == symbol {
		// '(' expression ')' | unaryOp term
		sym, err := ce.jt.Symbol()
		if err != nil {
			return TraceError(err)
		}
		if err := ce.advance(); err != nil {
			return TraceError(err)
		}
		if sym == "(" {
			// '(' expression ')'
			if err := ce.compileExpression(); err != nil {
				return TraceError(err)
			}
			if err := ce.eatSymbol(")"); err != nil {
				return TraceError(err)
			}
		} else if sym == "-" {
			// -term
			// push whatever is being negated to the top of the stack
			if err := ce.compileTerm(); err != nil {
				return TraceError(err)
			}

			// Negate it
			ce.cw.WriteArithmetic(COM_NEG)
		} else if sym == "~" {
			// ~term
			// push whatever is being bit-wise not-ed to the top of the stack
			if err := ce.compileTerm(); err != nil {
				return TraceError(err)
			}

			// Bit-wise not it
			ce.cw.WriteArithmetic(COM_NOT)
		} else {
			return TraceError(fmt.Errorf("invalid symbol in term, symbol must be one of \"(\" or \"-\" or \"~\""))
		}
		return nil
	} else if ce.jt.TokenType() == identifier {
		// varName | varName '[' expression ']' | subroutineCall
		var err error
		var peeked byte
		peeked, err = ce.jt.Peek()
		if err != nil {
			return TraceError(err)
		}
		if peeked == '.' || peeked == '(' {
			return ce.compileSubroutineCall()
		} else if peeked == '[' {
			// We're terminating at an array access
			varName, err := ce.getVarName()
			if err != nil {
				return TraceError(err)
			}
			// Set the THAT segment to the address of the element in the array that is being accessed
			if err := ce.setThatForArrayAccess(varName); err != nil {
				return TraceError(err)
			}
			// Push the value of the array element onto the stack
			ce.cw.WritePush(SEG_THAT, 0)
		} else {
			// Else we are at a variable
			varName, err := ce.getVarName()
			if err != nil {
				return TraceError(err)
			}

			err = ce.pushVarToStack(varName)
			if err != nil {
				return TraceError(err)
			}
		}
		return nil
	} else {
		return TraceError(fmt.Errorf("unknown error"))
	}
}

// Sets the THAT segment to the address of the element in the array that is being accessed.
//
// Expects the tokenizer to be at the beginning of:
// '[' expression ']'
// `varName` must be the name of the array being accessed,
// as in `varName[expression]`.
//
// Returns having eaten the closing ']' character.
func (ce *CompilationEngine) setThatForArrayAccess(varName string) error {
	// confirm that we're dealing with an Array
	type_, err := ce.st.TypeOf(varName)
	if err != nil {
		return TraceError(err)
	}
	if type_ != "Array" {
		return TraceError(fmt.Errorf("expected an array"))
	}
	// push the base address of the Array onto the stack
	err = ce.pushVarToStack(varName)
	if err != nil {
		return TraceError(err)
	}

	// eat the '['
	if err = ce.advance(); err != nil {
		return TraceError(err)
	}
	// compile the expression, it's result will then be on top of the stack
	if err = ce.compileExpression(); err != nil {
		return TraceError(err)
	}
	if err := ce.eatSymbol("]"); err != nil {
		return TraceError(err)
	}
	// add the expression result (the index) to the base address of the Array to get the address of the element we want to set
	ce.cw.WriteArithmetic(COM_ADD)
	// pop the result into the pointer segment to set THAT
	ce.cw.WritePop(SEG_POINTER, 1)

	return nil
}

func kindToSegment(kind Kind) Segment {
	switch kind {
	case KIND_STATIC:
		return SEG_STATIC
	case KIND_VAR:
		return SEG_LOCAL
	case KIND_ARG:
		return SEG_ARG
	case KIND_FIELD:
		return SEG_THIS
	default:
		panic(fmt.Sprintf("kindToSegment not implemented for kind %v", kind))
	}
}

// (expression (',' expression)* )?
// Caller should expect to be at the next token when this function returns.
// Returns the number of ',' separated expressions that were compiled
//
// MUST not use temp 1
func (ce *CompilationEngine) compileExpressionList() (uint, error) {
	var nArgs uint

	// Terminating condition: check if we are at a closing parenthesis
	sym, _ := ce.jt.Symbol()
	if ce.jt.TokenType() == symbol && sym == ")" {
		// The next token was a closing parenthesis; the caller is expected to account for it
		return nArgs, nil
	}

	// If we're not at a closing parenthesis, then we're at an expression
	if err := ce.compileExpression(); err != nil {
		return nArgs, TraceError(err)
	}
	// We've compiled one expression, so increment nArgs
	nArgs++

	// check for further expressions in the list
	for sym, _ := ce.jt.Symbol(); sym == ","; sym, _ = ce.jt.Symbol() {
		if err := ce.eatSymbol(","); err != nil {
			// Should be impossible to get here, since we already checked for a ","
			panic("impossible")
		}
		if err := ce.compileExpression(); err != nil {
			return nArgs, TraceError(err)
		}
		nArgs++
	}

	return nArgs, nil
}

// pushVarToStack pushes the variable with the given name onto the stack.
func (ce *CompilationEngine) pushVarToStack(varName string) error {
	kind, err := ce.st.KindOf(varName)
	if err != nil {
		return TraceError(err)
	}

	index, err := ce.st.IndexOf(varName)
	if err != nil {
		return TraceError(err)
	}

	ce.cw.WritePush(kindToSegment(kind), index)
	return nil
}
