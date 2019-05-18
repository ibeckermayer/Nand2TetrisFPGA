module PC_tb();
   reg         clk;		    // clk is internal
   reg  [15:0] in;		    // loaded from testvectors
   reg 	       reset, load, inc;    // loaded from testvectors
   wire [15:0] out;		    // output
   reg  [15:0] out_expected;	    // expected output
   reg  [31:0] vectornum, errors;   // bookkeeping variables
   reg  [34:0] testvectors[3538:0]; // array of testvectors

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
   always
     begin
	#5 clk = ~clk;		// 10ns period
     end

   // initialize testvectors and bookkeepers
   initial
     begin
	clk = 1;
	$readmemb("tvs/PC.tv", testvectors);
	vectornum= 0; errors = 0;
	{in, reset, load, inc, out_expected} = testvectors[vectornum]; // load test signals into device / out_expected
     end

   // PC is positive edge triggered
   always @(posedge clk)
     begin
	#1;			 // wait time for tick to register
	if (out != out_expected) // check that output is expected output
	  begin			 // if error, display error
	     $display("Error at  in=%d, inc=%d, load=%d, reset=%d", in, inc, load, reset);
	     $display("out=         %d", out);
	     $display("out_expected=%d", out_expected);
	     errors = errors + 1;
	  end
	{in, reset, load, inc, out_expected} = testvectors[vectornum]; // load test signals into device / out_expected
     end

   // check signals on negedge of clock
   always @(negedge clk)
     begin
	vectornum= vectornum+ 1;
	if (vectornum == 3539)
	  begin
	     $display("%d tests completed with %d errors", vectornum - 1, errors);
	     $finish;// End simulation
	  end
     end
endmodule // PC_tb
