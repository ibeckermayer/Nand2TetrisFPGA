module RAM 
       (input 		 clk,
        input [15:0] address,
        input 		 load,
        input [15:0] 	 in,
        output [15:0] 	 out,
        output [15:0] screen_out);

reg [15:0] 		   ram [16383:0]; // 16k RAM array
reg [15:0]         screen[24575:16384]; // 8k Screen memory map

always @(posedge clk) begin
    if (load)
        ram[address] <= in;
        screen[address] <= in;
end

assign out = ram[address];
assign screen_out = screen[24575];
endmodule
