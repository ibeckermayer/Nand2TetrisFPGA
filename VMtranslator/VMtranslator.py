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
VMFileName = str


class Parser:
    '''
    Parses *.vm files. 
    
    Supports full line comments, doesn't support inline comments i.e. `push argument 0 // this will break the parser`.

    Doesn't support error checking: we'll assume VMFileName contains only valid VM code.
    '''

    def __init__(self, filename: VMFileName):
        self.filename = filename
        self.lines: TextIO = open(self.filename, 'r')
        self.cur_line: str = ''  # The current line being read, as a string
        self.cur_line_split: List[str] = []  # cur_line tokenized into a list
        self.command_type: Optional[
            CommandType] = None  # The type of command cur_line is interpreted as

    def advance(self) -> str:
        '''
        Advances the parser on to the next line, and determines the command type of the line
        Returns the line, which will be empty string ('') at the end of the file (useful for looping)
        '''

        def set_command_type():
            '''
            Called after reading a new line in order to set the command type
            '''
            if (self.cur_line_split[0] == 'add' or
                    self.cur_line_split[0] == 'sub' or
                    self.cur_line_split[0] == 'neg' or
                    self.cur_line_split[0] == 'eq' or
                    self.cur_line_split[0] == 'gt' or
                    self.cur_line_split[0] == 'lt' or
                    self.cur_line_split[0] == 'and' or
                    self.cur_line_split[0] == 'or' or
                    self.cur_line_split[0] == 'not'):
                self.command_type = CT.ARITHMETIC
            elif ('push' == self.cur_line_split[0]):
                self.command_type = CT.PUSH
            elif ('pop' == self.cur_line_split[0]):
                self.command_type = CT.POP

        # Read the next line and strip leading/trailing spaces
        self.cur_line = self.lines.readline().strip(' ')

        # If EOF, return EOF
        if self.cur_line == '':
            return self.cur_line

        # If this is a blank or commented line, skip it and advance to the next line
        if self.cur_line[:2] == '//' or self.cur_line == '\n':
            self.advance()
        else:  # Otherwise set the command type
            self.cur_line = self.cur_line.strip('\n')
            # Split the line into tokens for future processing
            self.cur_line_split = self.cur_line.split()
            set_command_type()

        return self.cur_line


class CodeWriter:

    def __init__(self, output_filename: str):
        self.output_file = open(output_filename, 'w')
        # Initialize the stack pointer to 256 (see page 170 in the pdf of the book for the spec)
        self.SP = 256

    def write_arithmetic(self, parser: Parser):
        '''
        Say you have vm code like:
           push constant x
           push constant y
           add
        The stack should behave like so:
                ---------               ---------
                |  ...  |               |  ...  |
                ---------               ---------
                |   x   |      add      |  x+y  |
                ---------   ========>   ---------
                |   y   |               |       |  <-- SP
                ---------               ---------
        SP -->  |       |               |       |
                ---------               ---------
        '''
        if parser.cur_line_split[0] == "add":
            # // VM: add
            # // stack pointer at 258
            # @257 // Decrement the stack pointer (now pointing to y) and load it into the A reg
            # D=M // Load y into the D reg
            # @256 // Decrement the stack pointer (now pointing to x) and load it into the A reg
            # M=D+M // Replace x with y+x
            # // Increment stack pointer, now at 257
            self.output_file.write("// VM: add\n")
            self.output_file.write(f"// stack pointer at {self.SP}\n")
            self.SP -= 1
            self.output_file.write(
                f"@{self.SP} // Decrement the stack pointer (now pointing to y) and load it into the A reg\n"
            )
            self.output_file.write("D=M // Load y into the D reg\n")
            self.SP -= 1
            self.output_file.write(
                f"@{self.SP} // Decrement the stack pointer (now pointing to x) and load it into the A reg\n"
            )
            self.output_file.write("M=D+M // Replace x with y+x\n")
            self.SP += 1
            self.output_file.write(
                f"// Increment stack pointer, now at {self.SP}\n")
            self.output_file.write('\n')

        return

    def write_push(self, parser: Parser):

        if parser.cur_line_split[1] == "constant":
            '''
            // VM: push constant 2
            // stack pointer at 256
            @2 // Load 2 into the A reg
            D=A // Move 2 to the D reg
            @256 // Load the stack pointer 256 into the A reg
            M=D // RAM[256] = 2
            // Increment the stack pointer, now at 257
            '''
            val: str = parser.cur_line_split[2]
            self.output_file.write(
                f"// VM: {' '.join(parser.cur_line_split[:3])}\n")
            self.output_file.write(f"// stack pointer at {self.SP}\n")
            self.output_file.write(f"@{val} // Load {val} into the A reg\n")
            self.output_file.write(f"D=A // Move {val} to the D reg\n")
            self.output_file.write(
                f"@{self.SP} // Load the stack pointer {self.SP} into the A reg\n"
            )
            self.output_file.write(f"M=D // RAM[{self.SP}] = {val}\n")

            # Increment the stack pointer
            self.SP += 1
            self.output_file.write(
                f"// Increment the stack pointer, now at {self.SP}\n")
            self.output_file.write('\n')

    def write_pop(self, parser: Parser):
        return


class VMtranslator:
    '''
    Parses and translates *.vm specification compliant files into assembly code to be run on the Hack machine architecture
    '''

    def __init__(self, directory_or_filename: str):
        '''
        directory_or_filename: The directory or *.vm file to translate
        '''
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
                if directory_or_filename.split('.')[-1] == 'vm':
                    self.parsers.append(Parser(filename))

        # Create CodeWriter
        self.codewriter = CodeWriter(
            directory_or_filename.split('.')[0].strip('/') + '.asm')

    def run(self):
        '''
        Main function that calls each parser to run, passing each parsed token into codewriter which then 
        writes the corresponding assembly code to the output file.
        '''
        for parser in self.parsers:
            # TODO: is there something we need to do when we change parsers (aka change to a new .vm file)?
            while (parser.advance()):
                # print(parser.cur_line)
                # parser parses the input, codewriter translates tokens to assembly, this section
                # contains the logic for which codewrite function to pass the parser to
                if parser.command_type == CT.ARITHMETIC:
                    self.codewriter.write_arithmetic(parser)
                elif parser.command_type == CT.PUSH:
                    self.codewriter.write_push(parser)
                elif parser.command_type == CT.POP:
                    self.codewriter.write_pop(parser)


if __name__ == "__main__":
    import sys
    vmt = VMtranslator(sys.argv[1])
    vmt.run()