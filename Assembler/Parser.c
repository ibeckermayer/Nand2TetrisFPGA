#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define BUF_SIZE 1024

typedef enum {
  INIT,
  SKIP,
  A_COMMAND,
  C_COMMAND,
  L_COMMAND,
  SYNTAX_ERROR
} line_type;

typedef struct parser_s {
  FILE *input;
  const char *input_filename;
  FILE *output;
  const char *output_filename;
  char current_line_buf[BUF_SIZE];
  line_type current_line_type;
} parser_t;

// Takes X.asm and returns X.hack
const char *dot_hack_from_dot_asm(const char *dot_asm) {
  size_t asm_num_chars = strlen(dot_asm);

  // extra char for asm -> hack, plus null term
  char *dot_hack = (char *)calloc(asm_num_chars + 2, sizeof(char));
  // copy "X."
  strncpy(dot_hack, dot_asm, asm_num_chars - 3);
  // then make "X.hack"
  strncpy(dot_hack + (asm_num_chars - 3), "hack", 4);

  return dot_hack;
}

void remove_spaces(char *str_trimmed, char *str_untrimmed) {
  while (*str_untrimmed != '\0') {
    if (!(*str_untrimmed == ' ')) {
      *str_trimmed = *str_untrimmed;
      str_trimmed++;
    }
    str_untrimmed++;
  }
  *str_trimmed = '\0';
}

parser_t *Parser__create(const char *input_filename) {
  parser_t *parser = (parser_t *)malloc(sizeof(parser_t));
  // Register file names
  parser->input_filename = input_filename;
  parser->output_filename = dot_hack_from_dot_asm(input_filename);
  // Open input and output files
  parser->input = fopen(parser->input_filename, "r");
  parser->output = fopen(parser->output_filename, "w");
  parser->current_line_type = INIT;

  return parser;
}

void set_command_type(parser_t *parser) {
  char first_char = *(parser->current_line_buf);
  char second_char = *(parser->current_line_buf + 1);

  if (first_char == '\n') { // \n or // -> skip line
    parser->current_line_type = SKIP;
  } else if (first_char == '/') {
    if (second_char == '/') {
      parser->current_line_type = SKIP;
    } else {
      parser->current_line_type = SYNTAX_ERROR;
    }
  } else if (first_char == '@') {
    // TODO: check that subsequent chars are valid
    parser->current_line_type = A_COMMAND;
  } else if (first_char == '(') {
    // TODO: check that subsequent chars are valid
    parser->current_line_type = L_COMMAND;
  } else if (/* condition */) {
    // D, A, M, 0, 1, -, !
  }
}

void Parser__advance(parser_t *parser) {
  // Reads the next line into the parser->current_line_buf
  fgets(parser->current_line_buf, BUF_SIZE, parser->input);
  // parser->current_line_buf for both args will remove spaces in place
  remove_spaces(parser->current_line_buf, parser->current_line_buf);
  // \n or // -> skip line
  // @ -> A_COMMAND
  // ( -> L_COMMAND
  // D, A, M, 0, 1, -, !
  set_command_type(parser);
}

void Parser__destroy(parser_t *parser) {
  // no need to free input_filename which comes from argv
  free((void *)parser->output_filename);
  fclose(parser->input);
  fclose(parser->output);
  free(parser);
}

int main(int argc, char *argv[]) {
  parser_t *parser = Parser__create(argv[1]);
  // printf("input_filename = %s\n", parser->input_filename);
  // printf("output_filename = %s\n", parser->output_filename);
  // Parser__advance(parser);
  // printf("current_line = %s", parser->current_line_buf);
  Parser__advance(parser);
  // printf("current_line = %s", parser->current_line_buf);
  Parser__advance(parser);
  // printf("current_line = %s", (parser->current_line_buf + 1));
  Parser__destroy(parser);
  return 0;
}