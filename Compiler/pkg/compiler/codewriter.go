package compiler

import (
	"fmt"
	"os"
)

// CodeWriter is the part responsible for writing vm code
type CodeWriter struct {
	// The output file
	outputFile *os.File
	// className is the class name being compiled,
	// set at compileClass by compilation engine
	className string
}

// NewCodeWriter creates a new CodeWriter
func NewCodeWriter(jackFilePath string) (*CodeWriter, error) {
	// Create the output file
	outputFile, err := os.Create(fmt.Sprintf("%v.vm", jackFilePath[0:len(jackFilePath)-len(".jack")]))
	if err != nil {
		return nil, err
	}

	return &CodeWriter{outputFile, ""}, nil
}

// Close closes the output file
func (cw *CodeWriter) Close() error {
	return cw.outputFile.Close()
}

func (cw *CodeWriter) writeln(format string, a ...interface{}) {
	cw.outputFile.WriteString(fmt.Sprintf(format, a...))
}

func (cw *CodeWriter) WriteFunction(name string, nLocals uint) {
	cw.writeln("function %v.%v %v\n", cw.className, name, nLocals)
}
