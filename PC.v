module PC
  (
   input clk,
   input in [15:0],
   input inc,
   input load,
   input reset,
   output reg out [15:0]
  )

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
