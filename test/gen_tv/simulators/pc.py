from simulators import BaseSimulator


class PCSimulator(BaseSimulator):
  '''
    class for simulation program counter (PC)
    '''

  def __init__(self):
    self.out = None  # int

  def simulate_step(self, in_: int, reset: int, load: int, inc: int) -> int:
    if (reset):
      self.out = 0
    elif (load):
      self.out = in_
    elif (inc):
      self.out = self.out + 1
    # else: do nothing
    return self.out

  def build_line(self, in_: int, reset: int, load: int, inc: int) -> str:
    '''
        format {in[16]}_{reset}_{load}_{inc}_{out[16]}
        '''
    out = self.simulate_step(in_, reset, load, inc)
    out_s = self.int_to_bin_str(out, self.WIDTH)
    in_s = self.int_to_bin_str(in_, self.WIDTH)
    reset_s = self.int_to_bin_str(reset, 1)
    load_s = self.int_to_bin_str(load, 1)
    inc_s = self.int_to_bin_str(inc, 1)

    return in_s + '_' + reset_s + '_' + load_s + '_' + inc_s + '_' + out_s + '\n'
