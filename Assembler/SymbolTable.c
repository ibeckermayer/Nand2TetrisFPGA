#include "SymbolTable.h"
#include "uthash.h"
#include <stdbool.h>
#include <stdio.h>

symbol_table_entry_t *SymbolTable__create(void) {
  symbol_table_entry_t *symbol_table = NULL; // important! initialize to NULL

  // Initialize with predefined symbols acc. to spec in section 6.2.3
  SymbolTable__addEntry(&symbol_table, "SP", 0);
  SymbolTable__addEntry(&symbol_table, "LCL", 1);
  SymbolTable__addEntry(&symbol_table, "ARG", 2);
  SymbolTable__addEntry(&symbol_table, "THIS", 3);
  SymbolTable__addEntry(&symbol_table, "THAT", 4);
  SymbolTable__addEntry(&symbol_table, "SCREEN", 0x4000);
  SymbolTable__addEntry(&symbol_table, "KDB", 0x6000);
  char r[4] = {};
  for (int i = 0; i < 16; i++) {
    sprintf(r, "R%d", i);
    SymbolTable__addEntry(&symbol_table, r, i);
  }
  return symbol_table;
}

void SymbolTable__addEntry(symbol_table_entry_t **symbol_table,
                           const char *symbol, uint16_t value) {

  if (SymbolTable__contains(symbol_table, (char *)symbol)) {
    return;
  }
  symbol_table_entry_t *new_entry;
  new_entry = malloc(sizeof(symbol_table_entry_t));

  size_t symbol_length = strlen(symbol);

  // allocate memory for chars in symbol plus the null terminator
  new_entry->symbol = calloc(symbol_length + 1, sizeof(char));
  strncpy((char *)new_entry->symbol, symbol, symbol_length);
  new_entry->value = value;
  HASH_ADD_KEYPTR(hh, *symbol_table, new_entry->symbol, symbol_length,
                  new_entry);
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
    free((char *)entry->symbol);
    free(entry);
  }
}

// // TODO: Turn this into a test
// int main() {
//   symbol_table_entry_t *symbol_table = SymbolTable__create();
//   char *test1 = "test1";
//   char *test2 = "test2";
//   uint16_t val = 5;
//   SymbolTable__addEntry(&symbol_table, test1, val);
//   if (SymbolTable__contains(&symbol_table, "test1"))
//     printf("Contains test1\n");
//   if (!(SymbolTable__contains(&symbol_table, test2)))
//     printf("Does not contain test2\n");
//   printf("SymbolTable has address (%d) for symbol %s\n",
//          SymbolTable__getAddress(&symbol_table, test1), test1);
//   SymbolTable__destroy(&symbol_table);
//   return 0;
// }