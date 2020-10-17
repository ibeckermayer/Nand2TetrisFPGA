import os
from typing import List, TextIO, Optional
from enum import Enum


class CommandType(Enum):
    ARITHMETIC = 1
    PUSH = 2
    POP = 3
    # C_LABEL = 4
    # C_GOTO = 5
    # C_IF = 6
    # C_FUNCTION = 7
    # C_RETURN = 8
    # C_CALL = 9


# Give CommandType an alias to make it a bit easier on the eye
CT = CommandType


class Parser:
    '''
    Parses *.vm files. Supports full line comments, doesn't support inline comments i.e. `push argument 0 // this will break the parser`
    '''

    def __init__(self, filename: str):
        self.filename = filename
        self.lines: TextIO = open(self.filename, 'r')
        self.cur_line: str = ''
        self.command_type: Optional[CommandType]

    def advance(self) -> str:
        '''
        Advances the parser on to the next line, and determines the command type of the line
        Returns the line, which will be '' at the end of the file
        '''

        def set_command_type():
            '''
            Called after reading a new line in order to set the command type
            '''
            if (self.cur_line == 'add' or self.cur_line == 'sub' or
                    self.cur_line == 'neg' or self.cur_line == 'eq' or
                    self.cur_line == 'gt' or self.cur_line == 'lt' or
                    self.cur_line == 'and' or self.cur_line == 'or' or
                    self.cur_line == 'not'):
                self.command_type = CT.ARITHMETIC
            elif ('push' in self.cur_line):
                self.command_type = CT.PUSH
            elif ('pop' in self.cur_line):
                self.command_type = CT.POP

        # Read the next line and strip leading/trailing spaces
        self.cur_line = self.lines.readline().strip(' ')
        # If this is a blank or commented line, skip it and advance to the next
        if self.cur_line[:2] == '//' or self.cur_line == '\n':
            self.advance()
        else:
            self.cur_line = self.cur_line.strip('\n')
        set_command_type()
        return self.cur_line


class CodeWriter:

    def __init__(self, output_filename: str):
        self.file = open(output_filename, 'w')


class VMtranslator:

    def __init__(self, directory_or_filename: str):
        self.parsers: List[Parser] = []

        # If directory_or_filename is a filename, check that its a .vm file
        if '.' in directory_or_filename:
            if directory_or_filename.split('.')[-1] != 'vm':
                raise RuntimeError(
                    "Files passed to VMTranslator must end with the .vm extension"
                )
            self.parsers.append(
                Parser(directory_or_filename))  # only one parser
        else:
            # Else directory_or_filename is a directory, walk through each file and give it a parser
            for filename in os.listdir(directory_or_filename):
                self.parsers.append(Parser(filename))

        # Create CodeWriter with first parser's filename
        self.codewriter = CodeWriter(
            self.__dot_vm_to_dot_hack(directory_or_filename.split('.')[0].strip('/'))
