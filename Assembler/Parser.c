#include "Parser.h"
#include "SymbolTable.h"
#include <stdlib.h>
#include <string.h>

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
  parser->current_pass_type = FIRST_PASS;
  parser->machine_code_line_number = -1;
  parser->assembly_code_line_number = 0;
  parser->next_A_COMMAND_symbol_RAM_addr = 16;
  return parser;
}

int is_number(char c) { return (c >= '0' && c <= '9'); }

int is_valid_constant_non_number(char c) {
  // Constants must be non-negative and are written in decimal notation. A user-
  // defined symbol can be any sequence of letters, digits, underscore (_), dot
  // (.), dollar sign ($), and colon (:) that does not begin with a digit.
  return ((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' ||
          c == '.' || c == '$' || c == ':');
}

int is_line_end(char c) {
  // checks if the character means the relevant section of the line is over
  return (c == '\n' || c == '\0' || c == '/');
}

// Returns 1 if valid or 0 if invalid
// Caller should omit '@' symbol
int is_valid_A_COMMAND(char *cmd) {
  int ptr_offset = 1;
  // if first char is a number, check that all the rest are numbers (barring
  // comments)
  if (is_number(*cmd)) {
    while (!(is_line_end(*(cmd + ptr_offset)))) {
      if (is_number(*(cmd + ptr_offset))) {
        ptr_offset++;
      } else {
        return 0;
      }
    }
  } else if (is_valid_constant_non_number(*cmd)) {
    //  else if first char is_valid_constant_non_number, check that
    // all the rest are either that or numerical
    while (!(is_line_end(*(cmd + ptr_offset)))) {
      if (is_number(*(cmd + ptr_offset)) ||
          is_valid_constant_non_number(*(cmd + ptr_offset))) {
        ptr_offset++;
      } else {
        return 0;
      }
    }
  } else {
    return 0;
  }
  // If line end was a '/', check that there's another '/' to confirm its a
  // valid comment and not a syntax error
  if (*(cmd + ptr_offset) == '/') {
    if (*(cmd + ptr_offset + 1) == '/') {
      return 1;
    } else {
      return 0;
    }
  }
  return 1;
}

// Returns 1 if valid or 0 if invalid
// Caller should omit leading '(' symbol
int is_valid_L_COMMAND(char *cmd) {
  int ptr_offset = 1;

  // Must begin with valid symbol
  if (!(is_valid_constant_non_number(*cmd))) {
    return 0;
  }

  // Check up to what should be the closing ')'
  while (
      (!(is_line_end(*(cmd + ptr_offset))) && (*(cmd + ptr_offset) != ')'))) {
    if (is_valid_constant_non_number(*(cmd + ptr_offset)) ||
        is_number(*(cmd + ptr_offset))) {
      ptr_offset++;
    } else {
      return 0;
    }
  }

  // Check for closing ')'
  if (*(cmd + ptr_offset) != ')') {
    return 0;
  }

  // Now go on to the end and check for comment syntax
  while (!(is_line_end(*(cmd + ptr_offset)))) {
    ptr_offset++;
  }

  // If line end was a '/', check that there's another '/' to confirm its a
  // valid comment and not a syntax error
  if (*(cmd + ptr_offset) == '/') {
    if (*(cmd + ptr_offset + 1) == '/') {
      return 1;
    } else {
      return 0;
    }
  }
  return 1;
}

// Finds the pointer to the end of line according to is_line_end()
char *strchr_line_end(char *cmd) {
  while (!(is_line_end(*cmd))) {
    cmd++;
  }
  return cmd;
}

void set_command_type(parser_t *parser) {
  // \n or // -> skip line
  // @ -> A_COMMAND
  // ( -> L_COMMAND
  // D, A, M, 0, 1, -, ! -> C_COMMAND

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
    if (is_valid_A_COMMAND(parser->current_line_buf + 1)) {
      parser->current_line_type = A_COMMAND;
    } else {
      parser->current_line_type = SYNTAX_ERROR;
    }
  } else if (first_char == '(') {
    if (is_valid_L_COMMAND(parser->current_line_buf + 1)) {
      parser->current_line_type = L_COMMAND;
    } else {
      parser->current_line_type = SYNTAX_ERROR;
    }
  } else if (first_char == 'D' || first_char == 'A' || first_char == 'M' ||
             first_char == '0' || first_char == '1' || first_char == '-' ||
             first_char == '!') {
    //  Don't need to check syntax here, it can more easily be checked during
    //  translation to binary
    parser->current_line_type = C_COMMAND;
  } else {
    parser->current_line_type = SYNTAX_ERROR;
  }
}

