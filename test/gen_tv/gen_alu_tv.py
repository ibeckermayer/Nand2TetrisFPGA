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
from util import int_to_bin_str
from random import randint

OUTPUT_FILE = 'tvs/ALU.tv'  # expects to be run from directory above this
N = 1000  # number of times each function will be checked with random inputs for x and y

# strings for function values in the format (zx)(nx)(zy)(ny)(f)(no), taken from page 56
functions = [
    '101010',  # f(x,y) = 0
    '111111',  # f(x,y) = 1
    '111010',  # f(x,y) = -1
    '001100',  # f(x,y) = x
    '110000',  # f(x,y) = y
    '001101',  # f(x,y) = !x
    '110001',  # f(x,y) = !y
    '001111',  # f(x,y) = -x
    '110011',  # f(x,y) = -y
    '011111',  # f(x,y) = x+1
    '110111',  # f(x,y) = y+1
    '001110',  # f(x,y) = x-1
    '110010',  # f(x,y) = y-1
    '000010',  # f(x,y) = x+y
    '010011',  # f(x,y) = x-y
    '000111',  # f(x,y) = y-x
    '000000',  # f(x,y) = x&y
    '010101'  # f(x,y) = x|y
]


def gen_rand_int() -> int:
    return randint(-32768, 32767)


def bitwise_not(val: int) -> str:
    '''
    returns bitwise negated str representation of 16 bit integer
    '''
    strval = int_to_bin_str(val, 16)
    notval = ''
    for dig in strval:
        if dig == '0':
            notval += '1'
        else:
            notval += '0'
    return notval


def bitwise_and(x: int, y: int) -> str:
    strx = int_to_bin_str(x, 16)
    stry = int_to_bin_str(y, 16)
    andval = ''
    for i in range(16):
        if (strx[i] == '1' and stry[i] == '1'):
            andval += '1'
        else:
            andval += '0'
    return andval


def bitwise_or(x: int, y: int) -> str:
    strx = int_to_bin_str(x, 16)
    stry = int_to_bin_str(y, 16)
    orval = ''
    for i in range(16):
        if (strx[i] == '1' or stry[i] == '1'):
            orval += '1'
        else:
            orval += '0'
    return orval


def zr(out_str: str) -> str:
    if out_str == int_to_bin_str(0, 16):
        return '1'
    return '0'


def ng(out_str: str) -> str:
    return out_str[0]


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
    for function in functions:
        for i in range(N):
            f.write(
                ALU_logic_and_build_line(gen_rand_int(), gen_rand_int(),
                                         function))
