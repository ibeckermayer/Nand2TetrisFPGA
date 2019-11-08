'''
script for generating the testvector file for use in ALU_tb.v

Format will be:
{inM[WIDTH], instruction[WIDTH], reset}_{outM[WIDTH], writeM, addressM[WIDTH], pc[WIDTH]}
NOTE that clk and reset should be set internally in the testbench

Algorithm:
  - Follow the alu approach of blasting this with random numbers and seeing if both the python
    model and the verilog model agree. This begs the question what's testing what? We're also
    now running in to an unforseen complication wherein we can only load A-instructions that also
    happen to have values that have a c1-c6 that would be a valid ALU function. This is because
    of how the ALU logic is hardcoded in the python simulation. The verilog in ALU.v is a high
    enough level of abstraction that I'm not sure what it does on an invalid instruction. It would
    be interesting to check this model out in that Xilinx tool that allows you to see the block/wiring
    diagram and/or a signal-probe simulation, since that might give you more insight into how the high 
    level verilog is interpreted.
    
    This all makes me consider whether this is moment where formal verification might come into
    play -- I've heard that it's used effectively for digital design. In some sense it would also
    be a test of whether one's formal verification problem was set up correctly. I imagine that it
    might be better because formal verification might make it easier to specify and test all the possible
    edge cases. Perhaps I'm being lazy, but trying to do that using my python simulations for a machine as 
    complex as the CPU is seems extremely daunting.
'''
from simulators.cpu import CPUSimulator
from simulators.alu import ALUSimulator, UnkownALUFunction
from random import randint

OUTPUT_FILE = 'tvs/CPU.tv'  # expects to be run from directory above this
N = 10000
i = 0

# Flag in case you only want to generate A instructions to help
# narrow down the debugging process
GENERATE_ONLY_A_INSTRUCTIONS = False

cpusim = CPUSimulator()


def gen_random_instruction() -> str:
  '''
  Function for generating a random instruction since there are some rules instructions should play by.
  '''
  possible_instruction = cpusim.int_to_bin_str(
      randint(-32768, 32767), cpusim.WIDTH)

  if GENERATE_ONLY_A_INSTRUCTIONS:
    return '0' + possible_instruction[1:]

  if possible_instruction[0] == '0':
    # if this is an A instruction go ahead and return it right away
    return possible_instruction
  else:
    # else this is a C instruction, possible_instruction[1:3] == '1'
    # per the specification 4.2.3 The C-Instruction
    possible_instruction = possible_instruction[
        0] + '11' + possible_instruction[3:]
    # a[4:10] must be a valid function for the alu
    possible_instruction = possible_instruction[0:4] + \
      ALUSimulator.funcs[randint(0, len(ALUSimulator.funcs) - 1)] + \
      possible_instruction[10:]
    return possible_instruction


with open(OUTPUT_FILE, 'w') as f:
  while i < N:
    inM = randint(-32768, 32767)
    instruction = gen_random_instruction()
    if (i == 0):
      reset = True
    else:
      # TODO: make reset = bool(randint(0, 1)), this
      # current setup is just easier for debugging purposes
      reset = False
    try:
      f.write(cpusim.build_line(inM, instruction, reset))
      i += 1
    except UnkownALUFunction as e:
      continue
