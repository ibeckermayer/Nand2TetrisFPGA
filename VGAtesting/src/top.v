// FPGA VGA Graphics Part 1: Top Module (static squares)
// (C)2017-2018 Will Green - Licensed under the MIT License
// Learn more at https://timetoexplore.net/blog/arty-fpga-vga-verilog-01

`default_nettype none

module top(input wire clk,           // board clock: 100 MHz on Arty/Basys3/Nexys
           input wire reset,         // reset button
           output wire vga_hs,       // horizontal sync output
           output wire vga_vs,       // vertical sync output
           output wire [3:0] vga_r,  // 4-bit VGA red output
           output wire [3:0] vga_g,  // 4-bit VGA green output
           output wire [3:0] vga_b,//); // 4-bit VGA blue output
        // Testbench only outputs
        output wire pixel_clk_tb,
        output wire o_active_tb,
        output wire [9:0] x_tb,
        output wire [8:0] y_tb,
        output wire [15:0] screen_addr_tb,
        output wire [15:0] screen_out_tb,
        output wire [3:0] bit_index_tb
     );
    
    wire [9:0] x;  // current pixel x position: 10-bit value: 0-1023
    wire [8:0] y;  // current pixel y position:  9-bit value: 0-511
    wire o_active; // high when vga640x480 is drawing and active pixel
    wire [15:0] screen_addr;
    
    // wire to hold the output of the ROM32K that's holding the screen data
    wire [15:0] screen_out;
    
    VGA320x240_Controller vga_ctrl
    (
    .clk(clk),
    .reset(reset),
    .screen_in(screen_out),
    .screen_addr(screen_addr),
    .vga_hs(vga_hs),
    .vga_vs(vga_vs),
    .vga_r(vga_r),
    .vga_b(vga_b),
    .vga_g(vga_g),
        .vga_active_o(o_active_tb),
        .bit_index_o(bit_index_tb),
        .pixel_clk_o(pixel_clk_tb),
        .x_o(x),
        .y_o(y)
    );
    
    RAMROM #(15, "/home/ibeckermayer/Nand2TetrisFPGA/VGAtesting/src/every_other.txt") screen
    (
    .clk(clk),
    .address(screen_addr),
    .out(screen_out)
    );
    
     // test bench assignments
     assign o_active_tb    = o_active;
     assign x_tb           = x;
     assign y_tb           = y;
     assign screen_addr_tb = screen_addr;
     assign screen_out_tb  = screen_out;
endmodule
