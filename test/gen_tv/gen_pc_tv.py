'''
script for generating the testvector file for use in PC_tb.v

Format will be:
(in_v[16])(reset)(load)(inc)_(out_v[16]) in binary format

Algorithm:
- start with a reset, expect a 0
- count up to 1000 to check increment
- reset again to check reset at arbitrary count
- count up to 1000 again
- load 65000 to test loading from arbitrary count
- count up to 65535
- count up to 1000 again to check rollover
'''


def u_int_to_bin_str(num: int, width: int) -> str:
    '''
    converts a positive integer to an unsigned binary string representation
    num: number to be converted to binary
    width: bit-width of the output string (e.g. 16-bits, 1-bit)
    '''
    binary_string = '{:b}'.format(num)
    full_width = len(binary_string)
    if (width <= full_width):
        return binary_string[-width:]
    else:
        return ('0' * (width - full_width)) + binary_string


def PC_logic(in_v: int, reset: int, load: int, inc: int,
             previous_out_v: int) -> int:
    '''
    implements expected logic of PC module to generate test file
    returns integer that represents output of PC module for this input logic
    '''
    if (reset):
        return 0
    else if (load):
        return in_v
    else if (inc):
        return previous_out_v + 1
    else:
        return previous_out_v


# initial values
in_v = 1
reset = 1  # start with a reset
load = 1  # start with load 1, see that reset has precedence
inc = 1  # start with inc 1, see that reset has precedence
out_v = 0  # output should be 0 due to reset

