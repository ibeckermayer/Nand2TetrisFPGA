module vga_tb();
reg 	       clk, reset;
reg [15:0]  current_word;
reg [63:0] tocks;
wire hsync, vsync;
wire [11:0] rgb;
wire [15:0] current_word_address;
wire video_on_out;

screen_controller DUT (.clk(clk), .reset(reset), .current_word(current_word), 
        .hsync(hsync), .vsync(vsync), .rgb(rgb), .current_word_address(current_word_address), .video_on_out(video_on_out));

// reg [31:0] tocks;
// wire hsync, vsync, video_on, p_tick;
// vga_sync DUT (.clk(clk), .reset(reset), .hsync(hsync), .vsync(vsync), .video_on(video_on), .p_tick(p_tick));

always
begin
    #5 clk = ~clk;		// 10ns period
end

initial
begin
    clk = 1;
    current_word = 16'b1111000011110000;
    tocks = 0;
    reset = 1;
end

always @(posedge clk)
begin
    #1;			// wait time for register
    reset = 0;
end // always @ (posedge clk)

// tcl command: run 100000 ns

// check signals on negedge of clock
// always @(negedge clk)
// begin
//     tocks = tocks + 1;
//     if (tocks > 32767)
//         $finish;// End simulation
// end

endmodule // ROM_tb
