from simulators import BaseSimulator
from typing import Tuple


class CPUSimulator(BaseSimulator):
  '''
  class to simulate the CPU
  '''

  def _i(cls, i: int) -> int:
    '''
      Converts index from instruction variable in CPU.v to index for instruction variable
      in this file. Recall that the verilog file is reversed index, using this function
      makes it easier to track logic between files
      '''
    return 15 - i

  def __init__(self):
    # TODO: ought to be set to a better default value
    self._A: int = None
    self._D: int = None  # TODO: ought to be set to a better default value
    self._ALU: ALUSimulator = ALUSimulator()
    self._PC: PCSimulator = PCSimulator()

  def simulate_step(self, inM: int, instruction: str,
                    reset: bool) -> Tuple[int, bool, int, int]:
    '''
    simulates a single time step in CPU operation
    returns (outM, writeM, addressM, pc)

    TODO: Idea for an approach: first push everything through, then set everything clocked
    '''
    # _i(13)/_i(14) -- TODO: conspicuous that I can't find this being used in CPU.v.
    # This should be investigated when you have access to the book again

    # The first thing we should do is decide whether this is an A instruction or a C instruction
    # since that distinction largely determines the behavior of the CPU
    is_A = instruction[self._i(15)] == '0'
    is_C = not is_A

    # ## TODO: left off here

    # Let's do the simplest thing first: if this is an instruction, the value is loaded into the A register
    if instruction[self._i(15)] == '0':
      self._A = self.bin_str_to_int(instruction)

    # Next, we can simulate a step on all of the combinational logic

    ############################################################################################
    # Start by fleshing all the combinational logic out, meaning (???):
    # If A has an input B, set B first.
    a = instruction[self._i(12)]
    c1 = instruction[self._i(11)]
    c2 = instruction[self._i(10)]
    c3 = instruction[self._i(9)]
    c4 = instruction[self._i(8)]
    c5 = instruction[self._i(7)]
    c6 = instruction[self._i(6)]
    c = c1 + c2 + c3 + c4 + c5 + c6
    d1 = instruction[self._i(5)] if instruction[self._i(15)] == '1' else '0'
    d2 = instruction[self._i(4)] if instruction[self._i(15)] == '1' else '0'
    d3 = instruction[self._i(3)] if instruction[self._i(15)] == '1' else '0'
    j1 = instruction[self._i(2)]
    j2 = instruction[self._i(1)]
    j3 = instruction[self._i(0)]

    alu_out, alu_zr, alu_ng = self._ALU.simulate_step(self._D, A_or_inM, c)
    is_j1 = (True if alu_ng else
             False) if j1 else False  # is j1 true and alu_out < 0?
    is_j2 = (True if alu_zr else
             False) if j2 else False  # is j2 true and alu_out = 0?
    is_j3 = (True if (not (alu_ng) and not (alu_zr)) else
             False) if j3 else False  # is j3 true and alu_out > 0?
    jump = (is_j1 or is_j2 or is_j3)
    outM = alu_out
    writeM = d3 == '1'

    # Now flesh out all the sequential logic, meaning (???):
    # If A has an input B, set B first.
    # program counter gets incremented or set (from jump instruction)
    self._A = self.bin_str_to_int(instruction) if instruction[self._i(
        15)] == '0' else alu_out if d1 == '1' else self._A
    self._pc = self._PC.simulate_step(
        in_=self._A, reset=reset, load=jump, inc=1)
    addressM = self._A
    A_or_inM = inM if (a == '1') else self._A
    self._D = alu_out if d2 == '1' else self._D

    return (outM, writeM, addressM, self._pc)
