from numpy import binary_repr


def int_to_bin_str(val: int, width: int) -> str:
    '''
    converts an integer to it's twos complement representation of width=width
    val: number to be converted to binary
    width: bit-width of the output string (e.g. 16-bits, 1-bit)
    '''
    retval = binary_repr(val, width)
    if len(retval) > width:
        return retval[1:]
    else:
        return retval
