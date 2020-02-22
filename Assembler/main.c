#include "Parser.h"

int main(int argc, char *argv[]) {
  if (argc != 2) {
    printf("Error: Assembler must be called like `Assembler prog.asm`\n");
    return 1;
  }

  Parser__run(argv[1]);
  return 0;
}