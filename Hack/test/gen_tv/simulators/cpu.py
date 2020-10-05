from simulators import BaseSimulator
from simulators.alu import ALUSimulator
from simulators.pc import PCSimulator
from typing import Tuple, Optional


class CPUSimulator(BaseSimulator):
    """
    class to simulate the CPU
    """

    def _i(self, i: int) -> int:
        """
        Converts index from instruction variable in CPU.v to index for instruction variable
        in this file. Recall that the verilog file is reversed index, using this function
        makes it easier to track logic between files
        """
        return 15 - i

    def __init__(self):
        self._A: int = 0
        self._D: int = 0
        self._ALU: ALUSimulator = ALUSimulator()
        self._PC: PCSimulator = PCSimulator()
        self._pc: int = 0  # PC output

    def simulate_step(
        self, inM: int, instruction: str, reset: bool
    ) -> Tuple[int, bool, int, int, int, int, int]:
        """
        simulates a single time step in CPU operation
        returns (outM, writeM, addressM, pc)
        """

        # Set internal variables that will help determine the control logic of dependent on the dest and jump bits.
        is_A_instruction: bool = instruction[self._i(15)] == "0"
        is_C_instruction: bool = not is_A_instruction

        # ALU logic, which is run every step regardless of instruction type
        a = instruction[self._i(12)]
        A_or_inM = inM if (a == "1") else self._A
        c1 = instruction[self._i(11)]
        c2 = instruction[self._i(10)]
        c3 = instruction[self._i(9)]
        c4 = instruction[self._i(8)]
        c5 = instruction[self._i(7)]
        c6 = instruction[self._i(6)]
        c = c1 + c2 + c3 + c4 + c5 + c6
        alu_out, alu_zr, alu_ng = self._ALU.simulate_step(self._D, A_or_inM, c)

        # Calculate destination bits, setting them all to zero (don't save anything) if this is an A instruction so that nothing in
        # memory is mistakenly overwritten by an A instruction that happens to say so if interpreted as a C instruction.
        d1 = instruction[self._i(5)] if is_C_instruction else "0"
        d2 = instruction[self._i(4)] if is_C_instruction else "0"
        d3 = instruction[self._i(3)] if is_C_instruction else "0"

        # Similar logic to above for the jump bits.
        j1 = instruction[self._i(2)] if is_C_instruction else "0"
        j2 = instruction[self._i(1)] if is_C_instruction else "0"
        j3 = instruction[self._i(0)] if is_C_instruction else "0"

        # Calculate if we should jump.
        # is j1 true and alu_out < 0?
        # is_j1 = (True if alu_ng else False) if j1 == '1' else False
        is_j1 = (j1 == "1") and alu_ng
        # is j2 true and alu_out = 0?
        # is_j2 = (True if alu_zr else False) if j2 == '1' else False
        is_j2 = (j2 == "1") and alu_zr
        # is j3 true and alu_out > 0?
        # is_j3 = (True if (not (alu_ng) and not (alu_zr)) else False) if j3 == '1' else False
        is_j3 = (j3 == "1") and (not (alu_ng) and not (alu_zr))
        jump = is_j1 or is_j2 or is_j3

        self._pc = self._PC.simulate_step(in_=self._A, reset=reset, load=jump, inc=1)

        if d1 == "1":
            self._A = alu_out

        if d2 == "1":
            self._D = alu_out

        # If its an A-instruction, instruction is registered into A
        if is_A_instruction:
            self._A = self.bin_str_to_int(instruction)

        # Reset overrides everything else
        if reset:
            self._A = 0
            self._D = 0

        # Lastly calculate the outputs
        outM: int = alu_out
        writeM: bool = d3 is "1"
        addressM: int = self._A
        pc: int = self._pc

        return (outM, writeM, addressM, self._pc, alu_out, self._A, self._D)

    def build_line(self, inM: int, instruction: str, reset: bool) -> str:
        """
        format {inM[WIDTH], instruction[WIDTH], reset}_{outM[WIDTH], writeM, addressM[WIDTH], pc[WIDTH]_{alu_out[WIDTH], A[WIDTH], D[WIDTH]}}
        """
        outM, writeM, addressM, pc, alu_out, A, D = self.simulate_step(
            inM, instruction, reset
        )
        return (
            self.int_to_bin_str(inM, self.WIDTH)
            + instruction
            + self.int_to_bin_str(reset, 1)
            + "_"
            + self.int_to_bin_str(outM, self.WIDTH)
            + self.int_to_bin_str(writeM, 1)
            + self.int_to_bin_str(addressM, self.WIDTH)
            + self.int_to_bin_str(pc, self.WIDTH)
            + "_"
            + self.int_to_bin_str(alu_out, self.WIDTH)
            + self.int_to_bin_str(A, self.WIDTH)
            + self.int_to_bin_str(D, self.WIDTH)
            + "\n"
        )
