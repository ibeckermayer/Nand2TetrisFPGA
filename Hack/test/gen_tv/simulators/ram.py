from simulators import BaseSimulator


class RAMSimulator(BaseSimulator):
    """
    class for simulating RAM
    """

    def __init__(self):
        self.mem = {}  # dict(int, int) => {address: value}

    def simulate_step(self, address: int, in_: int, load: bool) -> int:
        """
        takes in address, input value, and load boolean and returns the expected output
        """
        if load:
            self.mem[address] = in_

        return self.mem.get(address)

    def build_line(self, address: int, in_: int, load: bool) -> str:
        """
        builds a line for the testvector file
        format: {address[WIDTH]}_{in[WIDTH]}_{load}_{out[WIDTH]}
        """
        out = self.simulate_step(address, in_, load)
        if out != None:
            out_s = self.int_to_bin_str(out, self.WIDTH)
        else:
            out_s = "x" * self.WIDTH  # verilog representation of unknown values
        address_s = self.int_to_bin_str(address, self.WIDTH)
        in_s = self.int_to_bin_str(in_, self.WIDTH)
        load_s = self.int_to_bin_str(load, 1)

        return address_s + "_" + in_s + "_" + load_s + "_" + out_s + "\n"
