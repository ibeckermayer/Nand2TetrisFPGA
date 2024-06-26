# Nand2TetrisFPGA

This project is an FPGA implementation of the Hack computer architecture described in the book [The Elements of Computing Systems](https://www.nand2tetris.org/), By Noam Nisan and Shimon Schocken.

Hardware components are described in Verilog and can be found in the `Hack` directory.

The assembler is written in C and is located in the `Assembler` directory.

## Vivado

Consistent version control with Vivado has wound up being a gigantic PITA to maintain transferability between different machines. The basic idea (which I've been forced to repeat manually several times) is to create a Vivado project, select the Basys 3 board, and then add all of the sources in `Hack/src` to the project and generate the bitstream.

See [this line](https://github.com/ibeckermayer/Nand2TetrisFPGA/blob/349bea802aa11eefb987400cac9005b7725ff33a/Hack/src/Hack.v#L21) in the `Hack.v` file for the program that Vivado will try to load onto the FPGA. You will have
to ensure that this file exists and that Vivado can find it.

## Progress Update (Oct 6, 2020)

While the photo deserves its title of `jank_ass_vga.jpg`, its also the result of a lot of effort. What's happening in here is that a program (compiled by the Assembler) is running on the Hack computer architecture, which I've synthesized onto [Basys3](https://store.digilentinc.com/basys-3-artix-7-fpga-trainer-board-recommended-for-introductory-users/) FPGA development board. The program simply writes the value of `0b0010101010101010` (value chosen for idiosyncratic technical reasons I will not go into here) in RAM registers 16384 to 21184. Those registers are read out by a circuit (`Hack/src/VGA/VGA320x240_Controller.v`) that converts the individual bits stored in those registers into pixels displayed on the screen. The counting is a little bit tricky, because I needed to halve the ordinary 640x480 VGA output (which would take up 19,200 words of the 32k RAM) down to 320x240 (4,800 words).

![alt text](https://raw.githubusercontent.com/ibeckermayer/Nand2TetrisFPGA/master/Misc/jank_ass_vga.png)
