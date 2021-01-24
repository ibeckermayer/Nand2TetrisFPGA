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
	jt     *JackTokenizer // A tokenizer set up to tokenize the file we want to compile
	out    *os.File       // The output file
	xmlEnc *xml.Encoder   // The xml encoder for testing
}

// NewCompilationEngine creates a new CompilationEngine.
func NewCompilationEngine(filename string) (*CompilationEngine, error) {
	jt, err := NewJackTokenizer(filename)
	if err != nil {
		return nil, err
	}

	xmlOutputFile, err := os.Create(fmt.Sprintf("%v_out.xml", filename[0:len(filename)-len(".jack")]))
	if err != nil {
		return nil, err
	}
	xmlEnc := xml.NewEncoder(xmlOutputFile)
	xmlEnc.Indent("", "  ")

	return &CompilationEngine{
		jt:     jt,
		out:    xmlOutputFile,
		xmlEnc: xmlEnc,
	}, nil
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

// CompileClass should be called immediately after the compilationEngine is created.
// 'class' className '{' classVarDec* subroutineDec* '}'
func (ce *CompilationEngine) CompileClass() error {
	defer ce.out.Close()
	defer ce.xmlEnc.Flush()
	// Advance and check that first token is "class"
	ce.jt.Advance()
	if kw, err := ce.jt.KeyWord(); err != nil || kw != "class" {
		return SyntaxError(fmt.Errorf("expected keyword \"class\""))
	}

	// Found class, open tag
	ce.openXMLTag("class")        // <class>
	defer ce.closeXMLTag("class") // defer closing </class>
	ce.marshaljt()                // <keyword> class </keyword>

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
		return SyntaxError(fmt.Errorf("expected the symbol \"{\" but got \"%v\" instead", sym))
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

		ce.CompileClassVarDec()

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

		ce.CompileSubroutine()

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

// TODO: You should start from here!!!
func (ce *CompilationEngine) CompileClassVarDec() {
	// panic("not implemented") // TODO: Implement
	return
}

func (ce *CompilationEngine) CompileSubroutine() {
	// panic("not implemented") // TODO: Implement
	return
}

func (ce *CompilationEngine) CompileParameterList() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) CompileVarDec() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) CompileStatements() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) CompileDo() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) CompileLet() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) CompileWhile() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) DompileReturn() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) CompileIf() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) CompileExpression() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) CompileTerm() {
	panic("not implemented") // TODO: Implement
}

func (ce *CompilationEngine) CompileExpressionList() {
	panic("not implemented") // TODO: Implement
}
