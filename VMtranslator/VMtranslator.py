import os
from typing import List, TextIO, Optional
from enum import Enum


class CommandType(Enum):
    SKIP = 0
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
            if self.cur_line[:2] == '//' or self.cur_line == '\n':
                self.command_type = CT.SKIP
                return

            self.cur_line = self.cur_line.strip('\n')
            # Split the line into tokens for future processing
            self.cur_line_split = self.cur_line.split()

            if (self.cur_line_split[0] == 'add' or self.cur_line_split[0] == 'sub' or
                    self.cur_line_split[0] == 'neg' or self.cur_line_split[0] == 'eq' or
                    self.cur_line_split[0] == 'gt' or self.cur_line_split[0] == 'lt' or
                    self.cur_line_split[0] == 'and' or self.cur_line_split[0] == 'or' or
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

        set_command_type()

        return self.cur_line


class CodeWriter:
    '''
    RAM Addresses
    16-255: Static variables (1 segment per VM file)
    256-2047: Stack (1 local and 1 argument segment per VM function)
    2048-16483: Heap (used to store objects and arrays)
    16384-24575: Memory mapped I/O

    See "Figure 7.6 The memory segments seen by every VM function"
    '''

    def __init__(self, output_filename: str):
        self.output_file = open(output_filename, 'w')
        # Initialize the stack pointer to 256 (see page 170 in the pdf of the book for the spec)
        self.SP = 256
        # TODO: Initial vals of LCL, ARG are inferred from projects/07/MemoryAccess/BasicTest/BasicTest.vm
        # and projects/07/MemoryAccess/BasicTest/BasicTest.cmp. Probably there are better initial values
        self.LCL = 300  # allocated per VM function
        self.ARG = 400  # allocated per VM function
        self.STATIC = 16  # allocated per VM file
        self.TEMP = 5  # constant, should never change
        self.PTR = 3  # constant, should never change
        self.eq_num = 0  # Used as an identifier in `eq` operations
        self.gt_num = 0  # Used as an identifier in `gt` operations
        self.lt_num = 0  # Used as an identifier in `lt` operations
        # TODO: Initial vals for this and that are inferred from projects/07/MemoryAccess/BasicTest/BasicTest.vm
        # and projects/07/MemoryAccess/BasicTest/BasicTest.cmp. Probably there are better initial values
        self.output_file.write(f"// init\n")
        self.output_file.write(f"@3000\n")  # initial THIS value
        self.output_file.write(f"D=A\n")
        self.output_file.write(f"@{self.PTR}\n")
        self.output_file.write(f"M=D\n")
        self.output_file.write(f"@3010\n")  # initial THAT value
        self.output_file.write(f"D=A\n")
        self.output_file.write(f"@{self.PTR+1}\n")
        self.output_file.write(f"M=D\n")
        self.output_file.write(f'\n')

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
            x*  |   x   |      add      |  x+y  |
                ---------   ========>   ---------
            y*  |   y   |               |       |  <-- SP
                ---------               ---------
        SP -->  |       |               |       |
                ---------               ---------
        
        NOTE: Any time dest=M, the stack pointer should be incremented (self.SP += 1)
        '''
        if parser.cur_line_split[0] == "add":
            # // add
            # // SP--
            # @SP // A = *y
            # D=M // D = y
            # // SP--
            # @SP // A = *x
            # M=D+M // RAM[*x] = y + x
            # // SP++
            self.output_file.write(f"// add\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"D=M\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=D+M\n")
            self.SP += 1
        elif parser.cur_line_split[0] == "sub":
            # // sub
            # // SP--
            # @SP // A = *y
            # D=M // D = y
            # // SP--
            # @SP // A = *x
            # M=M-D // RAM[*x] = x - y
            # // SP++
            self.output_file.write(f"// sub\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"D=M\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=M-D\n")
            self.SP += 1
        elif parser.cur_line_split[0] == "neg":
            # // neg
            # // SP--
            # @SP // A = *y
            # M=-M // Negate y and replace it with it's negated value
            # // Increment stack pointer, now at 258
            self.output_file.write(f"// neg\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=-M\n")
            self.SP += 1
        elif parser.cur_line_split[0] == "eq":
            # // eq
            # // SP--
            # @SP // A = *y
            # D=M // D = y
            # // SP--
            # @SP // A = *x
            # D=M-D // D = x - y, 0 if eq is True
            # @eq0True // Load instruction for true case into the A register
            # D;JEQ // If D == 0, x==y. Jump over the false case to the true
            # // false case: if x != y
            # @SP // A = *x
            # M=0 // RAM[*x] = 0 (False)
            # @eq0TrueEnd // Load instruction to skip the true case
            # 0;JMP // jump over the true case
            # (eq0True) // true case: elif x == y
            # @SP // A = *x
            # M=-1 // RAM[*x] = -1 (True)
            # (eq0TrueEnd)
            # // SP++
            self.output_file.write(f"// eq\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"D=M\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"D=M-D\n")
            self.output_file.write(f"@eq{self.eq_num}True\n")
            self.output_file.write(f"D;JEQ\n")
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=0\n")
            self.output_file.write(f"@eq{self.eq_num}TrueEnd\n")
            self.output_file.write(f"0;JMP\n")
            self.output_file.write(f"(eq{self.eq_num}True)\n")
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=-1\n")
            self.output_file.write(f"(eq{self.eq_num}TrueEnd)\n")
            self.SP += 1
            self.eq_num += 1
        elif parser.cur_line_split[0] == "gt":
            # // gt
            # // SP--
            # @SP // A = *y
            # D=M // D = y
            # // SP--
            # @SP // A = *x
            # D=M-D // D = x - y, positive if gt is True
            # @gt0True // Load instruction for true case into the A register
            # D;JGT // If D == 0, x==y. Jump over the false case to the true
            # // false case: if x != y
            # @SP // A = *x
            # M=0 // RAM[*x] = 0 (False)
            # @gt0TrueEnd // Load instruction to skip the true case
            # 0;JMP // jump over the true case
            # (gt0True) // true case: elif x == y
            # @SP // A = *x
            # M=-1 // RAM[*x] = -1 (True)
            # (gt0TrueEnd)
            # // SP++
            self.output_file.write(f"// gt\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"D=M\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"D=M-D\n")
            self.output_file.write(f"@gt{self.gt_num}True\n")
            self.output_file.write(f"D;JGT\n")
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=0\n")
            self.output_file.write(f"@gt{self.gt_num}TrueEnd\n")
            self.output_file.write(f"0;JMP\n")
            self.output_file.write(f"(gt{self.gt_num}True)\n")
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=-1\n")
            self.output_file.write(f"(gt{self.gt_num}TrueEnd)\n")
            self.SP += 1
            self.gt_num += 1
        elif parser.cur_line_split[0] == "lt":
            # // lt
            # // SP--
            # @SP // A = *y
            # D=M // D = y
            # // SP--
            # @SP // A = *x
            # D=M-D // D = x - y, positive if lt is True
            # @lt0True // Load instruction for true case into the A register
            # D;JLT // If D == 0, x==y. Jump over the false case to the true
            # // false case: if x != y
            # @SP // A = *x
            # M=0 // RAM[*x] = 0 (False)
            # @lt0TrueEnd // Load instruction to skip the true case
            # 0;JMP // jump over the true case
            # (lt0True) // true case: elif x == y
            # @SP // A = *x
            # M=-1 // RAM[*x] = -1 (True)
            # (lt0TrueEnd)
            # // SP++
            self.output_file.write(f"// lt\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"D=M\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"D=M-D\n")
            self.output_file.write(f"@lt{self.lt_num}True\n")
            self.output_file.write(f"D;JLT\n")
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=0\n")
            self.output_file.write(f"@lt{self.lt_num}TrueEnd\n")
            self.output_file.write(f"0;JMP\n")
            self.output_file.write(f"(lt{self.lt_num}True)\n")
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=-1\n")
            self.output_file.write(f"(lt{self.lt_num}TrueEnd)\n")
            self.SP += 1
            self.lt_num += 1
        elif parser.cur_line_split[0] == "and":
            # // and
            # // SP--
            # @SP // A = *y
            # D=M // D = y
            # // SP--
            # @SP // A = *x
            # M=D&M // RAM[*x] = x - y
            # // SP++
            self.output_file.write(f"// and\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"D=M\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=D&M\n")
            self.SP += 1
        elif parser.cur_line_split[0] == "or":
            # // or
            # // SP--
            # @SP // A = *y
            # D=M // D = y
            # // SP--
            # @SP // A = *x
            # M=D|M // RAM[*x] = x - y
            # // SP++
            self.output_file.write(f"// or\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"D=M\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=D|M\n")
            self.SP += 1
        elif parser.cur_line_split[0] == "not":
            # // not
            # // SP--
            # @SP // A = *y
            # M=!M // Negate y and replace it with it's negated value
            # // Increment stack pointer, now at 258
            self.output_file.write(f"// neg\n")
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=!M\n")
            self.SP += 1

        self.output_file.write(f'\n')

    def write_push(self, parser: Parser):
        # Write a comment with the VM code for reference/debugging
        self.output_file.write(f"// {' '.join(parser.cur_line_split[:3])}\n")

        def __write_push_vtreg(base: int, offset: int):
            '''
            vtreg is short for "virtually tracked register"

            To save instruction (ROM) space, the local (LCL), argument (ARG), static (STATIC), temp (TEMP),
            and pointer (PTR), pointers are all tracked virtually by the CodeWriter.

            This function provides the pattern for pushing the value pointed to by base + offset onto the stack
            ex: `push static 3` --> __write_push_vtreg(self.STATIC, 3)
            '''
            self.output_file.write(f"@{base + offset}\n")
            self.output_file.write(f"D=M\n")
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=D\n")
            self.SP += 1

        def __write_push_thisthat(thisthat: str, offset: int):
            '''
            Writes a push instruction for this or that
            '''
            ptr: int = self.PTR if thisthat == 'this' else self.PTR + 1
            # Load this or that pointer into the CPU and add offset
            self.output_file.write(f"@{ptr}\n")
            if offset > 0:
                self.output_file.write(f"D=M\n")
                self.output_file.write(f"@{offset}\n")
                self.output_file.write(f"A=D+A\n")
            else:
                self.output_file.write(f"A=M\n")
            # Now grab the value being pointed to
            self.output_file.write(f"D=M\n")
            # And push onto stack
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=D\n")
            self.SP += 1

        if parser.cur_line_split[1] == "constant":
            '''
            // push constant 2
            // SP = 256
            @2 // Load 2 into the A reg
            D=A // Move 2 to the D reg
            @256 // Load the stack pointer 256 into the A reg
            M=D // RAM[256] = 2
            // Increment the stack pointer, now at 257
            '''
            val: str = parser.cur_line_split[2]

            self.output_file.write(f"@{val}\n")
            self.output_file.write(f"D=A\n")
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"M=D\n")
            self.SP += 1
        elif parser.cur_line_split[1] == "local":
            __write_push_vtreg(self.LCL, int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "argument":
            __write_push_vtreg(self.ARG, int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "static":
            __write_push_vtreg(self.STATIC, int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "temp":
            __write_push_vtreg(self.TEMP, int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "pointer":
            __write_push_vtreg(self.PTR, int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "this" or parser.cur_line_split[1] == "that":
            __write_push_thisthat(parser.cur_line_split[1], int(parser.cur_line_split[2]))
        else:
            raise RuntimeError(f"Unkown push command: {' '.join(parser.cur_line_split[:3])}")

        self.output_file.write(f'\n')

    def write_pop(self, parser: Parser):
        self.output_file.write(f"// {' '.join(parser.cur_line_split[:3])}\n")

        def __write_pop_vtreg(base: int, offset: int):
            '''
            vtreg is short for "virtually tracked register"

            This function provides the pattern for popping the value at the top of the stack off into the
            register pointed at by base + offset
            ex: `pop static 3` --> __write_pop_vtreg(self.STATIC, 3)
            '''
            # Pop value from the top of the stack into the D reg
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"D=M\n")

            # Copy D reg into base + offset
            self.output_file.write(f"@{base + offset}\n")
            self.output_file.write(f"M=D\n")

        def __write_pop_thisthat(thisthat: str, offset: int):
            '''
            Writes a pop instruction for this or that
            '''
            ptr: int = self.PTR if thisthat == 'this' else self.PTR + 1
            # Load this or that pointer into the CPU and add offset, then save it in R13
            self.output_file.write(f"@{ptr}\n")
            self.output_file.write(f"D=M\n")
            if offset > 0:
                self.output_file.write(f"@{offset}\n")
                self.output_file.write(f"D=D+A\n")
            self.output_file.write(f"@R13\n")
            self.output_file.write(f"M=D\n")

            # Pop the top value off the stack and save it in D
            self.SP -= 1
            self.output_file.write(f"@{self.SP}\n")
            self.output_file.write(f"D=M\n")

            # Now grab the this or that pointer from R13, and set the memory it points to
            # to the previously-top-of-the-stack value stored in D
            self.output_file.write(f"@R13\n")
            self.output_file.write(f"A=M\n")
            self.output_file.write(f"M=D\n")

        if parser.cur_line_split[1] == "local":
            __write_pop_vtreg(self.LCL, int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "argument":
            __write_pop_vtreg(self.ARG, int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "static":
            __write_pop_vtreg(self.STATIC, int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "temp":
            __write_pop_vtreg(self.TEMP, int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "pointer":
            __write_pop_vtreg(self.PTR, int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "this" or parser.cur_line_split[1] == "that":
            __write_pop_thisthat(parser.cur_line_split[1], int(parser.cur_line_split[2]))
        else:
            raise RuntimeError(f"Unkown pop command: {' '.join(parser.cur_line_split[:3])}")

        self.output_file.write(f'\n')
        pass


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
                raise RuntimeError("Files passed to VMTranslator must end with the .vm extension")
            self.parsers.append(Parser(directory_or_filename))  # only one parser
        else:
            # Else directory_or_filename is a directory, walk through each file and give it a parser
            for filename in os.listdir(directory_or_filename):
                if directory_or_filename.split('.')[-1] == 'vm':
                    self.parsers.append(Parser(filename))

        # Create CodeWriter
        self.codewriter = CodeWriter(directory_or_filename.split('.')[0].strip('/') + '.asm')

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
                if parser.command_type == CT.SKIP:
                    continue
                elif parser.command_type == CT.ARITHMETIC:
                    self.codewriter.write_arithmetic(parser)
                elif parser.command_type == CT.PUSH:
                    self.codewriter.write_push(parser)
                elif parser.command_type == CT.POP:
                    self.codewriter.write_pop(parser)

        # Don't forget to close the output file when you're done
        self.codewriter.output_file.close()


if __name__ == "__main__":
    import sys
    vmt = VMtranslator(sys.argv[1])
    vmt.run()