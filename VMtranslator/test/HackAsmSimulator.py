'''
Program to simulate the Hack platform based on assembly input
'''
from enum import Enum
from typing import TextIO, List, Dict, Optional
from dataclasses import dataclass
from numpy import int16, uint16
from collections.abc import MutableSequence


class LineType(Enum):
    A_COMMAND = 1
    C_COMMAND = 2
    L_COMMAND = 3
    SKIP = 4
    EOF = 5


LT = LineType


class CommandType(Enum):
    A_COMMAND = 1
    C_COMMAND = 2
    END = 3  # psuedo command, signals simulation to stop


CT = CommandType  # Alias


@dataclass
class Instruction:
    '''
    val for type=A_COMMAND looks like {'val': '45'} for `@45`
    val for type=C_COMMAND looks like {'dest': 'D', 'comp': 'D-A', 'jump': ''} for `D=D-A`
    val for type=END is just an empty Dict
    line: the line in the .asm file corresponding to this instruction, useful for debugging
    line_num: the line number in the .asm file corresponding to this instruction, useful for debugging
    '''
    type: CommandType
    val: Dict[str, str]
    line: str
    line_num: int


class AsmParser:
    '''
    Used to step through a Hack assembly file and generate a list of Instruction to be fed into the HackExecutor

    Two pass parser: first steps through the file and builds the symbol table, then steps through again to create the instructions
    '''

    def __init__(self, filename: str):
        self.filename = filename
        self.file: TextIO = open(self.filename, 'r')
        self.cur_line: str
        self.asm_code_line_num: int = 0  # The line number in the .asm file
        # The next machine code line number. Gets incremented after each A_COMMAND or C_COMMAND,
        # so that subsequent L_COMMAND's know which line they should be set to
        self.machine_code_line_num: int = 0  # Only increments on executable lines, i.e. only A_COMMAND's and C_COMMAND's; 0 indexed
        self.line_type: LineType
        self.first_pass: bool = True  # True for first pass, false for second pass
        self.symbol_table: Dict[str, str] = {
            'R0': '0',
            'R1': '1',
            'R2': '2',
            'R3': '3',
            'R4': '4',
            'R5': '5',
            'R6': '6',
            'R7': '7',
            'R8': '8',
            'R9': '9',
            'R10': '10',
            'R11': '11',
            'R12': '12',
            'R13': '13',
            'R14': '14',
            'R15': '15',
            'SP': '0',
            'LCL': '1',
            'ARG': '2',
            'THIS': '3',
            'THAT': '4'
        }
        self.next_unknown_symbol = 16
        self.instructions: List[Instruction] = []  # Filled out in the second pass

    def advance(self) -> str:
        '''
        Advances to the next line of the file and updates the state of the parser.
        If this is the first pass, updates the symbol table as appropriate.
        If this is the second pass, replaces symbols with their respective values as appropriate.
        '''
        # Read the next line and strip leading/trailing spaces
        self.cur_line = self.file.readline().strip(' ')
        self.asm_code_line_num += 1

        # EOF
        if self.cur_line == '':
            self.line_type = LT.EOF
            # If this is the first pass, reset file, set first_pass flag to false, and reset machine_code_line_num and asm_code_line_num
            if self.first_pass:
                self.file.seek(0)
                self.first_pass = False
                self.machine_code_line_num = 0
                self.asm_code_line_num = 0
            else:
                self.append_instruction(CT.END, {})
            return self.cur_line

        # If this is a full line comment or empty line return a SKIP
        if self.cur_line == '\n' or self.cur_line[:2] == "//":
            self.line_type = LT.SKIP
            return self.cur_line

        # Strip away all comments and newlines so the line is only the first token
        # i.e. "@45 // load 45 into the A reg\n" becomes "@45"
        self.cur_line = self.cur_line.split('//')[0].strip()

        # L_COMMAND
        if self.cur_line[0] == "(":
            self.line_type = LT.L_COMMAND
            if self.first_pass:
                # On first pass, update symbol table with L_COMMAND's
                self.symbol_table[self.cur_line[1:-1]] = str(self.machine_code_line_num)
            return self.cur_line

        # A_COMMAND
        if self.cur_line[0] == '@':
            self.line_type = LT.A_COMMAND
            if self.first_pass:
                if self.cur_line[1].isalpha():
                    # On first pass, if this is a symbol, update the symbol table iif the symbol dne
                    if self.symbol_table.get(self.cur_line[1:]) == None:
                        self.symbol_table[self.cur_line[1:]] = str(self.next_unknown_symbol)
                        self.next_unknown_symbol += 1
            else:
                self.append_A_instr(self.cur_line)
            self.machine_code_line_num += 1
            return self.cur_line
        else:  # C_COMMAND
            self.line_type = LT.C_COMMAND
            if self.first_pass:
                pass
            else:
                self.append_C_instr(self.cur_line)
            self.machine_code_line_num += 1
            return self.cur_line

        # Build the instruction
        instruction = self.__build_instruction(self.cur_line, self.machine_code_line_num)

        return ''

    def append_instruction(self, type: CommandType, val: Dict[str, str]):
        '''
        Appends an Instruction to self.instructions
        '''
        self.instructions.append(Instruction(type, val, self.cur_line, self.asm_code_line_num))

    def append_A_instr(self, cur_line: str):
        if self.is_infinite_loop_terminator(cur_line):
            return
        val = cur_line[1:]  # remove '@'
        if val[0].isalpha():
            val = self.symbol_table[val]
        self.append_instruction(CT.A_COMMAND, {'val': val})

    def append_C_instr(self, cur_line: str):
        dest = ''
        comp = ''
        jump = ''
        buf = cur_line
        if '=' in cur_line:
            split_buf = buf.split('=')
            # Extract dest
            dest = split_buf[0]
            # Remove it from buf
            buf = split_buf[-1]
        if ';' in cur_line:
            split_buf = buf.split(';')
            # Extract jump
            jump = split_buf[-1]
            # Remove it from buf
            buf = split_buf[0]
        # Only comp is left over
        comp = buf

        self.append_instruction(CT.C_COMMAND, {'dest': dest, 'comp': comp, 'jump': jump})

    def is_infinite_loop_terminator(self, cur_line: str):
        '''
        NOTE: Major hack incoming

        A Sys.init VM function will typically setup and then call some other function, and then terminate with an infinite loop:
        {code[vm]
            function Sys.init 0
            push constant 4
            call Main.fibonacci 1   // computes the 4'th fibonacci element
            label INFLOOP
            goto INFLOOP            // loops infinitely
        }
        
        The assembly for the final two lines (the infinite loop) will wind up looking like the following
        {code[asm]
            (Sys.init$INFLOOP)
            @Sys.init$INFLOOP
            0;JMP
        }

        This makes sense for ordinary Hack usage, but for the purposes of our testing simulator (HackExecutor), we want the program to
        terminate or else the tests themselves will run in an infinite loop and never complete. To avoid this, we will designate the
        label `Sys.init$INFLOOP` a special keyword in our testing setup, and use it to indicate that we've reached the infinte loop
        and so should terminate the program (denoted by a CT.END instruction)
        '''
        if cur_line == "@Sys.init$INFLOOP":
            self.append_instruction(CT.END, {})
            return True
        return False

    def run(self) -> List[Instruction]:
        '''
        Walk through the input file and generate a list of instructions to be passed into the HackExecutor
        '''
        # Run first pass to fill out symbol table
        while self.advance():
            continue

        # Run second pass to fill out list of instructions
        while self.advance():
            continue

        # Return the list
        return self.instructions


