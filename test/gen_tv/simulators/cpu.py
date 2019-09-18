from simulators import BaseSimulator
from simulators.alu import ALUSimulator
from simulators.pc import PCSimulator
from typing import Tuple


class CPUSimulator(BaseSimulator):
  '''
  class to simulate the CPU
  '''

  def _i(self, i: int) -> int:
    '''
      Converts index from instruction variable in CPU.v to index for instruction variable
      in this file. Recall that the verilog file is reversed index, using this function
      makes it easier to track logic between files
      '''
    return 15 - i

  def __init__(self):
    self._A: int = 0
    self._D: int = 0
    self._ALU: ALUSimulator = ALUSimulator()
    self._PC: PCSimulator = PCSimulator()
    self._pc: int = 0  # PC output

  def simulate_step(self, inM: int, instruction: str,
                    reset: bool) -> Tuple[int, bool, int, int]:
    '''
    simulates a single time step in CPU operation
    returns (outM, writeM, addressM, pc)
    '''

    # _i(13)/_i(14) -- TODO: conspicuous that I can't find this being used in CPU.v.

    # ALU logic, which is run every step regardless of instruction type
    a = instruction[self._i(12)]
    A_or_inM = inM if (a == '1') else self._A
    c1 = instruction[self._i(11)]
    c2 = instruction[self._i(10)]
    c3 = instruction[self._i(9)]
    c4 = instruction[self._i(8)]
    c5 = instruction[self._i(7)]
    c6 = instruction[self._i(6)]
    c = c1 + c2 + c3 + c4 + c5 + c6
    alu_out, alu_zr, alu_ng = self._ALU.simulate_step(self._D, A_or_inM, c)

    # Set internal variables that will help determine the control logic of dependent on the dest and jump bits.
    is_A_instruction: bool = instruction[self._i(15)] == '0'
    is_C_instruction: bool = not is_A_instruction

    # Calculate destination bits, setting them all to zero (don't save anything) if this is an A instruction so that nothing in
    # memory is mistakenly overwritten by an A instruction that happens to say so if interpreted as a C instruction.
    d1 = instruction[self._i(5)] if is_C_instruction else '0'
    d2 = instruction[self._i(4)] if is_C_instruction else '0'
    d3 = instruction[self._i(3)] if is_C_instruction else '0'

    # Similar logic to above for the jump bits.
    j1 = instruction[self._i(2)] if is_C_instruction else '0'
    j2 = instruction[self._i(1)] if is_C_instruction else '0'
    j3 = instruction[self._i(0)] if is_C_instruction else '0'

    # Calculate if we should jump.
    # is j1 true and alu_out < 0?
    is_j1 = (True if alu_ng else False) if j1 else False
    # is j2 true and alu_out = 0?
    is_j2 = (True if alu_zr else False) if j2 else False
    # is j3 true and alu_out > 0?
    is_j3 = (True if
             (not (alu_ng) and not (alu_zr)) else False) if j3 else False
    jump = (is_j1 or is_j2 or is_j3)

    # Set all the outputs for this step.
    # addressM of this step is always whatever was in A in the last step.
    addressM: int = self._A
    # pc for this step is whatever the PC was set to or incremented to last step.
    pc: int = self._pc
    outM: int = alu_out
    writeM: bool = d3 is '1'

    # Now that we've calculated all the combinational logic, we are ready to calculate all the sequential logic.
    if (is_A_instruction):
      self._A = self.bin_str_to_int(instruction)
    elif (d1 == '1'):
      self._A = alu_out

    self._pc = self._PC.simulate_step(
        in_=self._A, reset=reset, load=jump, inc=1)

    if (d2 == '1'):
      self._D = alu_out

    return (outM, writeM, addressM, self._pc)

  def build_line(self, inM: int, instruction: str, reset: bool) -> str:
    '''
    format {inM[WIDTH], instruction[WIDTH], reset}_{outM[WIDTH], writeM, addressM[WIDTH], pc[WIDTH]}
    '''
    outM, writeM, addressM, pc = self.simulate_step(inM, instruction, reset)
    return \
      self.int_to_bin_str(inM, self.WIDTH) + \
      instruction + \
      self.int_to_bin_str(reset, 1) + \
      '_' + \
      self.int_to_bin_str(outM, self.WIDTH) + \
      self.int_to_bin_str(writeM, 1) + \
      self.int_to_bin_str(addressM, self.WIDTH) + \
      self.int_to_bin_str(pc, self.WIDTH) + '\n'
