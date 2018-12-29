module RAM8_tb();
   wire [15:0] out;
   reg [15:0]  in;
   reg [2:0]   address;
   reg 	       load;
   reg 	       clk;

   RAM8 DUT (
	     .out(out),
	     .in(in),
	     .address(address),
	     .load(load),
	     .clk(clk)
	     );

   initial begin
      clk = 0;
      load = 0;
      address = 0;
      in = 0;

      #10 load = 1; address = 0; in = 0;
      #10 load = 0; address = 0;
      #10 load = 1; address = 1; in = 1;
      #10 load = 0; address = 1;
      #10 load = 1; address = 2; in = 2;
      #10 load = 0; address = 2;
      #10 load = 1; address = 3; in = 3;
      #10 load = 0; address = 3;
      #10 load = 1; address = 4; in = 4;
      #10 load = 0; address = 4;
      #10 load = 1; address = 5; in = 5;
      #10 load = 1; address = 0; in = 0;
      #10 load = 0; address = 0;
      #10 load = 1; address = 1; in = 1;
      #10 load = 0; address = 1;
      #10 load = 1; address = 2; in = 2;
      #10 load = 0; address = 2;
      #10 load = 1; address = 3; in = 3;
      #10 load = 0; address = 3;
      #10 load = 1; address = 4; in = 4;
      #10 load = 0; address = 4;
      #10 load = 1; address = 5; in = 5;
   end // initial begin

   always #5 clk = ~clk;

   initial #250 $stop;

   initial
     $monitor("At time %t, clk = %0d, load = %0d, address = %0d, in = %0d, out = %0d",
              $time, clk, load, address, in, out);
endmodule // RAM8_tb
