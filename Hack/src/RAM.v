module RAMROM #(parameter ADDR_WIDTH = 16)       // ram width == 2^(ADDR_WIDTH)
       (input 		 clk,
        input [ADDR_WIDTH-1:0] address,
        input 		 load,
        input [15:0] 	 in,
        output reg [15:0] 	 out);

reg [15:0] 		   ram [2**ADDR_WIDTH-1:0]; // 2^ADDR_WIDTH array of 16-bit wide reg

always @(posedge clk) begin
    if (load)
        ram[address] <= in;
    out <= ram[address];  // synchronous RAM (BRAM)
end
endmodule // RAM8
