#include "uthash.h"
#include <stdbool.h>
#include <stdio.h>

// Structure for table, see https://troydhanson.github.io/uthash/userguide.html
typedef struct symbol_table_entry_s
{
  const char *symbol; // key
  uint16_t value;     // unsigned, 16-bit according to Hack spec
  UT_hash_handle hh;  // makes this structure hashable
} symbol_table_entry_t;

// Global table
symbol_table_entry_t *symbol_table = NULL; // important! initialize to NULL

// Add a new symbol to the table
void SymbolTable__addEntry(const char *symbol, size_t symbol_length, uint16_t value)
{
  symbol_table_entry_t *new_entry;
  new_entry = malloc(sizeof(symbol_table_entry_t));

  // allocate memory for chars in symbol plus the null terminator
  new_entry->symbol = calloc(symbol_length + 1, sizeof(char));
  strncpy((char *)new_entry->symbol, symbol, symbol_length);
  new_entry->value = value;
  HASH_ADD_KEYPTR(hh, symbol_table, new_entry->symbol, symbol_length, new_entry);
}

// Checks if table contains symbol
bool SymbolTable__contains(char *symbol)
{
  symbol_table_entry_t *s;
  HASH_FIND_STR(symbol_table, symbol, s);
  return s != NULL;
}

// Gets the address of a given symbol
uint16_t SymbolTable__getAddress(char *symbol)
{
  symbol_table_entry_t *s;
  HASH_FIND_STR(symbol_table, symbol, s);
  // WARN: will segfault if symbol DNE
  return s->value;
}

// Deletes the table and cleans up all the memory
void SymbolTable__destroy()
{
  symbol_table_entry_t *entry, *tmp;

  HASH_ITER(hh, symbol_table, entry, tmp)
  {
    HASH_DEL(symbol_table, entry); // delete; symbol_table advances to next
    free((char *)entry->symbol);
    free(entry);
  }
}

// TODO: Turn this into a test
int main()
{
  char *test1 = "test1";
  char *test2 = "test2";
  uint16_t val = 5;
  SymbolTable__addEntry(test1, 5, val);
  if (SymbolTable__contains("test1"))
    printf("Contains test1\n");
  if (!(SymbolTable__contains(test2)))
    printf("Does not contain test2\n");
  printf("SymbolTable has address (%d) for symbol %s\n", SymbolTable__getAddress(test1), test1);
  SymbolTable__destroy();
  return 0;
}