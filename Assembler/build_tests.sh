# Build
clang -o Test/A_COMMAND_test Parser.c Test/A_COMMAND_test.c -I./ -D TEST
clang -o Test/L_COMMAND_test Parser.c Test/L_COMMAND_test.c -I./ -D TEST
clang -o Test/SymbolTable_test SymbolTable.c Test/SymbolTable_test.c -I./ -D TEST