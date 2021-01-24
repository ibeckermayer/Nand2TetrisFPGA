package test

import (
	"testing"

	"github.com/ibeckermayer/Nand2TetrisFPGA/Compiler/pkg/compiler"
)

func runCeTest(jackFilePath string, t *testing.T) {
	// correctFilePath := fmt.Sprintf("%v.xml", jackFilePath[0:len(jackFilePath)-len(".jack")])

	ce, err := compiler.NewCompilationEngine(jackFilePath)
	if err != nil {
		fatalize(err, t)
	}

	err = ce.CompileClass()
	if err != nil {
		fatalize(err, t)
	}
}

func TestCeExpressionLessSquareMain(t *testing.T) {
	runCeTest("./ExpressionLessSquare/Main.jack", t)
}
