def u_int_to_bin_str(num: int, width: int) -> str:
    '''
    converts a positive integer to an unsigned binary string representation
    num: number to be converted to binary
    width: bit-width of the output string (e.g. 16-bits, 1-bit)
    '''
    binary_string = '{:b}'.format(num)
    full_width = len(binary_string)
    if (width <= full_width):
        return binary_string[-width:]
    else:
        return ('0' * (width - full_width)) + binary_string
