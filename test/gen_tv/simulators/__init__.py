from numpy import binary_repr


class BaseSimulator(object):
  '''
    Base class for hardware device simulators
    Contains useful methods for dealing with integers and their binary string analogs
    '''
  WIDTH = 16

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
      tmp2 = cls.int_to_bin_str(cls.bin_str_to_int(tmp1) + 1, len(bin_str))
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
    for i in range(cls.WIDTH):
      if (strx[i] == '1' and stry[i] == '1'):
        andval += '1'
      else:
        andval += '0'
    return andval

  @classmethod
  def bitwise_or(cls, strx: str, stry: str) -> str:
    orval = ''
    for i in range(cls.WIDTH):
      if (strx[i] == '1' or stry[i] == '1'):
        orval += '1'
      else:
        orval += '0'
    return orval