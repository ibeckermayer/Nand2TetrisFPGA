from numpy import binary_repr
from typing import Tuple

WIDTH = 16


class BaseSimulator(object):
    '''
    Base class for hardware device simulators
    Contains useful methods for dealing with integers and their binary string analogs
    '''

    @classmethod
    def int_to_bin_str(cls, val: int, width: int) -> str:
        '''
        converts an integer to it's twos complement representation of width=width
        val: number to be converted to binary
        width: bit-width of the output string (e.g. WIDTH-bits, 1-bit)
        '''
        retval = binary_repr(val, width)
        if len(retval) > width:
            return retval[-width:]
        else:
            return retval

    @classmethod
    def bin_str_to_int(cls, bin_str: str) -> int:
        '''
        inverse of int_to_bin_str
        NOTE: assumes bin_str is twos complement format if len(bin_str) > 1

        '''
        if (len(bin_str) == 1 or bin_str[0] == '0'):
            # single bit or positive value
            return int(bin_str, 2)
        elif (bin_str[1:] == '0' * (len(bin_str) - 1)):
            # edge case for lowest possible value
            return int(-((2**len(bin_str)) / 2))
        else:
            # flip all the bits
            tmp1 = cls.bitwise_not(bin_str)
            # add 1
            tmp2 = cls.int_to_bin_str(
                cls.bin_str_to_int(tmp1) + 1, len(bin_str))
            # interpret the result as a binary representation of the magnitude and add a negative sign
            tmp3 = cls.bin_str_to_int(tmp2)
            return -tmp3

    @classmethod
    def bitwise_not(cls, strval: str) -> str:
        notval = ''
        for dig in strval:
            if dig == '0':
                notval += '1'
            else:
                notval += '0'
        return notval

    @classmethod
    def bitwise_and(cls, strx: str, stry: str) -> str:
        andval = ''
        for i in range(WIDTH):
            if (strx[i] == '1' and stry[i] == '1'):
                andval += '1'
            else:
                andval += '0'
        return andval

    @classmethod
    def bitwise_or(cls, strx: str, stry: str) -> str:
        orval = ''
        for i in range(WIDTH):
            if (strx[i] == '1' or stry[i] == '1'):
                orval += '1'
            else:
                orval += '0'
        return orval


class RAMSimulator(BaseSimulator):
    '''
    class for simulating RAM
    '''

    def __init__(self):
        self.mem = {}  # dict(int, int) => {address: value}

    def simulate_step(self, address: int, in_: int, load: bool) -> int:
        '''
        takes in address, input value, and load boolean and returns the expected output
        '''
        if (load):
            self.mem[address] = in_
            # print(self.mem)
            # print(self.mem.get(address))

        return self.mem.get(address)

    def build_line(self, address: int, in_: int, load: bool) -> str:
        '''
        builds a line for the testvector file
        format: {address[WIDTH]}_{in[WIDTH]}_{load}_{out[WIDTH]}
        '''
        out = self.simulate_step(address, in_, load)
        if (out != None):
            out_s = self.int_to_bin_str(out, WIDTH)
        else:
            out_s = 'x' * WIDTH  # verilog representation of unknown values
        address_s = self.int_to_bin_str(address, WIDTH)
        in_s = self.int_to_bin_str(in_, WIDTH)
        load_s = self.int_to_bin_str(load, 1)

        return address_s + '_' + in_s + '_' + load_s + '_' + out_s + '\n'


class ROMSimulator(BaseSimulator):
    '''
    class for simulating ROM
    '''

    def __init__(self):
        self.mem = RAMSimulator()

    def load(self, address: int, in_: int) -> str:
        '''
        pseudo function of real ROM, necessary to load memory into ROM
        returns string of binary value loaded into ROM, to be used to build ROM input file
        '''
        return self.int_to_bin_str(self.mem.simulate_step(address, in_, True),
                                   WIDTH) + '\n'

    def simulate_step(self, address: int) -> int:
        '''
        takes in address, input value, and returns the expected output
        '''
        return self.mem.simulate_step(address, None, False)

    def build_line(self, address: int) -> str:
        '''
        builds a line for the testvector file
        format: {address[WIDTH]}_{out[WIDTH]}
        '''
        out = self.simulate_step(address)
        if (out != None):
            out_s = self.int_to_bin_str(out, WIDTH)
        else:
            out_s = 'x' * WIDTH  # verilog representation of unknown values
        address_s = self.int_to_bin_str(address, WIDTH)

        return address_s + '_' + out_s + '\n'


