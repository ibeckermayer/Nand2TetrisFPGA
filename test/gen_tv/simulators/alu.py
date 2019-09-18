from simulators import BaseSimulator
from typing import Tuple


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

  def simulate_step(self, x: int, y: int, func: str) -> Tuple[int, bool, bool]:
    '''
        returns (out, zr, nr)
        '''

    def _norm_zr_ng(out: int) -> Tuple[int, bool, bool]:
      '''
            helper to build the tuple for each function
            '''
      # account for overflow by normalizing to WIDTH-bit interpretation
      out_norm = self.bin_str_to_int(self.int_to_bin_str(out, self.WIDTH))
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
              self.bitwise_not(self.int_to_bin_str(x, self.WIDTH))))
    elif func == '110001':  # f(x,y) = !y
      return _norm_zr_ng(
          self.bin_str_to_int(
              self.bitwise_not(self.int_to_bin_str(y, self.WIDTH))))
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
              self.bitwise_and(
                  self.int_to_bin_str(x, self.WIDTH),
                  self.int_to_bin_str(y, self.WIDTH))))
    elif func == '010101':  # f(x,y) = x|y
      return _norm_zr_ng(
          self.bin_str_to_int(
              self.bitwise_or(
                  self.int_to_bin_str(x, self.WIDTH),
                  self.int_to_bin_str(y, self.WIDTH))))
    else:
      raise Exception('Unkown function')

  def build_line(self, x: int, y: int, func: str) -> str:
    '''
        builds a line for the testvector file
        format {x[WIDTH]}_{y[WIDTH]}_{zx, nx, zy, ny, f, no}_{out[WIDTH], zr, ng}
        '''
    out, zr, ng = self.simulate_step(x, y, func)
    out_zr_ng_s = self.int_to_bin_str(out, self.WIDTH) + self.int_to_bin_str(
        zr, 1) + self.int_to_bin_str(ng, 1)
    x_s = self.int_to_bin_str(x, self.WIDTH)
    y_s = self.int_to_bin_str(y, self.WIDTH)
    return x_s + '_' + y_s + '_' + func + '_' + out_zr_ng_s + '\n'
