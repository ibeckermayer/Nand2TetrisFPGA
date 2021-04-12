package compiler

import (
	"fmt"
	"os"
)

// CodeWriter is the part responsible for writing vm code
type CodeWriter struct {
	// The output file
	outputFile *os.File
}

// NewCodeWriter creates a new CodeWriter
func NewCodeWriter(jackFilePath string) (*CodeWriter, error) {
	// Create the output file
	outputFile, err := os.Create(fmt.Sprintf("%v.vm", jackFilePath[0:len(jackFilePath)-len(".jack")]))
	if err != nil {
		return nil, err
	}

	return &CodeWriter{outputFile}, nil
}

// Close closes the output file
func (cw *CodeWriter) Close() error {
	return cw.outputFile.Close()
}

func (cw *CodeWriter) writeln(format string, a ...interface{}) {
	cw.outputFile.WriteString(fmt.Sprintf(format+"\n", a...))
}

// WriteFunction writes a vm function command
func (cw *CodeWriter) WriteFunction(name string, nLocals int) {
	cw.writeln("function %v %v", name, nLocals)
}

// WriteCall writes a vm call command
func (cw *CodeWriter) WriteCall(name string, nArgs int) {
	cw.writeln("call %v %v", name, nArgs)
}

// WritePush writes a vm push command
func (cw *CodeWriter) WritePush(segment Segment, index int) {
	cw.writeln("push %v %v", segment, index)
}

// Write Pop writes a vm pop command
func (cw *CodeWriter) WritePop(segment Segment, index int) {
	cw.writeln("pop %v %v", segment, index)
}

// WriteArithmetic writes a vm arithmetic command
func (cw *CodeWriter) WriteArithmetic(cmd Command) {
	cw.writeln("%v", cmd)
}

// WriteReturn writes a vm return command
func (cw *CodeWriter) WriteReturn() {
	cw.writeln("return")
}

type Segment string

const (
	SEG_CONST   Segment = "constant"
	SEG_ARG     Segment = "argument"
	SEG_LOCAL   Segment = "local"
	SEG_STATIC  Segment = "static"
	SEG_THIS    Segment = "this"
	SEG_THAT    Segment = "that"
	SEG_POINTER Segment = "pointer"
	SEG_TEMP    Segment = "temp"
)

type Command string

const (
	COM_ADD Command = "add"
	COM_SUB Command = "sub"
	COM_NEG Command = "neg"
	COM_EQ  Command = "eq"
	COM_GT  Command = "gt"
	COM_LT  Command = "lt"
	COM_AND Command = "and"
	COM_OR  Command = "or"
	COM_NOT Command = "not"
)
