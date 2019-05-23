from numpy import binary_repr


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


# strings for function values in the format (zx)(nx)(zy)(ny)(f)(no), taken from page 56
alu_functions = [
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
