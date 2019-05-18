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
from util import u_int_to_bin_str

OUTPUT_FILE = 'tvs/PC.tv'  # expects to be run from directory above this


def PC_logic(in_v: int, reset: int, load: int, inc: int,
             previous_out_v: int) -> int:
    '''
    implements expected logic of PC module to generate test file
    returns integer that represents output of PC module for this input logic
    '''
    if (reset):
        return 0
    elif (load):
        return in_v
    elif (inc):
        return previous_out_v + 1
    else:
        return previous_out_v


def build_line(in_v: int, reset: int, load: int, inc: int,
               previous_out_v: int):
    '''
    builds a line for output to the tv file
    '''
    next_out = PC_logic(in_v, reset, load, inc, previous_out_v)
    line = u_int_to_bin_str(in_v, 16) + u_int_to_bin_str(
        reset, 1) + u_int_to_bin_str(load, 1) + u_int_to_bin_str(
            inc, 1) + '_' + u_int_to_bin_str(next_out, 16) + '\n'
    return (line, next_out)


# initial values
in_v = 1
reset = 1  # start with a reset
load = 1  # start with load 1, see that reset has precedence
inc = 1  # start with inc 1, see that reset has precedence
out_v = 0  # output should be 0 due to reset

with open(OUTPUT_FILE, 'w') as f:
    # start with a reset, expect a 0
    line, out_v = build_line(in_v, reset, load, inc, out_v)
    f.write(line)

    # count up to 1000 to check increment
    reset = 0
    load = 0
    for i in range(1000):
        line, out_v = build_line(in_v, reset, load, inc, out_v)
        f.write(line)

    # reset again to check reset at arbitrary count
    reset = 1
    line, out_v = build_line(in_v, reset, load, inc, out_v)
    f.write(line)

    # count up to 1000 again
    reset = 0
    load = 0
    for i in range(1000):
        line, out_v = build_line(in_v, reset, load, inc, out_v)
        f.write(line)

    # load 65000 to test loading from arbitrary count
    load = 1
    in_v = 65000
    line, out_v = build_line(in_v, reset, load, inc, out_v)
    f.write(line)

    # count up to 65535
    # count up to 1000 again to check rollover
    load = 0
    for i in range(1001 + 535):
        line, out_v = build_line(in_v, reset, load, inc, out_v)
        f.write(line)
