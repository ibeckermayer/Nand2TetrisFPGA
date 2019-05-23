'''
script for generating the testvector file for use in ALU_tb.v

Format will be:
(x)(y)(zx)(nx)(zy)(ny)(f)(no)_(out_v)(zr)(ng)

Algorithm:
For each combination of zx, nx, zy, ny, f, no listed on page 56 of the pdf of the book:
- check with:
  - x = 0, y = 0
  - check N times with random integers in [-32768, 32767]
'''
from util import (int_to_bin_str, bitwise_not, bitwise_and, bitwise_or, zr, ng,
                  alu_functions)
from random import randint

OUTPUT_FILE = 'tvs/ALU.tv'  # expects to be run from directory above this
N = 1000  # number of times each function will be checked with random inputs for x and y


def ALU_logic_and_build_line(x: int, y: int, function: str) -> str:
    '''
    builds a line for ALU.tv
    combined with the logic since it's simpler to implement bitwise logic on string repr of binary numbers
    '''
    input_str = int_to_bin_str(x, 16) + int_to_bin_str(y, 16) + function + '_'
    out_str = None  # will be set to output value (in binary string repr) depending on function
    if function == '101010':  # f(x,y) = 0
        out_str = int_to_bin_str(0, 16)
    elif function == '111111':  # f(x,y) = 1
        out_str = int_to_bin_str(1, 16)
    elif function == '111010':  # f(x,y) = -1
        out_str = int_to_bin_str(-1, 16)
    elif function == '001100':  # f(x,y) = x
        out_str = int_to_bin_str(x, 16)
    elif function == '110000':  # f(x,y) = y
        out_str = int_to_bin_str(y, 16)
    elif function == '001101':  # f(x,y) = !x
        out_str = bitwise_not(x)
    elif function == '110001':  # f(x,y) = !y
        out_str = bitwise_not(y)
    elif function == '001111':  # f(x,y) = -x
        out_str = int_to_bin_str(-x, 16)
    elif function == '110011':  # f(x,y) = -y
        out_str = int_to_bin_str(-y, 16)
    elif function == '011111':  # f(x,y) = x+1
        out_str = int_to_bin_str(x + 1, 16)
    elif function == '110111':  # f(x,y) = y+1
        out_str = int_to_bin_str(y + 1, 16)
    elif function == '001110':  # f(x,y) = x-1
        out_str = int_to_bin_str(x - 1, 16)
    elif function == '110010':  # f(x,y) = y-1
        out_str = int_to_bin_str(y - 1, 16)
    elif function == '000010':  # f(x,y) = x+y
        out_str = int_to_bin_str(x + y, 16)
    elif function == '010011':  # f(x,y) = x-y
        out_str = int_to_bin_str(x - y, 16)
    elif function == '000111':  # f(x,y) = y-x
        out_str = int_to_bin_str(y - x, 16)
    elif function == '000000':  # f(x,y) = x&y
        out_str = bitwise_and(x, y)
    elif function == '010101':  # f(x,y) = x|y
        out_str = bitwise_or(x, y)
    else:
        raise Exception('Unkown function')
    return input_str + out_str + zr(out_str) + ng(out_str) + '\n'


with open(OUTPUT_FILE, 'w') as f:
    for function in alu_functions:
        for i in range(N):
            f.write(
                ALU_logic_and_build_line(randint(-32768, 32767),
                                         randint(-32768, 32767), function))