class RAM32K(MutableSequence):
    '''
    Generally works like a list, but enforces the 15-bit width of the address space.
    I.e. ram[-32768] becomes ram[0]. This is in case we have a negative value in the A register
    and our simulator has a C_COMMAND that asks for `M` (either dest or comp).

    Also casting all values stored in the array to int16 for good measure
    '''

    def __init__(self):
        self._mem = [int16(0)] * 2**15

    def __to_uint15(self, index) -> int:
        '''
        To be called before any method that makes use of index
        '''
        # Cast index to uint16, format that to a 16-bit binary string representation,
        # lop off the MSB (with [-15:]), then cast that 15-bit binary string to an int
        return int(format(uint16(index), '016b')[-15:], 2)

    def __len__(self):
        return len(self._mem)

    def __delitem__(self, index):
        self._mem.__delitem__(self.__to_uint15(index))

    def insert(self, index, value):
        self._mem.insert(self.__to_uint15(index), int16(value))

    def __setitem__(self, index, value):
        self._mem.__setitem__(self.__to_uint15(index), int16(value))

    def __getitem__(self, index):
        return self._mem.__getitem__(self.__to_uint15(index))


class HackExecutor:
    '''
    Takes in the List[Instruction] generated by the AsmParser and then simulates them step by step
    '''

    def __init__(self, instructions: List[Instruction]):
        self.instructions = instructions  # Should only be A_COMMAND's and C_COMMAND's at this point
        self.pc: int16 = 0  # The program counter, used to index the instructions
        self.ram: RAM32K = RAM32K()  # 32k RAM initialized to 0
        self.A: int16 = 0  # A reg
        self.D: int16 = 0  # D reg
        self.ALU_output: int16 = 0

    # Give all numpy.int16 attributes getters and setters so that we don't need to remember to cast to int16 every set call.
    @property
    def pc(self) -> int16:
        return self._pc

    @pc.setter
    def pc(self, val):
        self._pc = int16(val)

    @property
    def A(self) -> int16:
        return self._A

    @A.setter
    def A(self, val):
        self._A = int16(val)

    @property
    def D(self) -> int16:
        return self._D

    @D.setter
    def D(self, val):
        self._D = int16(val)

    @property
    def ALU_output(self) -> int16:
        return self._ALU_output

    @ALU_output.setter
    def ALU_output(self, val):
        self._ALU_output = int16(val)

    def step(self) -> Instruction:
        '''
        Executes a single instruction
        '''
        ins: Instruction = self.instructions[self.pc]

        if ins.type == CT.A_COMMAND:
            # If this is an A command, load the value into the A register and increment the pc
            self.A = ins.val['val']
            self.pc += 1
        elif ins.type == CT.C_COMMAND:
            # NOTE: The order of the following commands matter, see respective docstrings
            self.__handle_comp(ins.val['comp'])
            self.__handle_jump(ins.val['jump'])
            self.__handle_dest(ins.val['dest'])
        # else: ins.type == CT.END
        return ins

    def __handle_comp(self, comp: str):
        '''
        Logic for handling the comp command, must be called before __handle_jump and __handle_dest
        since this function sets the ALU_output that each of those depend on
        '''
        if comp == '0':
            self.ALU_output = 0
        elif comp == '1':
            self.ALU_output = 1
        elif comp == '-1':
            self.ALU_output = -1
        elif comp == 'D':
            self.ALU_output = self.D
        elif comp == 'A':
            self.ALU_output = self.A
        elif comp == '!D':
            self.ALU_output = ~self.D
        elif comp == '!A':
            self.ALU_output = ~self.A
        elif comp == '-D':
            self.ALU_output = -self.D
        elif comp == '-A':
            self.ALU_output = -self.A
        elif comp == 'D+1':
            self.ALU_output = self.D + 1
        elif comp == 'A+1':
            self.ALU_output = self.A + 1
        elif comp == 'D-1':
            self.ALU_output = self.D - 1
        elif comp == 'A-1':
            self.ALU_output = self.A - 1
        elif comp == 'D+A':
            self.ALU_output = self.D + self.A
        elif comp == 'D-A':
            self.ALU_output = self.D - self.A
        elif comp == 'A-D':
            self.ALU_output = self.A - self.D
        elif comp == 'D&A':
            self.ALU_output = self.D & self.A
        elif comp == 'D|A':
            self.ALU_output = self.D | self.A
        elif comp == 'M':
            self.ALU_output = self.ram[self.A]
        elif comp == '!M':
            self.ALU_output = ~self.ram[self.A]
        elif comp == '-M':
            self.ALU_output = -self.ram[self.A]
        elif comp == 'M+1':
            self.ALU_output = self.ram[self.A] + 1
        elif comp == 'M-1':
            self.ALU_output = self.ram[self.A] - 1
        elif comp == 'D+M':
            self.ALU_output = self.D + self.ram[self.A]
        elif comp == 'D-M':
            self.ALU_output = self.D - self.ram[self.A]
        elif comp == 'M-D':
            self.ALU_output = self.ram[self.A] - self.D
        elif comp == 'D&M':
            self.ALU_output = self.D & self.ram[self.A]
        elif comp == 'D|M':
            self.ALU_output = self.D | self.ram[self.A]
        else:
            raise RuntimeError(f"Unkown comp command: {comp}")

    def __handle_jump(self, jump: str):
        '''
        Only to be called after __handle_comp, and before __handle_dest because:
        1) it depends on self.ALU_output being already set for this step (set by __handle_comp)
        2) it depends on the self.A of the previous step (self.A is modified for this step in __handle_dest, so we can't do it first)

        jump is either an empty string or one of JGT, JEQ, JGE, JLT, JNE, JLE, JMP

        Since jump optionally sets the program counter self.pc, we'll also handle incrementing that counter
        here too
        '''
        if jump == 'JGT' and self.ALU_output > 0:
            self.pc = self.A
        elif jump == 'JEQ' and self.ALU_output == 0:
            self.pc = self.A
        elif jump == 'JGE' and self.ALU_output >= 0:
            self.pc = self.A
        elif jump == 'JLT' and self.ALU_output < 0:
            self.pc = self.A
        elif jump == 'JNE' and self.ALU_output != 0:
            self.pc = self.A
        elif jump == 'JLE' and self.ALU_output <= 0:
            self.pc = self.A
        elif jump == 'JMP':
            self.pc = self.A
        else:
            # Don't forget to increment the counter if the jump condition isn't met
            self.pc += 1

    def __handle_dest(self, dest: str):
        '''
        Only to be called after __handle_comp and __handle_jump since:
        1) it depends on self.ALU_output being already set for this step (set by __handle_comp)
        2) it modifies self.A for this step, but self.A from the previous step is needed for __handle_comp and __handle_jump

        dest is either an empty string for null or some combination of A, M, and D.

        Our assembler cares about the order of A/M/D, as specified in Figure 4.4 of the book.
        For the sake of expediency here, we'll simply ignore that
        '''
        # M must be set before A, since M means ram[self.A_from_previous_step]
        if 'M' in dest:
            self.ram[self.A] = self.ALU_output

        if 'A' in dest:
            self.A = self.ALU_output

        if 'D' in dest:
            self.D = self.ALU_output


if __name__ == "__main__":
    import sys

    hack = HackExecutor(AsmParser(sys.argv[1]).run())
    for _ in hack.instructions:
        hack.step()
