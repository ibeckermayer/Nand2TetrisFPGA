package compiler

import "fmt"

type Kind string

const (
	KIND_STATIC Kind = "static"
	KIND_FIELD  Kind = "field"
	KIND_ARG    Kind = "arg"
	KIND_VAR    Kind = "var"
	KIND_NONE   Kind = "none"
)

// Entry is an entry in the symbol table
type Entry struct {
	// Name is the identifier name
	Name string
	// Type is the identifier type
	Type string
	// Kind is the kind of the identifier (see constants of prefix "KIND_")
	Kind Kind
	// Index is the index assigned to the identifier in the current scope
	Index uint
}

// symTable is the actual table mapping symbol names to their corresponding entry
type symTable map[string]Entry

// SymbolTable associates the identifier names found in the program with identifier properties needed for compilation.
// The SymbolTable has two nested scopes: class/subroutine
type SymbolTable struct {
	// table is where the class table and subroutine tables themselves are stored, in a map from Kind to *symTable.
	// Kind can be used to distinguish between class and subroutine scopes: KIND_STATIC and KIND_FIELD belong to the
	// class table, and KIND_ARG and KIND_VAR belong to the subroutine table. Storing the tables like this allows for
	// simpler and more concise class methods
	table map[Kind]*symTable
	// varCount is a map between each Kind and the number of variables of the given Kind defined in the current scope
	varCount map[Kind]uint
}

// NewSymbolTable creates a new empty symbol table
func NewSymbolTable() *SymbolTable {
	classTable := make(symTable)
	subroutineTable := make(symTable)

	return &SymbolTable{
		table: map[Kind]*symTable{
			KIND_STATIC: &classTable,
			KIND_FIELD:  &classTable,
			KIND_ARG:    &subroutineTable,
			KIND_VAR:    &subroutineTable,
		},

		varCount: map[Kind]uint{
			KIND_STATIC: 0,
			KIND_FIELD:  0,
			KIND_ARG:    0,
			KIND_VAR:    0,
		},
	}
}

// StartSubroutine starts a new subroutine, i.e. resets the SubroutineTable
func (s *SymbolTable) StartSubroutine() {
	subroutineTable := make(symTable)
	s.table[KIND_ARG] = &subroutineTable
	s.table[KIND_VAR] = &subroutineTable
	s.varCount[KIND_ARG] = 0
	s.varCount[KIND_VAR] = 0
}

// Define defines a new entry in the symbol table
func (s *SymbolTable) Define(name, type_ string, kind Kind) error {
	table := *(s.table[kind])

	// check if this symbol is already defined in the given scope
	_, ok := table[name]
	if ok {
		return fmt.Errorf("attempted redefinition of symbol: %v", name)
	}

	// add entry
	table[name] = Entry{
		Name:  name,
		Type:  type_,
		Kind:  kind,
		Index: s.varCount[kind],
	}

	// increment varCount for this type
	s.varCount[kind]++

	return nil
}

// VarCount returns the number of variables of the given kind already defined in the current scope
func (s *SymbolTable) VarCount(kind Kind) uint {
	return s.varCount[kind]
}

// Returns the kind of the named identifier in the given scope
func (s *SymbolTable) KindOf(name string) Kind {
	var kind Kind
	subroutineTable := *s.table[KIND_ARG]
	classTable := *s.table[KIND_STATIC]

	// Check subroutine scope first
	kind = subroutineTable[name].Kind
	if kind != "" {
		return kind
	}

	// If nothing was found, check class scope
	kind = classTable[name].Kind
	if kind != "" {
		return kind
	}

	// If nothing was found in either table return KIND_NONE
	return KIND_NONE
}

// Returns the type of the named identifier in the given scope
func (s *SymbolTable) TypeOf(name string) (string, error) {
	var type_ string
	subroutineTable := *s.table[KIND_ARG]
	classTable := *s.table[KIND_STATIC]

	// Check subroutine scope first
	type_ = subroutineTable[name].Type
	if type_ != "" {
		return type_, nil
	}

	// If nothing was found, check class scope
	type_ = classTable[name].Type
	if type_ != "" {
		return type_, nil
	}

	return "", fmt.Errorf("identifier %v was not found in any scope", name)
}

// IndexOf returns the index assigned to the named identifier
func (s *SymbolTable) IndexOf(name string) (uint, error) {
	subroutineTable := *s.table[KIND_ARG]
	classTable := *s.table[KIND_STATIC]

	// Check subroutine scope first
	entry, ok := subroutineTable[name]
	if ok {
		return entry.Index, nil
	}

	// If nothing was found, check class scope
	entry, ok = classTable[name]
	if ok {
		return entry.Index, nil
	}

	return 0, fmt.Errorf("identifier %v was not found in any scope", name)
}
