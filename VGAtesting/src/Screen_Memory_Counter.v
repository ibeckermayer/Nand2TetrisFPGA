// Module for counting screen memory addresses and individual bits in order to synchronize the screen memory's output with the
// 640x480 VGA signal. Our system is saving space by halving the 640x480 screen down to 320x240, meaning that the counter
// should tick up every other VGA pixel (halving the horizontal resolution), and should reset to the first bit of the previous
// VGA line every other line (to halve the vertical resolution).
module Screen_Memory_Counter(input wire clk,
                             input wire reset,
                             input wire pixel_clk,
                             input wire vga_active,
                             output wire [15:0] screen_addr,
                             output wire [3:0] bit_index);
    
    localparam H_DISPLAY                = 640 / 2; // horizontal display area
    localparam V_DISPLAY                = 480 / 2;   // vertical display area
    localparam MAX_ADDRESS_REG          = ((H_DISPLAY * V_DISPLAY) / 16) - 1; // 16 is bit-width of current_word
    localparam ADDRESSES_PER_SCREEN_ROW = H_DISPLAY / 16;
    localparam FIRST_SCREEN_REG_ADDR    = 0; // Address of the first screen register
    
    // registers to track the screen address and bit of the corresponding word to read out
    reg [15:0] screen_addr_internal;
    reg [3:0] bit_index_internal;
    reg add_to_bit_index_internal; // determines whether we should add to the bit_index_internal on the pixel tick
    reg [14:0] row_beginning_addr; // first address corresponding to the current vga row
    reg [14:0] row_end_addr; // last address corresponding to the current vga row
    reg reset_row;
    
    always @(posedge clk)
    begin
        if (reset)
        begin
            screen_addr_internal      <= FIRST_SCREEN_REG_ADDR;
            bit_index_internal        <= 0;
            add_to_bit_index_internal <= 0;
            row_beginning_addr        <= FIRST_SCREEN_REG_ADDR;
            row_end_addr              <= FIRST_SCREEN_REG_ADDR + ADDRESSES_PER_SCREEN_ROW - 1;
            reset_row                 <= 1;
        end
        else
        begin
            if (pixel_clk && vga_active)
            begin
                add_to_bit_index_internal <= !add_to_bit_index_internal; // Always flip our add_to_bit_index_internal tracker
                
                if (add_to_bit_index_internal == 1)
                begin
                    // count the bit index up on every other active pixel
                    // each bit in memory will be drawn on two pixels next to each other, to halve our resoution from 640 to 320
                    bit_index_internal <= bit_index_internal + 1;
                end
                
                if (bit_index_internal == 4'b1111 && add_to_bit_index_internal == 1) // active pixel tick before the index rolls over to 0, check if we're at the end of a row
                begin
                    // screen_addr_internal state machine
                    if (screen_addr_internal == row_end_addr) // we drew the final pixel of a 640 pixel wide row
                    begin
                        reset_row <= !reset_row;
                        // if we've only drawn this row once, draw it again in order to halve our 480 resolution to 240
                        if (reset_row)
                        begin
                            screen_addr_internal <= row_beginning_addr;
                        end
                        else
                        begin
                            if (screen_addr_internal == MAX_ADDRESS_REG)
                            begin
                                // if this is the last register and we've already drawn it twice (reset_row == 0), reset the machine to pixel 0
                                screen_addr_internal <= 0;
                                row_beginning_addr   <= 0;
                                row_end_addr         <= ADDRESSES_PER_SCREEN_ROW - 1;
                                
                            end
                            else
                            begin
                                // else add to the screen_addr_internal, and update the row_beginning_addr and row_end_addr
                                screen_addr_internal <= screen_addr_internal + 1;
                                row_beginning_addr   <= row_beginning_addr + ADDRESSES_PER_SCREEN_ROW;
                                row_end_addr         <= row_end_addr + ADDRESSES_PER_SCREEN_ROW;
                            end
                        end
                    end
                    else
                    begin
                        screen_addr_internal <= screen_addr_internal + 1; // else continue on to the next address
                    end
                end
            end
        end
    end
    
    assign screen_addr = screen_addr_internal;
    assign bit_index   = bit_index_internal;
endmodule
