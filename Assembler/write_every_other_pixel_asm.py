# 0b101010101010101 = 21845
EVERY_OTHER = 21845
FIRST_SCREEN_REG = 16384
with open("write_every_other_pixel.asm", "w") as f:
    f.write(f"@{EVERY_OTHER} // Load 0b101010101010101 into the A register\n")
    f.write("D=A //Copy 0b101010101010101 into the D register\n")
    for i in range(FIRST_SCREEN_REG, int(FIRST_SCREEN_REG + (320 * 240 / 16))):
        f.write(f"@{i} // Load {i} into the A register\n")
        f.write(f"M=D // RAM[{i}] = 0b101010101010101\n")
    f.write("(END)\n")
    f.write("@END\n")
    f.write("0;JMP // Infinite loop\n")