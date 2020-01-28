#ifndef _PARSER_H_
#define _PARSER_H_
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

// Structer to manage parser state
typedef struct parser_s {
  FILE *input;
  const char *input_filename;
  FILE *output;
  const char *output_filename;
  char current_line_buf[PARSER_BUF_SIZE];
  line_type current_line_type;
} parser_t;

// Creates parser_t object for input_filename, opens input FILE for reading
// Builds output_filename and opens output FILE for writing
// Initializes parser_t.current_line_type to INIT
parser_t *Parser__create(const char *input_filename);

// Reads the next line into the parser->current_line_buf
// Processes the line to determine parser->current_line_type
void Parser__advance(parser_t *parser);

// Cleanup
void Parser__destroy(parser_t *parser);
#endif