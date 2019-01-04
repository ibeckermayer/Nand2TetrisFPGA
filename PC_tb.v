module PC_tb();
   reg         clk;
   reg [15:0]  in;
   reg 	       inc;
   reg 	       load;
   reg 	       reset;
   wire [15:0] out;

   PC DUT
     (
      .clk(clk),
      .in(in),
      .inc(inc),
      .load(load),
      .reset(reset),
      .out(out)
      );

   initial begin
      clk = 0;
      in = 20;
      inc = 0;
      load = 0;
      reset = 0;

      #30 load = 1;
      #10 inc = 1; load = 0;
      #50 reset = 1;
      #10 reset = 0;
   end // initial begin

   always #5 clk = ~clk;

   initial #250 $stop;

   initial
     $monitor("At time %t, clk = %0d, load = %0d, inc = %0d, in = %0d, reset = 0%d, out = %0d",
              $time, clk, load, inc, in, reset, out);


endmodule // PC_tb