void set_machine_code_line_number(parser_t *parser) {
  if (parser->current_line_type == A_COMMAND ||
      parser->current_line_type == C_COMMAND) {
    parser->machine_code_line_number++;
  }
}

void syntax_error(parser_t *parser) {
#ifndef TEST
  if (parser->current_line_type == SYNTAX_ERROR) {
    fprintf(stderr, "Syntax Error on line: %d",
            parser->assembly_code_line_number);
    Parser__delete(parser, 1);
    exit(1);
  }
#endif
}

void finalize_extract(parser_t *parser, char *first_char, char *last_char) {
  size_t len = last_char - first_char;
  strncpy(parser->current_command_buf, first_char, len);
  *(parser->current_command_buf + len) = '\0';
}

void extract_A_COMMAND(parser_t *parser) {
  char *command_end = strchr_line_end(parser->current_line_buf);
  finalize_extract(parser, parser->current_line_buf + 1, command_end);
}

void extract_L_COMMAND(parser_t *parser) {
  char *command_end = strchr(parser->current_line_buf, ')');
  finalize_extract(parser, parser->current_line_buf + 1, command_end);
}

void extract_C_COMMAND(parser_t *parser) {
  char *command_end = strchr_line_end(parser->current_line_buf);
  finalize_extract(parser, parser->current_line_buf, command_end);
}

void extract_command(parser_t *parser) {
  if (parser->current_line_type == SKIP ||
      parser->current_line_type == SYNTAX_ERROR) {
    return;
  }

  if (parser->current_line_type == A_COMMAND) {
    extract_A_COMMAND(parser);
  } else if (parser->current_line_type == L_COMMAND) {
    extract_L_COMMAND(parser);
  } else if (parser->current_line_type == C_COMMAND) {
    extract_C_COMMAND(parser);
  } else {
    printf("Unexpected error in extract_command\n");
    syntax_error(parser);
  }
}

// Fills parser->current_command_buf with the command
void Parser__advance(parser_t *parser) {
  // Reads the next line into the parser->current_line_buf
  if (fgets(parser->current_line_buf, PARSER_BUF_SIZE, parser->input) == NULL) {
    parser->current_line_type = END_OF_FILE;
    return;
  }

  // Assembly code line number always advances, because we're reading a new line
  parser->assembly_code_line_number++;

  // parser->current_line_buf for both args will remove spaces in place
  remove_spaces(parser->current_line_buf, parser->current_line_buf);
  set_command_type(parser);             // sets parser->current_line_type
  extract_command(parser);              // sets parser->current_command_buf
  set_machine_code_line_number(parser); // sets parser->machine_code_line_number
}

// Adds current_command_buf to symbol table. Should only be called when
// (parser->current_line_type == L_COMMAND || parser->current_line_type ==
// A_COMMAND) && is_valid_constant_non_number(*(parser->current_command_buf)) &&
// parser->current_pass_type == FIRST_PASS
void Parser__update_symbol_table(parser_t *parser) {
  if (parser->current_line_type == L_COMMAND) {
    SymbolTable__addEntry(&(parser->symbol_table), parser->current_command_buf,
                          parser->machine_code_line_number + 1);
  } else if (parser->current_line_type == A_COMMAND) {
    SymbolTable__addEntry(&(parser->symbol_table), parser->current_command_buf,
                          parser->next_A_COMMAND_symbol_RAM_addr++);
  } else {
    printf("Unexpected symbole table update\n");
    syntax_error(parser);
  }
}

// Function that runs through the full process of assembling to machine code
void Parser__assemble(const char *input_filename) {
  parser_t *parser = Parser__create(input_filename);
  while (parser->current_line_type != END_OF_FILE) {
    Parser__advance(parser);

    if (parser->current_line_type == SKIP) {
      continue;
    }

    if (parser->current_line_type == SYNTAX_ERROR) {
      syntax_error(parser);
    }

    if ((parser->current_line_type == L_COMMAND ||
         parser->current_line_type == A_COMMAND) &&
        is_valid_constant_non_number(*(parser->current_command_buf))) {
      if (parser->current_pass_type == FIRST_PASS) {
        Parser__update_symbol_table(parser);
      }
    }

    // Extract relevant bits into current_command_buf
    // if is L_COMMAND || A_COMMAND and first letter is valid char and is first
    // pass, update symbol table
  }
}

void Parser__destroy(parser_t *parser, int is_error) {
  // no need to free input_filename which comes from argv
  free((void *)parser->output_filename);
  fclose(parser->input);
  fclose(parser->output);
  if (is_error) {
    remove(parser->output_filename);
  }
  free(parser);
}