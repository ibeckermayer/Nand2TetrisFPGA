// A program that loads 0x1234 into RAM[0x0000] and RAM[0x1234], and then repeats.
//
// This program is designed for the Hack platform when it has the 7 segment display attached to the
// RAM's output. This program has the property that the A register will only ever have 0x0000 or 0x1234 
// loaded into it. The A register is wired to the CPU.addressM output, which in turn is wired to the 
// RAM.address input, which selects RAM.out which is attached to our 7 segment display. So if we put 
// 0x1234 into both RAM[0] and RAM[0x1234] and then repeat, our 7 segment display should always just show 
// us "1234". I'm using this as a sanity check to demonstrate to myself that something is actually working 
// how I think it is on hardware (up to now all testing has been software simulations).
@4660 // Load 0x1234 into the A register
D=A   // Copy 0x1234 into the D register
M=D   // RAM[0x1234] = 0x1234
@0    // Load 0x0000 into the A register
M=D   // RAM[0] = 0x1234
0;JMP // Jump back to instruction 0 (0x0000 still loaded in the A register) and repeat
