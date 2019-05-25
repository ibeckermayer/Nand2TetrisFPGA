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
        width: bit-width of the output string (e.g. 16-bits, 1-bit)
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
        for i in range(16):
            if (strx[i] == '1' and stry[i] == '1'):
                andval += '1'
            else:
                andval += '0'
        return andval

    @classmethod
    def bitwise_or(cls, strx: str, stry: str) -> str:
        orval = ''
        for i in range(16):
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

        return self.mem.get(address)

    def build_line(self, address: int, in_: int, load: bool) -> str:
        '''
        builds a line for the testvector file
        '''
        address_s = self.int_to_bin_str(address, WIDTH)
        in_s = self.int_to_bin_str(in_, WIDTH)
        load_s = self.int_to_bin_str(load, WIDTH)
        # build output
        out = self.simulate_step(address, in_, load)
        if (out):
            out_s = self.int_to_bin_str(out, WIDTH)
        else:
            out_s = 'x' * 16  # verilog representation of unknown values

        return address_s + in_s + load_s + '_' + out_s + '\n'


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
            # account for overflow by normalizing to 16-bit interpretation
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
        return (x[16])(y[16])(zx)(nx)(zy)(ny)(f)(no)_(out[16])(zr)(ng)
        '''
        out, zr, ng = self.simulate_step(x, y, func)
        return self.int_to_bin_str(x, WIDTH) + self.int_to_bin_str(
            y, WIDTH) + func + '_' + self.int_to_bin_str(
                out, WIDTH) + self.int_to_bin_str(zr, 1) + self.int_to_bin_str(
                    ng, 1) + '\n'
