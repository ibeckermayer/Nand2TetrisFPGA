#ifndef _PARSER_H_
#define _PARSER_H_
#include "SymbolTable.h"
#include <stdio.h>

#define PARSER_BUF_SIZE 1024

// Enum of command types
typedef enum {
  INIT,
  SKIP,
  A_COMMAND,
  C_COMMAND,
  L_COMMAND,
  SYNTAX_ERROR,
  END_OF_FILE
} line_type;

// Enum for tracking first and second pass
typedef enum { FIRST_PASS, SECOND_PASS } pass_type;

// Structer to manage parser state
typedef struct parser_s {
  FILE *input;
  const char *input_filename;
  FILE *output;
  const char *output_filename;
  char current_line_buf[PARSER_BUF_SIZE];    // Holds the whole line
  char current_command_buf[PARSER_BUF_SIZE]; // Holds just the relvant parts of
                                             // the command.
  line_type current_line_type;
  pass_type current_pass_type;
  int16_t machine_code_line_number; // default -1, so first value becomes 0 on
                                    // Parser__advance
  int assembly_code_line_number;    // default 0, so first value becomes 1 on
                                    // Parser__advance
  symbol_table_entry_t *symbol_table;
} parser_t;

// Creates parser_t object for input_filename, opens input FILE for reading
// Builds output_filename and opens output FILE for writing
// Initializes parser_t.current_line_type to INIT
parser_t *Parser__create(const char *input_filename);

// Reads the next line into the parser->current_line_buf
// Processes the line to determine parser->current_line_type
void Parser__advance(parser_t *parser);

// Updates the symbol table if the current line calls for it
void Parser__update_symbol_table(parser_t *);

// Cleanup
void Parser__destroy(parser_t *parser, int is_error);
#endif