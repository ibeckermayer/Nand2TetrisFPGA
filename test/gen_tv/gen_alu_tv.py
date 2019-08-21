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
from util import ALUSimulator
from random import randint

OUTPUT_FILE = 'tvs/ALU.tv'  # expects to be run from directory above this
N = 1000  # number of times each function will be checked with random inputs for x and y

alusim = ALUSimulator()

with open(OUTPUT_FILE, 'w') as f:
    for func in alusim.funcs:
        for i in range(N):
            x = randint(-32768, 32767)
            y = randint(-32768, 32767)
            f.write(alusim.build_line(x, y, func))
