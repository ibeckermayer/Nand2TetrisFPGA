'''
Program to simulate the Hack platform based on assembly input
'''
from enum import Enum
from typing import TextIO, List, Dict, Optional
from dataclasses import dataclass


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


class Parser:
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


if __name__ == "__main__":
    import sys
    p = Parser(sys.argv[1])
    print(p.run())
