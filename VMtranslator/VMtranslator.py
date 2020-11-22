import os
from typing import List, TextIO, Optional
from enum import Enum


class CommandType(Enum):
    SKIP = 0
    ARITHMETIC = 1
    PUSH = 2
    POP = 3
    LABEL = 4
    GOTO = 5
    IF_GOTO = 6
    FUNCTION = 7
    RETURN = 8
    CALL = 9


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
        self.cur_line_number = 0
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
            elif (self.cur_line_split[0] == 'push'):
                self.command_type = CT.PUSH
            elif (self.cur_line_split[0] == 'pop'):
                self.command_type = CT.POP
            elif (self.cur_line_split[0] == 'label'):
                self.command_type = CT.LABEL
            elif (self.cur_line_split[0] == 'if-goto'):
                self.command_type = CT.IF_GOTO
            elif (self.cur_line_split[0] == 'goto'):
                self.command_type = CT.GOTO
            elif (self.cur_line_split[0] == 'function'):
                self.command_type = CT.FUNCTION
            elif (self.cur_line_split[0] == 'return'):
                self.command_type = CT.RETURN
            elif (self.cur_line_split[0] == 'call'):
                self.command_type = CT.CALL

        # Read the next line and strip leading/trailing spaces
        self.cur_line = self.lines.readline().strip(' ')
        self.cur_line_number += 1

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

    # constant, should never change
    TEMP = 5

    def __init__(self, output_filename: str):
        self.output_file = open(output_filename, 'w')
        self.cur_parser_filename = ''  # Gets updated per parser, used for naming static variables
        self.cur_func_name = ''  # Gets updated per function declaration, used for naming labels within functions
        self.eq_num = 0  # Used as an identifier in `eq` operations
        self.gt_num = 0  # Used as an identifier in `gt` operations
        self.lt_num = 0  # Used as an identifier in `lt` operations

        # Initialization assembly code
        self.output_file.write(f"// init\n")
        # Initialize the stack pointer to address 256
        self.set_reg("SP", 256)
        # Call Sys.init
        self.output_file.write(f"@Sys.init\n")
        self.output_file.write(f"0;JMP\n")
        self.output_file.write(f"\n")

    def static_symbol(self, suffix: str) -> str:
        '''
        From section 7.3 of the book:
        "According to the Hack machine language specification, when a new symbol is encountered for the first time in an assembly 
        program, the assembler allocates a new RAM address to it, starting at address 16. This convention can be exploited to represent
        each static variable number j in a VM file f as the assembly language symbol f.j. For example, suppose that the file 
        Xxx.vm contains the command push static 3. This command can be translated to the Hack assembly commands @Xxx.3 and D=M, followed 
        by additional assembly code that pushes D’s value to the stack. This implementation of the static segment is somewhat tricky, but it works."
        '''
        return self.cur_parser_filename.split('/')[-1].split('.')[0] + f".{suffix}"

    def prefix_w_cur_func_name(self, label: str) -> str:
        '''
        Prefixes label with f"{self.cur_func_name}$". From the book (Figure 8.6):
        "Each `label b` command in a VM function `f` should generate a globally uniques symbol `f$b`..."

        If self.cur_func_name isn't set (empty string) then just return label (this is kinda jank but expedient for testing)
        '''
        if self.cur_func_name == '':
            return label
        return f"{self.cur_func_name}${label}"

    def create_ret_addr(self, cur_func_name: str, cur_line_number: int) -> str:
        '''
        Creates a unique return-address symbol for implementing the `call f n` stack machine function call operation.
        Ensures that the symbol is unique by using the function name and the (letter encoded) line number of the calling function.
        '''

        # 'ra' short for "return address"
        return f"ra_{cur_func_name}_{cur_line_number}"

    def set_reg(self, symbol: str, value: int):
        '''
        Sets the register symbol to value
        i.e. if you want to set the stack pointer to 256, call `set_reg("SP", 256)
        '''
        self.output_file.write(f"@{value}\n")
        self.output_file.write(f"D=A\n")
        self.output_file.write(f"@{symbol}\n")
        self.output_file.write(f"M=D\n")

    def SP_pp(self, load_SP_into_A):
        '''
        SP++
        If load_SP_into_A, loads the value of SP into the A register upon completion
        '''
        self.output_file.write(f"@SP\n")
        self.output_file.write(f"M=M+1\n")
        if load_SP_into_A:
            self.output_file.write(f"A=M\n")

    def SP_mm(self, load_SP_into_A):
        '''
        SP--
        If load_SP_into_A, loads the value of SP into the A register upon completion
        '''
        self.output_file.write(f"@SP\n")
        self.output_file.write(f"M=M-1\n")
        if load_SP_into_A:
            self.output_file.write(f"A=M\n")

    def push_value(self, symbol: str, offset: int):
        '''
        Pushes the value pointed to by symbol onto the stack, adjusting for offset. 
        
        This is the ordinary functioning of the stack machine,
        i.e. if we encounter `push local 2` we call `self.push_value("LCL", 2)`
        
        In this example, we will find the base pointer stored in the LCL register, add 2 to it to create our final pointer, 
        then push the value pointed to by that pointer onto the stack
        '''
        # Build the pointer in the A register
        self.output_file.write(f"@{symbol}\n")
        if offset > 0:
            self.output_file.write(f"D=M\n")
            self.output_file.write(f"@{offset}\n")
            self.output_file.write(f"A=D+A\n")
        else:
            self.output_file.write(f"A=M\n")

        # Grab the value being pointed to and store it in the D register
        self.output_file.write(f"D=M\n")

        # Push that onto the stack
        self.load_SP_into_A()
        self.output_file.write(f"M=D\n")

        # Increment the stack pointer
        self.SP_pp(load_SP_into_A=False)

    def push_pointer(self, symbol: str):
        '''
        Pushes the pointer referenced by symbol onto the stack. For example, `push pointer 0` should result in a call
        to self.push_pointer("THIS").

        In this example, we will find the pointer stored in the THIS register and push it onto the stack
        '''
        # Grab the pointer referenced by symbol and store it in the D register
        self.output_file.write(f"@{symbol}\n")
        self.output_file.write(f"D=M\n")

        # Push it onto the stack
        self.load_SP_into_A()
        self.output_file.write(f"M=D\n")

        # Increment the stack pointer
        self.SP_pp(load_SP_into_A=False)

    def pop_value(self, symbol: str, offset: int):
        '''
        Pops the value on the top of the stack into the register stored in symbol, adjusting for offset.

        This is the ordinary functioning of the stack machine, i.e. if we encounter `pop this 2` we can call
        `this.pop_value("THIS", 2)`

        In this example, this will find the value pointed to by the THIS register, add 2 to it to create our pointer, and then pop the
        value off the top of the stack and store it in the pointer
        '''
        # Load pointer stored at address symbol, add offset, and save it in R13
        self.output_file.write(f"@{symbol}\n")
        self.output_file.write(f"D=M\n")
        if offset > 0:
            self.output_file.write(f"@{offset}\n")
            self.output_file.write(f"D=D+A\n")
        self.output_file.write(f"@R13\n")
        self.output_file.write(f"M=D\n")

        # Pop the top value off the stack and save it in D
        self.SP_mm(load_SP_into_A=True)
        self.output_file.write(f"D=M\n")

        # Now grab the pointer from R13, and set the memory it points to
        # to the previously-top-of-the-stack value stored in D
        self.output_file.write(f"@R13\n")
        self.output_file.write(f"A=M\n")
        self.output_file.write(f"M=D\n")

    def pop_pointer(self, symbol: str):
        '''
        Pops the value on the top of the stack into the symbol register itself.

        For example, if we encounter `pop pointer 0`, we would call `self.pop_pointer("THIS")
        '''
        # Load the value at the top of the stack into D
        self.SP_mm(load_SP_into_A=True)
        self.output_file.write(f"D=M\n")

        # And move it into the register `symbol`
        self.output_file.write(f"@{symbol}\n")
        self.output_file.write(f"M=D\n")

    def load_SP_into_A(self):
        '''
        Loads the value of SP into the A register
        '''
        self.output_file.write(f"@SP\n")
        self.output_file.write(f"A=M\n")

    def load_SP_into_D(self):
        '''
        Loads the value of SP into the D register
        '''
        self.output_file.write(f"@SP\n")
        self.output_file.write(f"D=M\n")

    def goto_label(self, label: str):
        '''
        Loads address `label` into the A register and then jumps to that instruction
        '''
        self.output_file.write(f"@{label}\n")
        self.output_file.write(f"0;JMP\n")

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
            # load_SP_into_A() // A = *y
            # D=M // D = y
            # // SP--
            # load_SP_into_A() // A = *x
            # M=D+M // RAM[*x] = y + x
            # // SP++
            self.output_file.write(f"// add\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"D=M\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"M=D+M\n")
            self.SP_pp(load_SP_into_A=False)
        elif parser.cur_line_split[0] == "sub":
            # // sub
            # // SP--
            # load_SP_into_A() // A = *y
            # D=M // D = y
            # // SP--
            # load_SP_into_A() // A = *x
            # M=M-D // RAM[*x] = x - y
            # // SP++
            self.output_file.write(f"// sub\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"D=M\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"M=M-D\n")
            self.SP_pp(load_SP_into_A=False)
        elif parser.cur_line_split[0] == "neg":
            # // neg
            # // SP--
            # load_SP_into_A() // A = *y
            # M=-M // Negate y and replace it with it's negated value
            # // Increment stack pointer, now at 258
            self.output_file.write(f"// neg\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"M=-M\n")
            self.SP_pp(load_SP_into_A=False)
        elif parser.cur_line_split[0] == "eq":
            # // eq
            # // SP--
            # load_SP_into_A() // A = *y
            # D=M // D = y
            # // SP--
            # load_SP_into_A() // A = *x
            # D=M-D // D = x - y, 0 if eq is True
            # @eq0True // Load instruction for true case into the A register
            # D;JEQ // If D == 0, x==y. Jump over the false case to the true
            # // false case: if x != y
            # load_SP_into_A() // A = *x
            # M=0 // RAM[*x] = 0 (False)
            # @eq0TrueEnd // Load instruction to skip the true case
            # 0;JMP // jump over the true case
            # (eq0True) // true case: elif x == y
            # load_SP_into_A() // A = *x
            # M=-1 // RAM[*x] = -1 (True)
            # (eq0TrueEnd)
            # // SP++
            self.output_file.write(f"// eq\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"D=M\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"D=M-D\n")
            self.output_file.write(f"@eq{self.eq_num}True\n")
            self.output_file.write(f"D;JEQ\n")
            self.load_SP_into_A()
            self.output_file.write(f"M=0\n")
            self.output_file.write(f"@eq{self.eq_num}TrueEnd\n")
            self.output_file.write(f"0;JMP\n")
            self.output_file.write(f"(eq{self.eq_num}True)\n")
            self.load_SP_into_A()
            self.output_file.write(f"M=-1\n")
            self.output_file.write(f"(eq{self.eq_num}TrueEnd)\n")
            self.SP_pp(load_SP_into_A=False)
            self.eq_num += 1
        elif parser.cur_line_split[0] == "gt":
            # // gt
            # // SP--
            # load_SP_into_A() // A = *y
            # D=M // D = y
            # // SP--
            # load_SP_into_A() // A = *x
            # D=M-D // D = x - y, positive if gt is True
            # @gt0True // Load instruction for true case into the A register
            # D;JGT // If D == 0, x==y. Jump over the false case to the true
            # // false case: if x != y
            # load_SP_into_A() // A = *x
            # M=0 // RAM[*x] = 0 (False)
            # @gt0TrueEnd // Load instruction to skip the true case
            # 0;JMP // jump over the true case
            # (gt0True) // true case: elif x == y
            # load_SP_into_A() // A = *x
            # M=-1 // RAM[*x] = -1 (True)
            # (gt0TrueEnd)
            # // SP++
            self.output_file.write(f"// gt\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"D=M\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"D=M-D\n")
            self.output_file.write(f"@gt{self.gt_num}True\n")
            self.output_file.write(f"D;JGT\n")
            self.load_SP_into_A()
            self.output_file.write(f"M=0\n")
            self.output_file.write(f"@gt{self.gt_num}TrueEnd\n")
            self.output_file.write(f"0;JMP\n")
            self.output_file.write(f"(gt{self.gt_num}True)\n")
            self.load_SP_into_A()
            self.output_file.write(f"M=-1\n")
            self.output_file.write(f"(gt{self.gt_num}TrueEnd)\n")
            self.SP_pp(load_SP_into_A=False)
            self.gt_num += 1
        elif parser.cur_line_split[0] == "lt":
            # // lt
            # // SP--
            # load_SP_into_A() // A = *y
            # D=M // D = y
            # // SP--
            # load_SP_into_A() // A = *x
            # D=M-D // D = x - y, positive if lt is True
            # @lt0True // Load instruction for true case into the A register
            # D;JLT // If D == 0, x==y. Jump over the false case to the true
            # // false case: if x != y
            # load_SP_into_A() // A = *x
            # M=0 // RAM[*x] = 0 (False)
            # @lt0TrueEnd // Load instruction to skip the true case
            # 0;JMP // jump over the true case
            # (lt0True) // true case: elif x == y
            # load_SP_into_A() // A = *x
            # M=-1 // RAM[*x] = -1 (True)
            # (lt0TrueEnd)
            # // SP++
            self.output_file.write(f"// lt\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"D=M\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"D=M-D\n")
            self.output_file.write(f"@lt{self.lt_num}True\n")
            self.output_file.write(f"D;JLT\n")
            self.load_SP_into_A()
            self.output_file.write(f"M=0\n")
            self.output_file.write(f"@lt{self.lt_num}TrueEnd\n")
            self.output_file.write(f"0;JMP\n")
            self.output_file.write(f"(lt{self.lt_num}True)\n")
            self.load_SP_into_A()
            self.output_file.write(f"M=-1\n")
            self.output_file.write(f"(lt{self.lt_num}TrueEnd)\n")
            self.SP_pp(load_SP_into_A=False)
            self.lt_num += 1
        elif parser.cur_line_split[0] == "and":
            # // and
            # // SP--
            # load_SP_into_A() // A = *y
            # D=M // D = y
            # // SP--
            # load_SP_into_A() // A = *x
            # M=D&M // RAM[*x] = x - y
            # // SP++
            self.output_file.write(f"// and\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"D=M\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"M=D&M\n")
            self.SP_pp(load_SP_into_A=False)
        elif parser.cur_line_split[0] == "or":
            # // or
            # // SP--
            # load_SP_into_A() // A = *y
            # D=M // D = y
            # // SP--
            # load_SP_into_A() // A = *x
            # M=D|M // RAM[*x] = x - y
            # // SP++
            self.output_file.write(f"// or\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"D=M\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"M=D|M\n")
            self.SP_pp(load_SP_into_A=False)
        elif parser.cur_line_split[0] == "not":
            # // not
            # // SP--
            # load_SP_into_A() // A = *y
            # M=!M // Negate y and replace it with it's negated value
            # // Increment stack pointer, now at 258
            self.output_file.write(f"// neg\n")
            self.SP_mm(load_SP_into_A=True)
            self.output_file.write(f"M=!M\n")
            self.SP_pp(load_SP_into_A=False)

        self.output_file.write(f'\n')

    def push_constant(self, val: str):
        '''
        Pushes val onto the stack and incrememnts the stack pointer
        '''
        self.output_file.write(f"@{val}\n")
        self.output_file.write(f"D=A\n")
        self.load_SP_into_A()
        self.output_file.write(f"M=D\n")
        self.SP_pp(load_SP_into_A=False)

    def write_push(self, parser: Parser):
        # Write a comment with the VM code for reference/debugging
        self.output_file.write(f"// {' '.join(parser.cur_line_split[:3])}\n")

        if parser.cur_line_split[1] == "constant":
            self.push_constant(parser.cur_line_split[2])
        elif parser.cur_line_split[1] == "local":
            self.push_value("LCL", int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "argument":
            self.push_value("ARG", int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "static":
            self.push_pointer(self.static_symbol(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "temp":
            self.push_pointer(f"{self.TEMP + int(parser.cur_line_split[2])}")
        elif parser.cur_line_split[1] == "pointer":
            if int(parser.cur_line_split[2]) == 0:
                self.push_pointer("THIS")
            elif int(parser.cur_line_split[2]) == 1:
                self.push_pointer("THAT")
            else:
                raise RuntimeError("Invalid command")
        elif parser.cur_line_split[1] == "this":
            self.push_value("THIS", int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "that":
            self.push_value("THAT", int(parser.cur_line_split[2]))
        else:
            raise RuntimeError(f"Unkown push command: {' '.join(parser.cur_line_split[:3])}")

        self.output_file.write(f'\n')

    def write_pop(self, parser: Parser):
        self.output_file.write(f"// {' '.join(parser.cur_line_split[:3])}\n")

        if parser.cur_line_split[1] == "local":
            self.pop_value("LCL", int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "argument":
            self.pop_value("ARG", int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "static":
            self.pop_pointer(self.static_symbol(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "temp":
            self.pop_pointer(f"{self.TEMP + int(parser.cur_line_split[2])}")
        elif parser.cur_line_split[1] == "pointer":
            if int(parser.cur_line_split[2]) == 0:
                self.pop_pointer("THIS")
            elif int(parser.cur_line_split[2]) == 1:
                self.pop_pointer("THAT")
            else:
                raise RuntimeError("Invalid command")
        elif parser.cur_line_split[1] == "this":
            self.pop_value("THIS", int(parser.cur_line_split[2]))
        elif parser.cur_line_split[1] == "that":
            self.pop_value("THAT", int(parser.cur_line_split[2]))
        else:
            raise RuntimeError(f"Unkown pop command: {' '.join(parser.cur_line_split[:3])}")

        self.output_file.write(f'\n')
        pass

    def write_label(self, parser: Parser):
        self.output_file.write(f"({self.prefix_w_cur_func_name(parser.cur_line_split[1])})\n")
        self.output_file.write(f"\n")

    def write_if_goto(self, parser: Parser):
        '''
        This command effects a conditional goto operation. The stack’s topmost value is popped; 
        if the value is not zero, execution continues from the location marked by the label; otherwise, 
        execution continues from the next command in the program. 
        The jump destination must be located in the same function.
        '''
        self.output_file.write(f"// {' '.join(parser.cur_line_split[:2])}\n")
        # Pop the top of the stack into the D register
        self.SP_mm(load_SP_into_A=True)
        self.output_file.write(f"D=M\n")

        # Load the label into the A register
        self.output_file.write(f"@{self.prefix_w_cur_func_name(parser.cur_line_split[1])}\n")

        # If popped value is non-zero, jump to label
        self.output_file.write(f"D;JNE\n")

        self.output_file.write(f"\n")

    def write_goto(self, parser: Parser):
        '''
        This command effects an unconditional goto operation, causing execution to continue from the location marked by the label. 
        The jump destination must be located in the same function.
        '''
        self.output_file.write(f"// {' '.join(parser.cur_line_split[:2])}\n")

        self.goto_label(self.prefix_w_cur_func_name(parser.cur_line_split[1]))

        self.output_file.write(f"\n")

    def write_function(self, parser: Parser):
        '''
        `function f k`: declaring a function `f` that has `k` local variables
        ```psuedocode
        (f)
        repeat k times:
        PUSH 0
        ```
        See Figure 8.5 on p. 193
        '''
        self.output_file.write(f"// {' '.join(parser.cur_line_split[:2])}\n")

        f = parser.cur_line_split[1]
        k = int(parser.cur_line_split[2])

        # Set cur_func_name so that labels within this function can be written as `functionName$label`
        self.cur_func_name = f

        # Create label for the function itself
        self.output_file.write(f"({f})\n")

        # Push k local variables onto the stack, initialized to 0
        for _ in range(k):
            self.push_constant("0")

    def write_call(self, parser: Parser):
        '''
        `call f n`: calling a function `f` after `n` arguments have been pushed onto the stack
        ```psuedocode
        push return-address
        push LCL
        push ARG
        push THIS
        push THAT
        ARG = SP - n - 5
        LCL = SP
        goto f
        (return-address)
        ```
        See Figure 8.5 on p. 193
        '''
        self.output_file.write(f"// {' '.join(parser.cur_line_split[:3])}\n")

        f = parser.cur_line_split[1]
        n = parser.cur_line_split[2]
        ret_addr = self.create_ret_addr(self.cur_func_name, parser.cur_line_number)

        self.push_constant(ret_addr)
        self.push_pointer("LCL")
        self.push_pointer("ARG")
        self.push_pointer("THIS")
        self.push_pointer("THAT")

        # Load SP into the D register
        self.load_SP_into_D()

        # LCL = SP
        self.output_file.write(f"@LCL\n")
        self.output_file.write(f"M=D\n")

        # ARG = SP - n - 5
        self.output_file.write(f"@{n}\n")
        self.output_file.write(f"D=D-A\n")
        self.output_file.write(f"@{5}\n")
        self.output_file.write(f"D=D-A\n")
        self.output_file.write(f"@ARG\n")
        self.output_file.write(f"M=D\n")

        self.goto_label(f)
        self.output_file.write(f"({ret_addr})\n")
        self.output_file.write(f"\n")

    def write_return(self, parser: Parser):
        '''
        ```psuedocode
        FRAME = LCL
        RET = *(FRAME - 5)
        *ARG = pop()
        SP = ARG+1
        THAT = *(FRAME - 1)
        THIS = *(FRAME - 2)
        ARG = *(FRAME - 3)
        LCL = *(FRAME - 4)
        goto RET
        ```
        See Figure 8.5 on p. 193
        '''
        self.output_file.write(f"// {parser.cur_line_split[0]}\n")
        # FRAME = LCL
        # RET = *(FRAME - 5)
        self.output_file.write(f"@LCL\n")
        self.output_file.write(f"D=M\n")
        self.output_file.write(f"@5\n")
        self.output_file.write(f"A=D-A\n")
        self.output_file.write(f"D=M\n")
        # Save in register 14
        self.output_file.write(f"@R14\n")
        self.output_file.write(f"M=D\n")

        # *ARG = pop()
        # reposition the return value from the current top of stack to the caller's top of stack (which is the current ARG)
        self.pop_value("ARG", 0)

        # SP = ARG+1
        self.output_file.write(f"@ARG\n")
        self.output_file.write(f"D=M+1\n")
        self.output_file.write(f"@SP\n")
        self.output_file.write(f"M=D\n")

        # THAT = *(FRAME - 1)
        # THIS = *(FRAME - 2)
        # ARG = *(FRAME - 3)
        # LCL = *(FRAME - 4)
        for destination in ['THAT', 'THIS', 'ARG', 'LCL']:
            self.output_file.write(f"@LCL\n")
            # LCL = FRAME-1 so that the loop works, A = FRAME-1
            self.output_file.write(f"AM=M-1\n")
            # D = *(FRAME-1), the saved value
            self.output_file.write(f"D=M\n")
            # destination = D
            self.output_file.write(f"@{destination}\n")
            self.output_file.write(f"M=D\n")

        # goto RET
        self.output_file.write(f"@R14\n")
        self.output_file.write(f"A=M\n")
        self.output_file.write(f"0;JMP\n")

        self.output_file.write(f"\n")


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
                if filename.split('.')[-1] == 'vm':
                    self.parsers.append(Parser(f"{directory_or_filename}/{filename}"))

        # Create CodeWriter
        self.codewriter = CodeWriter(directory_or_filename.split('.')[0].strip('/') + '.asm')

    def run(self):
        '''
        Main function that calls each parser to run, passing each parsed token into codewriter which then 
        writes the corresponding assembly code to the output file.
        '''
        for parser in self.parsers:
            # Update cur_parser_filename so codewriter knows how to name static vars (section 7.3 in the book)
            self.codewriter.cur_parser_filename = parser.filename
            while (parser.advance()):
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
                elif parser.command_type == CT.LABEL:
                    self.codewriter.write_label(parser)
                elif parser.command_type == CT.IF_GOTO:
                    self.codewriter.write_if_goto(parser)
                elif parser.command_type == CT.GOTO:
                    self.codewriter.write_goto(parser)
                elif parser.command_type == CT.FUNCTION:
                    self.codewriter.write_function(parser)
                elif parser.command_type == CT.RETURN:
                    self.codewriter.write_return(parser)
                elif parser.command_type == CT.CALL:
                    self.codewriter.write_call(parser)

        # Don't forget to close the output file when you're done
        self.codewriter.output_file.close()


if __name__ == "__main__":
    import sys
    vmt = VMtranslator(sys.argv[1])
    vmt.run()