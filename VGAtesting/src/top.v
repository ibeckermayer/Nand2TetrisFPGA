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
           output wire [3:0] vga_b); // 4-bit VGA blue output
    //    Testbench only outputs
    //    output wire pixel_clk_tb,
    //    o_active_tb,
    //    output wire [9:0] x_tb,
    //    output wire [8:0] y_tb,
    //    output wire [15:0] screen_addr_tb,
    //    screen_out_tb,
    //    output wire [3:0] bit_index_tb
    //);
    
    localparam H_DISPLAY       = 640; // horizontal display area
    localparam V_DISPLAY       = 480;   // vertical display area
    localparam MAX_ADDRESS_REG = ((H_DISPLAY * V_DISPLAY) / 16) - 1; // 16 is bit-width of current_word
    
    // generate a 25 MHz pixel clock
    reg [15:0] cnt;
    reg pixel_clk;
    always @(posedge clk)
    begin
        if (reset)
        begin
            cnt       <= 0;
            pixel_clk <= 0;
        end
        else
            {pixel_clk, cnt} <= cnt + 16'h4000;  // divide by 4: (2^16)/4 = 0x4000
    end
    
    
    wire [9:0] x;  // current pixel x position: 10-bit value: 0-1023
    wire [8:0] y;  // current pixel y position:  9-bit value: 0-511
    wire o_active; // high when vga640x480 is drawing and active pixel
    
    vga640x480 display (
    .clk(clk),
    .pixel_clk(pixel_clk),
    .reset(reset),
    .o_hs(vga_hs),
    .o_vs(vga_vs),
    .o_active(o_active),
    .o_x(x),
    .o_y(y)
    );
    
    // registers to track the screen address and bit of the corresponding word to read out
    reg [15:0] screen_addr;
    reg [3:0] bit_index;
    
    always @(posedge clk)
    begin
        if (reset)
        begin
            screen_addr <= 0;
            bit_index   <= 0;
        end
            if (pixel_clk && o_active)
            begin
                bit_index <= bit_index + 1; // bit index is always counted on active pixels
                if (bit_index == 4'b1111) // active pixel tick before the index rolls over to 0
                begin
                    if (screen_addr == MAX_ADDRESS_REG)
                        screen_addr <= 0; // if this is the maximum reg for the screen, reset to zero
                    else
                        screen_addr <= screen_addr + 1; // else continue on to the next address
                end
            end
    end
    
    // wire to hold the output of the ROM32K that's holding the screen data
    wire [15:0] screen_out;
    
    RAMROM #(16, "/home/ibeckermayer/Nand2TetrisFPGA/VGAtesting/src/every_other.txt") screen
    (
    .clk(clk),
    .address(screen_addr),
    .out(screen_out)
    );
    
    assign {vga_r, vga_g, vga_b} = o_active ? (screen_out[bit_index] == 1 ? 12'b111100000000: 12'b000011110000) : 12'b000000000000;
    
    // // test bench assignments
    // assign pixel_clk_tb   = pixel_clk;
    // assign o_active_tb    = o_active;
    // assign x_tb           = x;
    // assign y_tb           = y;
    // assign screen_addr_tb = screen_addr;
    // assign screen_out_tb  = screen_out;
    // assign bit_index_tb   = bit_index;
endmodule
