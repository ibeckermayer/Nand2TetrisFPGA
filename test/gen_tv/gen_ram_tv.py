'''
script for generating the testvector file for use in RAM_tb.v

Format will be:
(address[16])(in_v[16])(load)_(out_v[16]) in binary format

Algorithm:
- address and in_v and out_v = 0..32767, load = 1: load each address with its value (test load)
- go through again with load, in_v = 0 and read each value (test read only)
- load 65535 into register 500, next read it out (largest value into arbitrary register)
- load 1 into register 0, read it out (setting up to test rollover)
- load 65536 into register 32768 (65535 rollsover to 0, 32768 is out of bounds but gets addressed to highest address of 32767)
- explicit read of register 0 to confirm no rollover
- confirm highest address register got 0. addressing doesn't rollover,
  out of bounds gets sent to the highest number register
'''
from util import int_to_bin_str

OUTPUT_FILE = 'tvs/RAM.tv'  # expects to be run from directory above this


def build_line(address: int, in_v: int, load: int, out_v: int) -> str:
    '''
    builds a line for output to the tv file
    '''
    return int_to_bin_str(address,
                          16) + int_to_bin_str(in_v, 16) + int_to_bin_str(
                              load, 1) + '_' + int_to_bin_str(out_v, 16) + '\n'


with open(OUTPUT_FILE, 'w') as f:
    # address and in_v and out_v = 0..32767, load = 1: load each address with its value (test load)
    for i in range(32768):
        f.write(build_line(i, i, 1, i))

    # go through again with load, in_v = 0 and read each value (test read only)
    for i in range(32768):
        f.write(build_line(i, 0, 0, i))

    # load 65535 into register 500, next read it out (largest value into arbitrary register)
    f.write(build_line(500, 65535, 1, 65535))
    f.write(build_line(500, 0, 0, 65535))

    # load 1 into register 0, read it out (setting up to test rollover)
    f.write(build_line(0, 1, 1, 1))
    f.write(build_line(0, 0, 0, 1))

    # load 65536 into register 32768 (65535 rollsover to 0, 32768 is out of bounds but gets addressed to highest address of 32767)
    f.write(build_line(32768, 65536, 1, 65536))

    # explicit read of register 0 to confirm no rollover
    f.write(build_line(0, 0, 0, 1))

    # confirm highest address register got 0. addressing doesn't rollover,
    # out of bounds gets sent to the highest number register
    f.write(build_line(32767, 0, 0, 0))