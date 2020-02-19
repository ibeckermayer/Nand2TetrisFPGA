#include "SymbolTable.h"
#include "Test.h"

void test_initialization(symbol_table_entry_t *symbol_table) {
  // Tests that default register values are registered in the table on create
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
  char r[4] = {};
  for (int i = 0; i < 16; i++) {
    sprintf(r, "R%d", i);
    if (SymbolTable__contains(&symbol_table, r) &&
        SymbolTable__getValue(&symbol_table, r) == i) {
      PASS;
    } else {
      FAIL;
    }
  }
}

void test_no_overwrite(symbol_table_entry_t *symbol_table) {
  // Tests that values in symbol table CAN NOT be overwritten
  // Trying to overwrite R15 with 0
  SymbolTable__addEntry(&symbol_table, "R15", 0);
  if (SymbolTable__getValue(&symbol_table, "R15") == 15) {
    PASS;
  } else {
    FAIL;
  }
}

void test_add_symbol(symbol_table_entry_t *symbol_table) {
  SymbolTable__addEntry(&symbol_table, "entry", 1000);

  if (SymbolTable__getValue(&symbol_table, "entry") == 1000) {
    PASS;
  } else {
    FAIL;
  }
}

int main() {
  symbol_table_entry_t *symbol_table = SymbolTable__create();

  test_initialization(symbol_table);
  test_no_overwrite(symbol_table);
  test_add_symbol(symbol_table);

  SymbolTable__destroy(&symbol_table);
  printf("\n");
}