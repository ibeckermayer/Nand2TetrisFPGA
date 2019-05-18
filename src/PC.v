module PC
  (
   input 	     clk,
   input [15:0]      in,
   input 	     inc,
   input 	     load,
   input 	     reset,
   output reg [15:0] out
   );

   always @(posedge clk) begin
      if (reset)
	out <= 0;
      else if (load)
	out <= in;
      else if (inc)
	out <= out + 1;
      else
	out <= out;
   end

endmodule // PC
