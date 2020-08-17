set -e

COMPILER=""
if command -v clang &> /dev/null; then
    COMPILER=clang
elif command -v gcc &> /dev/null; then 
    COMPILER=gcc
else 
    printf "This program requires clang or gcc to be installed\n" 1>&2
    exit 1
fi 

"$COMPILER" -o Test/Parser__advance_test Parser.c SymbolTable.c Test/Parser__advance_test.c -I./ -D TEST
"$COMPILER" -o Test/SymbolTable_test SymbolTable.c Test/SymbolTable_test.c -I./ -D TEST