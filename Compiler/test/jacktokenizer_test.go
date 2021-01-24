package test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"testing"

	"os"

	jtz "github.com/ibeckermayer/Nand2TetrisFPGA/Compiler/pkg/compiler"
)

func fatalize(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func runTokenizerTest(jackFilePath string, t *testing.T) {
	correctFilePath := fmt.Sprintf("%vT.xml", jackFilePath[0:len(jackFilePath)-len(".jack")])
	tmpFilePath := fmt.Sprintf("%v_outT.xml", correctFilePath[0:len(correctFilePath)-len(".xml")])

	// Create jt
	jt, err := jtz.NewJackTokenizer(jackFilePath)
	fatalize(err, t)

	// Create tmp file to write to
	f, err := os.Create(tmpFilePath)
	fatalize(err, t)
	defer func() {
		f.Close()
		os.Remove(tmpFilePath)
	}()

	// Pass tmp file to xml encoder
	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")

	// Open <tokens> tag
	enc.EncodeToken(
		xml.StartElement{Name: xml.Name{Space: "", Local: "tokens"}, Attr: []xml.Attr{}})

	// Walk through file and encode jt state as XML
	for jt.Advance(); jt.HasMoreTokens(); jt.Advance() {
		jt.MarshalXML(enc,
			xml.StartElement{Name: xml.Name{Space: "not used", Local: "not used"}, Attr: []xml.Attr{}})
	}

	// Close </tokens> tag and flush encoder to tmp file
	enc.EncodeToken(xml.EndElement{Name: xml.Name{Space: "", Local: "tokens"}})
	enc.Flush()
	f.WriteString("\n")

	// Check that the xml file we just wrote is what we expect it to be
	b1, err := ioutil.ReadFile(tmpFilePath)
	fatalize(err, t)
	b2, err := ioutil.ReadFile(correctFilePath)
	fatalize(err, t)

	if !(bytes.Equal(b1, b2)) {
		t.Fatal("Files weren't equal!")
	}
}

func TestArrayTestMain(t *testing.T) {
	runTokenizerTest("./ArrayTest/Main.jack", t)
}

func TestExpressionLessSquareMain(t *testing.T) {
	runTokenizerTest("./ExpressionLessSquare/Main.jack", t)
}

func TestExpressionLessSquareSquare(t *testing.T) {
	runTokenizerTest("./ExpressionLessSquare/Square.jack", t)
}

func TestExpressionLessSquareSquareGame(t *testing.T) {
	runTokenizerTest("./ExpressionLessSquare/SquareGame.jack", t)
}

func TestSquareMain(t *testing.T) {
	runTokenizerTest("./Square/Main.jack", t)
}

func TestSquareSquare(t *testing.T) {
	runTokenizerTest("./Square/Square.jack", t)
}

func TestSquareSquareGame(t *testing.T) {
	runTokenizerTest("./Square/SquareGame.jack", t)
}
