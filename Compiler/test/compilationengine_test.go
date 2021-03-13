package test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ibeckermayer/Nand2TetrisFPGA/Compiler/pkg/compiler"
)

func runCeTest(jackFilePath string, t *testing.T) {
	ce := &compiler.CompilationEngine{JackFilePath: jackFilePath}

	err := ce.Run()
	if err != nil {
		t.Fatal(err)
	}

	correctFilePath := strings.Replace(jackFilePath, ".jack", ".xml", -1)
	outFilePath := strings.Replace(jackFilePath, ".jack", "_out.xml", -1)

	b1, err := ioutil.ReadFile(correctFilePath)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := ioutil.ReadFile(outFilePath)
	if err != nil {
		t.Fatal(err)
	}

	if !(bytes.Equal(b1, b2)) {
		t.Fatal("Files weren't equal!")
	}
}

func TestCeExpressionLessSquareMain(t *testing.T) {
	runCeTest("./ExpressionLessSquare/Main.jack", t)
}

func TestCeExpressionLessSquareSquare(t *testing.T) {
	runCeTest("./ExpressionLessSquare/Square.jack", t)
}

func TestCeExpressionLessSquareSquareGame(t *testing.T) {
	runCeTest("./ExpressionLessSquare/SquareGame.jack", t)
}

func TestCeArrayTestMain(t *testing.T) {
	runCeTest("./ArrayTest/Main.jack", t)
}

func TestCeSquareMain(t *testing.T) {
	runCeTest("./Square/Main.jack", t)
}

func TestCeSquareSquare(t *testing.T) {
	runCeTest("./Square/Square.jack", t)
}

func TestCeSquareSquareGame(t *testing.T) {
	runCeTest("./Square/SquareGame.jack", t)
}
