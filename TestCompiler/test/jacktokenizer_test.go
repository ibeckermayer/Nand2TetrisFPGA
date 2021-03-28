package test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"testing"

	"os"

	"github.com/ibeckermayer/Nand2TetrisFPGA/Compiler/pkg/compiler"
)

func runTokenizerTest(jackFilePath string, t *testing.T) {
	correctFilePath := fmt.Sprintf("%vT.xml", jackFilePath[0:len(jackFilePath)-len(".jack")])
	tmpFilePath := fmt.Sprintf("%v_out.xml", correctFilePath[0:len(correctFilePath)-len(".xml")])

	// Create jt
	jt, err := compiler.NewJackTokenizer(jackFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// Create tmp file to write to
	f, err := os.Create(tmpFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Pass tmp file to xml encoder
	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")

	// Open <tokens> tag
	enc.EncodeToken(
		xml.StartElement{Name: xml.Name{Space: "", Local: "tokens"}, Attr: []xml.Attr{}})

	// Walk through file and encode jt state as XML
	for err := jt.Advance(); err == nil && jt.HasMoreTokens(); err = jt.Advance() {
		if err != nil {
			if err != nil {
				t.Fatal(err)
			}
		}
		jt.MarshalXML(enc,
			xml.StartElement{Name: xml.Name{Space: "not used", Local: "not used"}, Attr: []xml.Attr{}})
	}

	// Close </tokens> tag and flush encoder to tmp file
	enc.EncodeToken(xml.EndElement{Name: xml.Name{Space: "", Local: "tokens"}})
	enc.Flush()
	f.WriteString("\n")

	// Check that the xml file we just wrote is what we expect it to be
	b1, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := ioutil.ReadFile(correctFilePath)
	if err != nil {
		t.Fatal(err)
	}

	if !(bytes.Equal(b1, b2)) {
		t.Fatal("Files weren't equal!")
	}
}

func TestTokenizerArrayTestMain(t *testing.T) {
	runTokenizerTest("./ArrayTest/Main.jack", t)
}

func TestTokenizerExpressionLessSquareMain(t *testing.T) {
	runTokenizerTest("./ExpressionLessSquare/Main.jack", t)
}

func TestTokenizerExpressionLessSquareSquare(t *testing.T) {
	runTokenizerTest("./ExpressionLessSquare/Square.jack", t)
}

func TestTokenizerExpressionLessSquareSquareGame(t *testing.T) {
	runTokenizerTest("./ExpressionLessSquare/SquareGame.jack", t)
}

func TestTokenizerSquareMain(t *testing.T) {
	runTokenizerTest("./Square/Main.jack", t)
}

func TestTokenizerSquareSquare(t *testing.T) {
	runTokenizerTest("./Square/Square.jack", t)
}

func TestTokenizerSquareSquareGame(t *testing.T) {
	runTokenizerTest("./Square/SquareGame.jack", t)
}
