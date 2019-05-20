'''
script for generating the testvector file for use in ROM_tb.v

ROM_MEMORY_INPUT Algorithm:
- generate input for ROM that begins with 32767 in address 0, and ends with 0 in address 32767

Format will be:
(address)_(out_v)

Algorithm:
- go through each address in ram and expect 32767 - address
'''
from util import u_int_to_bin_str

# expects to be run from directory above this
ROM_MEMORY_INPUT = 'tvs/ROM_input.tv'  # ROM32K must be loaded from a file
OUTPUT_FILE = 'tvs/ROM.tv'  # file for the testbench

# generate ROM_MEMORY_INPUT
with open(ROM_MEMORY_INPUT, 'w') as f:
    for i in range(32768):
        f.write(u_int_to_bin_str(32767 - i, 16) + '\n')


def build_line(address: int, out_v: int) -> str:
    return u_int_to_bin_str(address, 16) + '_' + u_int_to_bin_str(out_v,
                                                                  16) + '\n'


with open(OUTPUT_FILE, 'w') as f:
    for i in range(32768):
        f.write(build_line(i, 32767 - i))
