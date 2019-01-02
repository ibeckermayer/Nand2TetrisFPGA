module RAM8(
            input             clk,
            input [2:0]       address,
            input             load,
            input [15:0]      in,
            output reg [15:0] out
            );

   reg [15:0]                 ram [7:0]; // 8-element array of 16-bit wide reg

   always @(posedge clk) begin
      if (load)
        ram[address] <= in;
      out <= ram[address];      // clocked assignment
   end
endmodule // RAM8
