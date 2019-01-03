module RAM #(parameter RAM_WIDTH = 8, parameter ADDR_WIDTH = 3)
   (
    input 		   clk,
    input [ADDR_WIDTH-1:0] address,
    input 		   load,
    input [15:0] 	   in,
    output reg [15:0] 	   out
    );

   reg [15:0] 		   ram [RAM_WIDTH-1:0]; // 8-element array of 16-bit wide reg

   always @(posedge clk) begin
      if (load)
        ram[address] <= in;
      out <= ram[address];      // clocked assignment
   end
endmodule // RAM8
