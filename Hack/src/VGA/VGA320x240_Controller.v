module VGA320x240_Controller(input wire clk,
                             input wire reset,
                             input wire [14:0] screen_in,
                             output wire [14:0] screen_addr,
                             output wire vga_hs,
                             output wire vga_vs,
                             output wire [3:0] vga_r,        // 4-bit VGA red output
                             output wire [3:0] vga_g,        // 4-bit VGA green output
                             output wire [3:0] vga_b,        // 4-bit VGA blue output
                             output wire vga_active_o,       // for tb
                             output wire [3:0] bit_index_o,  // for tb
                             output wire pixel_clk_o,        // for tb
                             output wire [9:0] x_o,          // for tb
                             output wire [8:0] y_o);         // for tb
    
    // generate a 25 MHz pixel clock
    wire pixel_clk;
    Pixelclk pixelclk
    (
    .clk(clk),
    .reset(reset),
    .pixel_clk(pixel_clk)
    );
    
    wire [9:0] x;  // current pixel x position: 10-bit value: 0-1023
    wire [8:0] y;  // current pixel y position:  9-bit value: 0-511
    wire vga_active; // high when vga640x480 is drawing and active pixel
    
    VGA640x480Synch synchronizer (
    .clk(clk),
    .pixel_clk(pixel_clk),
    .reset(reset),
    .o_hs(vga_hs),
    .o_vs(vga_vs),
    .o_active(vga_active),
    .o_x(x),
    .o_y(y)
    );
    
    wire [3:0] bit_index;
    Screen_Memory_Counter smc(
    .clk(clk),
    .reset(reset),
    .pixel_clk(pixel_clk),
    .vga_active(vga_active),
    .screen_addr(screen_addr),
    .bit_index(bit_index)
    );
    
    assign vga_active_o          = vga_active;
    assign bit_index_o           = bit_index;
    assign pixel_clk_o           = pixel_clk;
    assign x_o                   = x;
    assign y_o                   = y;
    assign {vga_r, vga_g, vga_b} = vga_active ? (screen_in[bit_index] == 1 ? 12'b111100000000: 12'b000011110000) : 12'b000000000000;
endmodule
