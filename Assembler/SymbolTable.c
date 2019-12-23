#include "uthash.h"

typedef struct symbol_table_entry_s
{
  char *symbol;      // key
  uint16_t value;    // unsigned, 16-bit according to Hack spec
  UT_hash_handle hh; // makes this structure hashable
} symbol_table_entry_t;

symbol_table_entry_t *symbol_table = NULL; // important! initialize to NULL

void add_symbol(char *symbol, size_t symbol_length, uint16_t value)
{
  symbol_table_entry_t *new_entry;
  new_entry = malloc(sizeof(symbol_table_entry_t));

  // allocate memory for chars in symbol plus the null terminator
  // TODO: Does this need to be freed or will uthash figure it out (I doubt it)?
  new_entry->symbol = calloc(symbol_length + 1, sizeof(char));
  strncpy(new_entry->symbol, symbol, symbol_length);
  new_entry->value = value;
  HASH_ADD_STR(symbol_table, symbol, new_entry);
}