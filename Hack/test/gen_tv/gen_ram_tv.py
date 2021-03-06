"""
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
"""
from simulators.ram import RAMSimulator

OUTPUT_FILE = "tvs/RAM.tv"  # expects to be run from directory above this

ramsim = RAMSimulator()

with open(OUTPUT_FILE, "w") as f:
    # address and in_v and out_v = 0..32767, load = 1: load each address with its value (test load)
    for i in range(32768):
        f.write(ramsim.build_line(i, i, True))

    # go through again with load, in_v = 0 and read each value (test read only)
    for i in range(32768):
        f.write(ramsim.build_line(i, 0, False))