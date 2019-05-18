module PC_tb();
   reg         clk;		 // clk is internal
   reg [15:0]  in;		 // loaded from testvectors
   reg 	       reset, load, inc; // loaded from testvectors
   wire [15:0] out;		 // output

   // instantiate test device
   PC DUT
     (
      .clk(clk),
      .in(in),
      .inc(inc),
      .load(load),
      .reset(reset),
      .out(out)
      );

   // generate clock signal
   always begin
      #5 clk = ~clk;		// 10ns period
   end

   
