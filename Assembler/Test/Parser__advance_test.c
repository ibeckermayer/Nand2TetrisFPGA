#include "Parser.h"
#include "Test.h"

void A_COMMAND_test() {
  printf("\033[0mRunning tests for A_COMMAND's\n");
  parser_t *parser = Parser__create("Test/A_COMMAND_test.asm");
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
      if (parser->current_line_type == SKIP) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 3) {
      if (parser->current_line_type == A_COMMAND &&
          strcmp("TEST", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == 0 &&
          parser->assembly_code_line_number == 3) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 4) {
      if (parser->current_line_type == A_COMMAND &&
          strcmp("TEST", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == 1 &&
          parser->assembly_code_line_number == 4) {
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
      if (parser->current_line_type == A_COMMAND &&
          strcmp("TEST1", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == 2 &&
          parser->assembly_code_line_number == 6) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 7) {
      if (parser->current_line_type == A_COMMAND &&
          strcmp("TEST1", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == 3 &&
          parser->assembly_code_line_number == 7) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 8) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 9) {
      if (parser->current_line_type == A_COMMAND &&
          strcmp("1", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == 4 &&
          parser->assembly_code_line_number == 9) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 10) {
      if (parser->current_line_type == A_COMMAND &&
          strcmp("1", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == 5 &&
          parser->assembly_code_line_number == 10) {
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
      if (parser->current_line_type == A_COMMAND &&
          strcmp("10", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == 6 &&
          parser->assembly_code_line_number == 12) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 13) {
      if (parser->current_line_type == A_COMMAND &&
          strcmp("10", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == 7 &&
          parser->assembly_code_line_number == 13) {
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
      if (parser->current_line_type == SYNTAX_ERROR) {
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
      if (parser->current_line_type == A_COMMAND &&
          strcmp("_.$:", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == 8 &&
          parser->assembly_code_line_number == 18) {
        PASS;
      } else {
        printf("%s\n", parser->current_command_buf);
        FAIL;
      }
    } else if (line == 19) {
      if (parser->current_line_type == A_COMMAND &&
          strcmp("_.$:", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == 9 &&
          parser->assembly_code_line_number == 19) {
        PASS;
      } else {
        printf("%s\n", parser->current_command_buf);
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
    } else if (line == 22) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    } else if (line == 23) {
      if (parser->current_line_type == SYNTAX_ERROR) {
        PASS;
      } else {
        FAIL;
      }
    }
  }
  // Delete as error in order to delete .Hack files
  Parser__destroy(parser, 1);
  printf("\n");
}

void L_COMMAND_test() {
  printf("\033[0mRunning tests for L_COMMAND's\n");
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
      if (parser->current_line_type == L_COMMAND &&
          strcmp("TEST", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == -1 &&
          parser->assembly_code_line_number == 2) {
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
      if (parser->current_line_type == L_COMMAND &&
          strcmp("TEST", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == -1 &&
          parser->assembly_code_line_number == 4) {
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
      if (parser->current_line_type == L_COMMAND &&
          strcmp("TEST1", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == -1 &&
          parser->assembly_code_line_number == 6) {
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
      if (parser->current_line_type == L_COMMAND &&
          strcmp("TEST1", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == -1 &&
          parser->assembly_code_line_number == 8) {
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
      if (parser->current_line_type == L_COMMAND &&
          strcmp("_.$:", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == -1 &&
          parser->assembly_code_line_number == 16) {
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
      if (parser->current_line_type == L_COMMAND &&
          strcmp("_.$:", parser->current_command_buf) == 0 &&
          parser->machine_code_line_number == -1 &&
          parser->assembly_code_line_number == 18) {
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
  // Delete as error in order to delete .Hack files
  Parser__destroy(parser, 1);
  printf("\n");
}

int main() {
  A_COMMAND_test();
  L_COMMAND_test();
}