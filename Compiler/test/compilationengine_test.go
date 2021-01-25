package test

import (
	"testing"

	"github.com/ibeckermayer/Nand2TetrisFPGA/Compiler/pkg/compiler"
)

func runCeTest(jackFilePath string, t *testing.T) {
	ce := &compiler.CompilationEngine{JackFilePath: jackFilePath}

	err := ce.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCeExpressionLessSquareMain(t *testing.T) {
	runCeTest("./ExpressionLessSquare/Main.jack", t)
}
