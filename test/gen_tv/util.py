from numpy import binary_repr
from random import randint


def int_to_bin_str(val: int, width: int) -> str:
    '''
    converts an integer to it's twos complement representation of width=width
    val: number to be converted to binary
    width: bit-width of the output string (e.g. 16-bits, 1-bit)
    '''
    retval = binary_repr(val, width)
    if len(retval) > width:
        return retval[1:]
    else:
        return retval


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
