module screen_controller
(
    input wire clk, reset,
    input wire [15:0] current_word,
    output wire hsync, vsync, video_on_out,
    output wire [11:0] rgb,
    output wire [15:0] current_word_address
);

localparam H_DISPLAY = 640; // horizontal display area
localparam V_DISPLAY = 480;   // vertical display area
localparam MAX_ADDRESS_REG = ((H_DISPLAY * V_DISPLAY) / 16) - 1; // 16 is bit-width of current_word

wire video_on; // True if this pixel is to be displayed on the screen (i.e. not in front or back porches or during retrace)
wire p_tick; // 25MHz clock coming fmemory vga_sync

// Instantiate vga_sync
vga_sync vga_sync_unit (.clk(clk), .reset(reset), .hsync(hsync), .vsync(vsync), .video_on(video_on), .p_tick(p_tick));

reg [15:0] current_word_address_reg; // tracks which word in memory we should read for the current pixel (current_word)
assign current_word_address = current_word_address_reg; // output should attach to memory address line
reg [3:0] current_word_index_reg; // tracks which bit in current_word we should translate into a pixel value

always @(posedge p_tick) begin
    if (reset) begin
        current_word_address_reg <= 0;
        current_word_index_reg <= 0;
    end else if (video_on) begin // if this is a pixel_tick for a pixel that should be drawn
        // Counter increment logic:
        // Increment the bit index. It maxes out at 4'b1111 and then rolls over to 4'b0000
        current_word_index_reg <= current_word_index_reg + 1; 
        if (current_word_index_reg == 4'b1111 && current_word_address_reg == MAX_ADDRESS_REG)
            // We've reached the final pixel of the screen, set current_word_address_reg to 0
            current_word_address_reg <= 0;
        else if (current_word_index_reg == 4'b1111)
            // We've reached the final pixel encoded in this word, increment current_word_address
            current_word_address_reg <= current_word_address_reg + 1;
     end
end

// 1 == white pixel, 0 == black pixel
assign rgb = video_on ? (current_word[current_word_index_reg] == 1'b1 ? 12'b000011110000 : 12'b000000001111) : 12'b000000000000;
assign video_on_out = video_on;


endmodule