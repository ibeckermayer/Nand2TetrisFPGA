#include "SymbolTable.h"
#include "uthash.h"
#include <stdbool.h>
#include <stdio.h>

symbol_table_entry_t *SymbolTable__create(void) {
  symbol_table_entry_t *symbol_table = NULL; // important! initialize to NULL

  // Initialize with predefined symbols acc. to spec in section 6.2.3
  SymbolTable__addEntry(&symbol_table, "SP", 0, 0);
  SymbolTable__addEntry(&symbol_table, "LCL", 1, 0);
  SymbolTable__addEntry(&symbol_table, "ARG", 2, 0);
  SymbolTable__addEntry(&symbol_table, "THIS", 3, 0);
  SymbolTable__addEntry(&symbol_table, "THAT", 4, 0);
  SymbolTable__addEntry(&symbol_table, "SCREEN", 0x4000, 0);
  SymbolTable__addEntry(&symbol_table, "KDB", 0x6000, 0);
  char r[4] = {};
  for (int i = 0; i < 16; i++) {
    sprintf(r, "R%d", i);
    SymbolTable__addEntry(&symbol_table, r, i, 0);
  }
  return symbol_table;
}

void SymbolTable__addEntry(symbol_table_entry_t **symbol_table,
                           const char *new_symbol, uint16_t value,
                           int replace) {

  if (SymbolTable__contains(symbol_table, (char *)new_symbol) && !(replace)) {
    return;
  }
  symbol_table_entry_t *new_entry = malloc(sizeof(symbol_table_entry_t));

  // allocate memory for chars in symbol plus the null terminator
  strcpy(new_entry->symbol, new_symbol);
  new_entry->value = value;
  if (!replace) {
    HASH_ADD_STR(*symbol_table, symbol, new_entry);
  } else {
    symbol_table_entry_t *replaced_item_ptr;
    HASH_REPLACE_STR(*symbol_table, symbol, new_entry, replaced_item_ptr);
    if (replaced_item_ptr != NULL) {
      free(replaced_item_ptr);
    }
  }
}

bool SymbolTable__contains(symbol_table_entry_t **symbol_table, char *symbol) {
  symbol_table_entry_t *s;
  HASH_FIND_STR(*symbol_table, symbol, s);
  return s != NULL;
}

uint16_t SymbolTable__getValue(symbol_table_entry_t **symbol_table,
                               char *symbol) {
  symbol_table_entry_t *s;
  HASH_FIND_STR(*symbol_table, symbol, s);
  // WARN: will segfault if symbol DNE
  return s->value;
}

void SymbolTable__destroy(symbol_table_entry_t **symbol_table) {
  symbol_table_entry_t *entry, *tmp;

  HASH_ITER(hh, *symbol_table, entry, tmp) {
    HASH_DEL(*symbol_table, entry); // delete; symbol_table advances to next
    free(entry);
  }
}
