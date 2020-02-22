set -e

clang -o Assembler main.c Parser.c SymbolTable.c -I./
