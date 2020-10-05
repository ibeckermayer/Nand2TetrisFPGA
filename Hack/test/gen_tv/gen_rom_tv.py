"""
script for generating the testvector file for use in ROM_tb.v

ROM_MEMORY_INPUT Algorithm:
- generate input for ROM that begins with 32767 in address 0, and ends with 0 in address 32767

Format will be:
(address)_(out_v)

Algorithm:
- go through each address in ram and expect 32767 - address
"""
from simulators.rom import ROMSimulator

# expects to be run from directory above this
ROM_MEMORY_INPUT = "tvs/ROM_input.tv"  # ROM32K must be loaded from a file
OUTPUT_FILE = "tvs/ROM.tv"  # file for the testbench

romsim = ROMSimulator()

with open(ROM_MEMORY_INPUT, "w") as f:
    for i in range(32768):
        f.write(romsim.load(i, 32767 - i))

with open(OUTPUT_FILE, "w") as f:
    for i in range(32768):
        f.write(romsim.build_line(i))
