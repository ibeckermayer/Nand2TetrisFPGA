module RAM #(parameter RAM_WIDTH = 8,
             parameter ADDR_WIDTH = 3)       // RAM_WIDTH == 2^(ADDR_WIDTH)
            (input 		 clk,
             input [ADDR_WIDTH-1:0] address,
             input 		 load,
             input [15:0] 	 in,
             output [15:0] 	 out);
    
    reg [15:0] 		   ram [RAM_WIDTH-1:0]; // RAM_WIDTH-element array of 16-bit wide reg
    
    always @(posedge clk) begin
        if (load)
            ram[address] <= in;
            end
            assign out = ram[address];	// asynchronous RAM
            endmodule // RAM8
