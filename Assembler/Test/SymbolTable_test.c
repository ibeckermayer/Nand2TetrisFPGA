#include "SymbolTable.h"
#include "Test.h"

void test_initialization(symbol_table_entry_t *symbol_table) {
  if (SymbolTable__contains(&symbol_table, "SP") &&
      SymbolTable__getValue(&symbol_table, "SP") == 0) {
    PASS;
  } else {
    FAIL;
  }
  if (SymbolTable__contains(&symbol_table, "LCL") &&
      SymbolTable__getValue(&symbol_table, "LCL") == 1) {
    PASS;
  } else {
    FAIL;
  }
  if (SymbolTable__contains(&symbol_table, "ARG") &&
      SymbolTable__getValue(&symbol_table, "ARG") == 2) {
    PASS;
  } else {
    FAIL;
  }
  if (SymbolTable__contains(&symbol_table, "THIS") &&
      SymbolTable__getValue(&symbol_table, "THIS") == 3) {
    PASS;
  } else {
    FAIL;
  }
  if (SymbolTable__contains(&symbol_table, "THAT") &&
      SymbolTable__getValue(&symbol_table, "THAT") == 4) {
    PASS;
  } else {
    FAIL;
  }
  if (SymbolTable__contains(&symbol_table, "SCREEN") &&
      SymbolTable__getValue(&symbol_table, "SCREEN") == 0x4000) {
    PASS;
  } else {
    FAIL;
  }
  if (SymbolTable__contains(&symbol_table, "KDB") &&
      SymbolTable__getValue(&symbol_table, "KDB") == 0x6000) {
    PASS;
  } else {
    FAIL;
  }

  // TODO: test R0-15 for the sake of it
}

int main() {
  symbol_table_entry_t *symbol_table = SymbolTable__create();

  test_initialization(symbol_table);
  // TODO: test attempt to overwrite a symbol
  // TODO: test adding arbitrary symbol

  SymbolTable__destroy(&symbol_table);
  printf("\n");
}