class ALUSimulator(BaseSimulator):
    '''
    class for simulating the ALU
    '''
    # strings for function values in the format (zx)(nx)(zy)(ny)(f)(no), taken from page 56
    funcs = [
        '101010',  # f(x,y) = 0
        '111111',  # f(x,y) = 1
        '111010',  # f(x,y) = -1
        '001100',  # f(x,y) = x
        '110000',  # f(x,y) = y
        '001101',  # f(x,y) = !x
        '110001',  # f(x,y) = !y
        '001111',  # f(x,y) = -x
        '110011',  # f(x,y) = -y
        '011111',  # f(x,y) = x+1
        '110111',  # f(x,y) = y+1
        '001110',  # f(x,y) = x-1
        '110010',  # f(x,y) = y-1
        '000010',  # f(x,y) = x+y
        '010011',  # f(x,y) = x-y
        '000111',  # f(x,y) = y-x
        '000000',  # f(x,y) = x&y
        '010101'  # f(x,y) = x|y
    ]

    def simulate_step(self, x: int, y: int,
                      func: str) -> Tuple[int, bool, bool]:
        '''
        returns (out, zr, nr)
        '''

        def _norm_zr_ng(out: int) -> Tuple[int, bool, bool]:
            '''
            helper to build the tuple for each function
            '''
            # account for overflow by normalizing to WIDTH-bit interpretation
            out_norm = self.bin_str_to_int(self.int_to_bin_str(out, WIDTH))
            return (out_norm, out_norm == 0, out_norm < 0)

        if func == '101010':  # f(x,y) = 0
            return _norm_zr_ng(0)
        elif func == '111111':  # f(x,y) = 1
            return _norm_zr_ng(1)
        elif func == '111010':  # f(x,y) = -1
            return _norm_zr_ng(-1)
        elif func == '001100':  # f(x,y) = x
            return _norm_zr_ng(x)
        elif func == '110000':  # f(x,y) = y
            return _norm_zr_ng(y)
        elif func == '001101':  # f(x,y) = !x
            return _norm_zr_ng(
                self.bin_str_to_int(
                    self.bitwise_not(self.int_to_bin_str(x, WIDTH))))
        elif func == '110001':  # f(x,y) = !y
            return _norm_zr_ng(
                self.bin_str_to_int(
                    self.bitwise_not(self.int_to_bin_str(y, WIDTH))))
        elif func == '001111':  # f(x,y) = -x
            return _norm_zr_ng(-x)
        elif func == '110011':  # f(x,y) = -y
            return _norm_zr_ng(-y)
        elif func == '011111':  # f(x,y) = x+1
            return _norm_zr_ng(x + 1)
        elif func == '110111':  # f(x,y) = y+1
            return _norm_zr_ng(y + 1)
        elif func == '001110':  # f(x,y) = x-1
            return _norm_zr_ng(x - 1)
        elif func == '110010':  # f(x,y) = y-1
            return _norm_zr_ng(y - 1)
        elif func == '000010':  # f(x,y) = x+y
            return _norm_zr_ng(x + y)
        elif func == '010011':  # f(x,y) = x-y
            return _norm_zr_ng(x - y)
        elif func == '000111':  # f(x,y) = y-x
            return _norm_zr_ng(y - x)
        elif func == '000000':  # f(x,y) = x&y
            return _norm_zr_ng(
                self.bin_str_to_int(
                    self.bitwise_and(self.int_to_bin_str(x, WIDTH),
                                     self.int_to_bin_str(y, WIDTH))))
        elif func == '010101':  # f(x,y) = x|y
            return _norm_zr_ng(
                self.bin_str_to_int(
                    self.bitwise_or(self.int_to_bin_str(x, WIDTH),
                                    self.int_to_bin_str(y, WIDTH))))
        else:
            raise Exception('Unkown function')

    def build_line(self, x: int, y: int, func: str) -> str:
        '''
        builds a line for the testvector file
        format {x[WIDTH]}_{y[WIDTH]}_{zx, nx, zy, ny, f, no}_{out[WIDTH], zr, ng}
        '''
        out, zr, ng = self.simulate_step(x, y, func)
        out_zr_ng_s = self.int_to_bin_str(out, WIDTH) + self.int_to_bin_str(
            zr, 1) + self.int_to_bin_str(ng, 1)
        x_s = self.int_to_bin_str(x, WIDTH)
        y_s = self.int_to_bin_str(y, WIDTH)
        return x_s + '_' + y_s + '_' + func + '_' + out_zr_ng_s + '\n'
