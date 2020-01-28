#include "Parser.h"

int main() {
  parser_t *parser = Parser__create("Test/A_COMMAND_test.asm");

  while (parser->current_line_type != END_OF_FILE) {
    Parser__advance(parser);
    if (parser->current_line_type == A_COMMAND) {
      printf("PASS!\n");
    } else {
      printf("FAIL!\n");
    }
  }
}