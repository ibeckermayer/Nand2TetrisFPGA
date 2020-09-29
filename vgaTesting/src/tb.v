module tb();
reg clk, reset;

always
begin
    #5 clk = ~clk;		// 10ns period
end

initial
begin
    clk = 1;
    reset = 1;
end

always @(posedge clk)
begin
    #20 reset = 0;
end

wire vga_hs, vga_vs, pixel_clk_tb, o_active_tb;
wire [3:0] vga_r, vga_g, vga_b, bit_index_tb;
wire [8:0] y_tb;
wire [9:0] x_tb;
wire [15:0] screen_addr_tb, screen_out_tb;
wire screen_out_of_bit_index = screen_out_tb[bit_index_tb];


top DUT (
    .clk(clk),             // board clock: 100 MHz on Arty/Basys3/Nexys
    .reset(reset),         // reset button
    .vga_hs(vga_hs),       // horizontal sync output
    .vga_vs(vga_vs),       // vertical sync output
    .vga_r(vga_r),    // 4-bit VGA red output
    .vga_g(vga_g),    // 4-bit VGA green output
    .vga_b(vga_b),     // 4-bit VGA blue output

    // Outputs for test bench
    .pixel_clk_tb(pixel_clk_tb), .o_active_tb(o_active_tb),
    .x_tb(x_tb),
    .y_tb(y_tb),
    .screen_addr_tb(screen_addr_tb), .screen_out_tb(screen_out_tb),
    .bit_index_tb(bit_index_tb)
);
endmodule