package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ibeckermayer/Nand2TetrisFPGA/Compiler/pkg/compiler"
)

// isJack determines whether a string ends in ".jack" or not
func isJack(in string) bool {
	split := strings.Split(in, ".")
	return split[len(split)-1] == "jack"
}

// jackFilter filters a list of strings down to only those which end in ".jack"
func jackFilter(ss []string) (ret []string) {
	for _, s := range ss {
		if isJack(s) {
			ret = append(ret, s)
		}
	}
	return
}

func main() {
	// First argument should be .jack file or directory of .jack files. Currently
	// do not support multiple directory builds.
	if len(os.Args) < 2 {
		fmt.Printf("program requires the first argument be a path to .jack file or directory that contains at least one .jack file")
		os.Exit(1)
	}
	toCompile := os.Args[1]
	var files []string

	// try to list all the files in a directory by the name of the first command line argument,
	// or just list the single jack file argument
	if err := filepath.Walk(toCompile, func(path string, info os.FileInfo, err error) error {
		if path != toCompile && info.IsDir() {
			// only walk a single level of directories
			return filepath.SkipDir
		}
		if isJack(path) {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
	if len(files) == 0 {
		log.Printf("invalid compilation input: \"%v\"; first argument must be either a .jack file or a directory with at least one .jack file in it", toCompile)
		os.Exit(1)
	}

	// for each jack file
	for _, filePath := range files {
		ce, err := compiler.NewCompilationEngine(filePath)
		err = ce.Run()
		if err != nil {
			log.Printf("Error running the compilation engine: %v", err)
			os.Exit(1)
		}
	}
}
