'''
Program to simulate the Hack platform based on assembly input
'''
from enum import Enum
from typing import TextIO, List, Dict, Optional
from dataclasses import dataclass
from numpy import int16, uint16
from collections.abc import MutableSequence


class CommandType(Enum):
    A_COMMAND = 1
    C_COMMAND = 2
    L_COMMAND = 3  # Pseudo instruction for i.e. `(LOOP)`
    SKIP = 4  # Pseudo instruction for empty lines or fully commented lines
    EOF = 5  # Pseudo instruction to denote we're at the EOF


CT = CommandType  # Alias


@dataclass
class Instruction:
    '''
    val for type=A_COMMAND looks like {'val': '45'} for `@45`
    val for type=C_COMMAND looks like {'dest': 'D', 'comp': 'D-A', 'jump': ''} for `D=D-A`
    '''
    type: CommandType
    val: Dict[str, str]


class AsmParser:
    '''
    Used to step through a Hack assembly file and generate a list of Instruction to be fed into the HackExecutor
    '''

    def __init__(self, filename: str):
        self.filename = filename
        self.file: TextIO = open(self.filename, 'r')
        self.cur_line: str
        self.cur_line_num: int = 0  # Only increments on executable lines, i.e. not comments/blank lines/L commands
        self.line_type: CommandType
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

    def advance(self) -> Instruction:
        '''
        To be called in a loop to advance through the file
        returning None represents EOF
        '''
        # Read the next line and strip leading/trailing spaces
        self.cur_line = self.file.readline().strip(' ')

        # If this is the end of the file, return EOF instruction
        if self.cur_line == '':
            return Instruction(CT.EOF, {})

        # If this is a full line comment or empty line return a SKIP
        if self.cur_line == '\n' or self.cur_line[:2] == "//":
            return Instruction(CT.SKIP, {})

        # Strip away all comments and newlines so the line is only the first token
        # i.e. "@45 // load 45 into the A reg\n" becomes "@45"
        self.cur_line = self.cur_line.split('//')[0].strip()
        # Build the instruction
        instruction = self.__build_instruction(self.cur_line, self.cur_line_num)

        return instruction

    def __build_instruction(self, cur_line: str, cur_line_num: int) -> Instruction:
        '''
        cur_line has been stripped of all extra comments and spaces
        If this is an L command, updates the symbol table and returns None
        Otherwise deciphers the instruction and returns a corresponding Instruction object
        '''

        if cur_line[0] == "(":  # L_COMMAND, i.e. `(LOOP)`
            symbol = cur_line[1:-1]  # extract symbol
            self.symbol_table[symbol] = str(cur_line_num)  # update symbol table
            return Instruction(CT.L_COMMAND, {})
        elif cur_line[0] == '@':  # A_COMMAND
            # increment the cur_line_num to track future L_COMMAND's
            self.cur_line_num = cur_line_num + 1

            if cur_line[1].isdigit():
                # Easy case: the next digit is numerical
                return Instruction(CT.A_COMMAND, {'val': cur_line[1:]})
            elif cur_line[1].isalpha():
                # Trickier case: the next digit is a letter, so we're dealing with a symbol
                sym_val = self.symbol_table.get(cur_line[1:])
                if sym_val is not None:
                    # This symbol has already been declared and registered in the symbol table,
                    # return an A_COMMAND with its value
                    return Instruction(CT.A_COMMAND, {'val': sym_val})
                else:
                    # This is a new symbol, assign it a value
                    new_sym = cur_line[1:]
                    new_sym_val = str(self.next_unknown_symbol)
                    self.symbol_table[new_sym] = new_sym_val
                    self.next_unknown_symbol += 1
                    return Instruction(CT.A_COMMAND, {'val': new_sym_val})
        else:  # C_COMMAND `dest=comp;jump`, dest and jump are optional
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
            return Instruction(CT.C_COMMAND, {'dest': dest, 'comp': comp, 'jump': jump})

        # Should never reach here
        raise RuntimeError(
            "__build_instruction reached a point in the control flow it shouldn't have!")
        return Instruction(CT.SKIP, {})

    def run(self) -> List[Instruction]:
        '''
        Walk through the input file and generate a list of instructions to be passed into the HackExecutor
        '''
        instructions: List[Instruction] = []
        while True:
            instruction: Instruction = self.advance()
            if instruction.type == CT.A_COMMAND or instruction.type == CT.C_COMMAND:
                instructions.append(instruction)
            elif instruction.type == CT.EOF:
                break

        # Should never reach here
        return instructions


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

    def step(self):
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
        else:
            raise RuntimeError(f"Instruction of type {ins.type} somehow got into the executor")

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
