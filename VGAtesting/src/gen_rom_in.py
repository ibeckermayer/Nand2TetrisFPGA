# generate ram array that should be translated into alternating black and white stripes
# 640 x 480 screen, 16-bit width rom array
# 16 rows black, then 16 rows white (30 horizontal lines in total)
def write_row(file, bin_value):
    '''
    Writes a row of pixels as a single binary value
    '''
    for i in range(int(640/16)):
        if bin_value:
            file.write('1' * 16 + '\n')
        else:
            file.write('0' * 16 + '\n')

def write_every_other_bit_switched_row(file, start_val):
    '''
    Writes a row where every other bit is flipped, starting with start
    '''
    this_val = start_val
    for i in range(int(640/16)):
        for i in range(16):
            if this_val:
                file.write('1')
            else:
                file.write('0')
            this_val = not(this_val)
        file.write('\n')
                


def write_half_and_half_txt():
    with open('half_and_half.txt', 'w') as f:
        bin_value = False # first line will be opposite bin value because 0 % 16 == 0
        for i in range(480):
            if i < 480/2:
                write_row(f, True)
            else:
                write_row(f, False)


def write_every_other_txt():
    with open('every_other.txt', 'w') as f:
        start_val = True
        for i in range(480):
            write_every_other_bit_switched_row(f, start_val)
            start_val = not(start_val)

write_half_and_half_txt()



        






    


