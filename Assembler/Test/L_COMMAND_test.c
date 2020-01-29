#include "Parser.h"
#include "Test.h"

int main() {
  parser_t *parser = Parser__create("Test/L_COMMAND_test.asm");
  int line = 0;

  while (parser->current_line_type != END_OF_FILE) {
    Parser__advance(parser);
    ++line;
    if (line == 1) {
      if (parser->current_line_type == SKIP) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 2) {
      if (parser->current_line_type == L_COMMAND) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 3) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 4) {
      if (parser->current_line_type == L_COMMAND) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 5) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 6) {
      if (parser->current_line_type == L_COMMAND) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 7) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 8) {
      if (parser->current_line_type == L_COMMAND) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 9) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 10) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 11) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 12) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 13) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 14) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 15) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 16) {
      if (parser->current_line_type == L_COMMAND) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 17) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 18) {
      if (parser->current_line_type == L_COMMAND) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 19) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 20) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 21) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    }
  }
}