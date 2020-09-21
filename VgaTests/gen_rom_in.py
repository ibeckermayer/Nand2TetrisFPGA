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

with open('all_white_rom.txt', 'w') as f:
    # bin_value = False # first line will be opposite bin value because 0 % 16 == 0
    bin_value = True
    for i in range(480):
        # if i%16 == 0:
        #     bin_value = not(bin_value)
        write_row(f, bin_value)


