'''
script for generating the testvector file for use in PC_tb.v

Algorithm:
- start with a reset, expect a 0
- count up to 1000 to check increment
- reset again to check reset at arbitrary count
- count up to 1000 again
- load 65000 to test loading from arbitrary count
- count up to 65535
- count up to 1000 again to check rollover
'''
from util import PCSimulator

OUTPUT_FILE = 'tvs/PC.tv'  # expects to be run from directory above this

pcsim = PCSimulator()

# initial values
in_v = 1
reset = 1  # start with a reset
load = 1  # start with load 1, see that reset has precedence
inc = 1  # start with inc 1, see that reset has precedence

with open(OUTPUT_FILE, 'w') as f:
  f.write(pcsim.build_line(in_v, reset, load, inc))

  # count up to 1000 to check increment
  reset = 0
  load = 0
  for i in range(1000):
    f.write(pcsim.build_line(in_v, reset, load, inc))

  # reset again to check reset at arbitrary count
  reset = 1
  f.write(pcsim.build_line(in_v, reset, load, inc))

  # count up to 1000 again
  reset = 0
  load = 0
  for i in range(1000):
    f.write(pcsim.build_line(in_v, reset, load, inc))

  # load 65000 to test loading from arbitrary count
  load = 1
  in_v = 65000
  f.write(pcsim.build_line(in_v, reset, load, inc))

  # count up to 65535
  # count up to 1000 again to check rollover
  load = 0
  for i in range(1001 + 535):
    f.write(pcsim.build_line(in_v, reset, load, inc))
