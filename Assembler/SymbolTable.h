#ifndef _SymbolTable_H_
#define _SymbolTable_H_
#include "uthash.h"
#include <stdbool.h>
#include <stdio.h>

// Structure for table, see https://troydhanson.github.io/uthash/userguide.html
typedef struct symbol_table_entry_s {
  const char *symbol; // key
  uint16_t value;     // unsigned, 16-bit according to Hack spec
  UT_hash_handle hh;  // makes this structure hashable
} symbol_table_entry_t;

// Creates pointer to symbol_table_entry_t equal to NULL, which can be thought
// of as the "root" object for the table. All other functions will take a
// POINTER to this pointer to manipulate the table. From the documentation: "The
// reason it’s necessary to deal with a pointer to the hash pointer is simple:
// the hash macros modify it (in other words, they modify the pointer itself not
// just what it points to)."
symbol_table_entry_t *SymbolTable__create(void);

// Add a new symbol to the table
void SymbolTable__addEntry(symbol_table_entry_t **symbol_table,
                           const char *symbol, uint16_t value);

// Checks if table contains symbol
bool SymbolTable__contains(symbol_table_entry_t **symbol_table, char *symbol);

// Gets the address of a given symbol
uint16_t SymbolTable__getValue(symbol_table_entry_t **symbol_table,
                               char *symbol);

// Deletes the table and cleans up all the memory
void SymbolTable__destroy(symbol_table_entry_t **symbol_table);

#endif