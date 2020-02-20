set -e

clang -o Test/Parser__advance_test Parser.c SymbolTable.c Test/Parser__advance_test.c -I./ -D TEST
clang -o Test/SymbolTable_test SymbolTable.c Test/SymbolTable_test.c -I./ -D TEST