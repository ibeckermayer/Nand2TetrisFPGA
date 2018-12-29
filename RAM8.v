module RAM8(out, address, in, load, clk);
   output[15:0] out;
   input [15:0] in;
   input [2:0]  address;
   input        load, clk;

   reg          out;            // out is a register
   reg [15:0]   ram [7:0];      // 8-element array of 16-bit wide reg

   always @(posedge clk) begin
      if (load)
        ram[address] <= in;
      out <= ram[address];      // clocked assignment
   end
endmodule // RAM8
