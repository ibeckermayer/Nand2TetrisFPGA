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
    //    Testbench only outputs
       output wire pixel_clk_tb,
       o_active_tb,
       output wire [9:0] x_tb,
       output wire [8:0] y_tb,
       output wire [15:0] screen_addr_tb,
       screen_out_tb,
       output wire [3:0] bit_index_tb
    );
    
    localparam H_DISPLAY                = 640; // horizontal display area
    localparam V_DISPLAY                = 480;   // vertical display area
    localparam MAX_ADDRESS_REG          = (((H_DISPLAY / 2) * (V_DISPLAY / 2)) / 16) - 1; // 16 is bit-width of current_word
    localparam ADDRESSES_PER_SCREEN_ROW = H_DISPLAY / 2 / 16;
    
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
    reg add_to_bit_index; // determines whether we should add to the bit_index on the pixel tick
    reg [14:0] row_beginning_addr;
    reg [14:0] row_end_addr;
    reg reset_row;
    
    always @(posedge clk)
    begin
        if (reset)
        begin
            screen_addr        <= 0;
            bit_index          <= 0;
            add_to_bit_index   <= 0;
            row_beginning_addr <= 0;
            row_end_addr       <= ADDRESSES_PER_SCREEN_ROW - 1;
            reset_row          <= 1;
        end
        else
        begin
            if (pixel_clk && o_active)
            begin
                add_to_bit_index <= !add_to_bit_index; // Always flip our add_to_bit_index tracker
                
                if (add_to_bit_index == 1)
                begin
                    // count the bit index up on every other active pixel
                    // each bit in memory will be drawn on two pixels next to each other, to halve our resoution from 640 to 320
                    bit_index <= bit_index + 1;
                end
                
                if (bit_index == 4'b1111 && add_to_bit_index == 1) // active pixel tick before the index rolls over to 0, check if we're at the end of a row
                begin
                    // screen_addr state machine
                    if (screen_addr == row_end_addr) // we drew the final pixel of a 640 pixel wide row
                    begin
                        reset_row   <= !reset_row;
                        // if we've only drawn this row once, draw it again in order to halve our 480 resolution to 240
                        if (reset_row)
                        begin
                            screen_addr <= row_beginning_addr;
                        end
                        else
                        begin
                            if (screen_addr == MAX_ADDRESS_REG)
                            begin
                                // if this is the last register and we've already drawn it twice (reset_row == 0), reset the machine to pixel 0
                                screen_addr        <= 0;
                                row_beginning_addr <= 0;
                                row_end_addr       <= ADDRESSES_PER_SCREEN_ROW - 1;
                                
                            end
                            else
                            begin
                                // else add to the screen_addr, and update the row_beginning_addr and row_end_addr
                                screen_addr        <= screen_addr + 1;
                                row_beginning_addr <= row_beginning_addr + ADDRESSES_PER_SCREEN_ROW;
                                row_end_addr       <= row_end_addr + ADDRESSES_PER_SCREEN_ROW;
                            end
                        end
                    end
                    else
                    begin
                        screen_addr <= screen_addr + 1; // else continue on to the next address
                    end
                end
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
    
    // test bench assignments
    assign pixel_clk_tb   = pixel_clk;
    assign o_active_tb    = o_active;
    assign x_tb           = x;
    assign y_tb           = y;
    assign screen_addr_tb = screen_addr;
    assign screen_out_tb  = screen_out;
    assign bit_index_tb   = bit_index;
endmodule
