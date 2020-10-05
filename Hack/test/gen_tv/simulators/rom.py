from simulators import BaseSimulator
from simulators.ram import RAMSimulator


class ROMSimulator(BaseSimulator):
    """
    class for simulating ROM
    """

    def __init__(self):
        self.mem = RAMSimulator()

    def load(self, address: int, in_: int) -> str:
        """
        pseudo function of real ROM, necessary to load memory into ROM
        returns string of binary value loaded into ROM, to be used to build ROM input file
        """
        return (
            self.int_to_bin_str(self.mem.simulate_step(address, in_, True), self.WIDTH)
            + "\n"
        )

    def simulate_step(self, address: int) -> int:
        """
        takes in address, input value, and returns the expected output
        """
        return self.mem.simulate_step(address, None, False)

    def build_line(self, address: int) -> str:
        """
        builds a line for the testvector file
        format: {address[WIDTH]}_{out[WIDTH]}
        """
        out = self.simulate_step(address)
        if out != None:
            out_s = self.int_to_bin_str(out, self.WIDTH)
        else:
            out_s = "x" * self.WIDTH  # verilog representation of unknown values
        address_s = self.int_to_bin_str(address, self.WIDTH)

        return address_s + "_" + out_s + "\n"
