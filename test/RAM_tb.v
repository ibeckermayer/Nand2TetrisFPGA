module RAM_tb();
   wire [15:0] out;
   reg [15:0]  out_expected;
   reg [15:0]  in;
   reg [15:0]  address;
   reg 	       load;
   reg 	       clk;
   reg [31:0]  vectornum, errors;   // bookkeeping variables
   reg [48:0]  testvectors[65542:0]; // array of testvectors; size determined from `wc -l tvs/RAM.tv`

   RAM #(32767, 16) DUT (	// 32K RAM as in final architecture,
				// still has 16 bit address but shouldn't
				// ever go beyond 0111111111111111 owing
				// to the nature of the A-instruction
	     .out(out),
	     .in(in),
	     .address(address),
	     .load(load),
	     .clk(clk)
	     );

   // generate clock signal
   always
     begin
	#5 clk = ~clk;		// 10ns period
     end

   // initialize clk, testvectors and bookkeepers
   initial
     begin
	clk = 1;
	$readmemb("tvs/RAM.tv", testvectors);
	vectornum= 0; errors = 0;
	{address, in, load, out_expected} = testvectors[vectornum]; // load test signals into device / out_expected
     end

   always @(posedge clk)
     begin
	#1;			// wait time for register
	if (out != out_expected) // check that output is expected output
	  begin			 // if error, display error
	     $display("Error at test vector line %d", vectornum+1);
	     $display("in=%d, address=%d, load=%d", in, address, load);
	     $display("out=         %d", out);
	     $display("out_expected=%d", out_expected);
	     errors = errors + 1;
	  end
	{address, in, load, out_expected} = testvectors[vectornum]; // load test signals into device / out_expected
     end // always @ (posedge clk)

   // check signals on negedge of clock
   always @(negedge clk)
     begin
	vectornum = vectornum + 1;
	if (vectornum > 65542)
	  begin
	     $display("%d tests completed with %d errors", vectornum, errors);
	     $finish;// End simulation
	  end
     end
endmodule // RAM_tb
