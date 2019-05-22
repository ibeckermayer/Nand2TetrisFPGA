module ALU_tb();
   reg clk;
   reg  [15:0]  x;
   reg  [15:0]  y;
   reg 	 zx;
   reg 	 nx;
   reg 	 zy;
   reg 	 ny;
   reg 	 f;
   reg 	 no;
   wire  [15:0] out;
   wire 	 zr;
   wire 	 ng;
   reg [17:0] 	 out_expected; // expected output with extra places for zr and ng
   reg [31:0] 	 vectornum, errors; // bookkeeping variables
   reg [55:0] 	 testvectors[18000-1:0]; // array of testvectors; size determined from `wc -l tvs/ALU.tv`


   ALU DUT
     (
      .x(x),
      .y(y),
      .zx(zx),
      .nx(nx),
      .zy(zy),
      .ny(ny),
      .f(f),
      .no(no),
      .out(out),
      .zr(zr),
      .ng(ng)
      );

   // generate clock signal
   always
     begin
	#5 clk = ~clk;		// 10ns period
     end


   initial
     begin
	clk = 1;
	$readmemb("tvs/ALU.tv", testvectors);
	vectornum= 0; errors = 0;
	{x, y, zx, nx, zy, ny, f, no, out_expected} = testvectors[vectornum];
     end // initial begin

   always @(posedge clk)
     begin
	#1;			 // wait time for tick to register
	if ({out, zr, ng} != out_expected) // check that output is expected output
	  begin			 // if error, display error
	     $display("Error at test vector line %d", vectornum+1);
	     $display("x=%b, y=%b, zx=%b, nx=%b, zy=%b, ny=%b, f=%b, no=%b", x, y, zx, nx, zy, ny, f, no);
	     $display("out=         %b", out);
	     $display("zr=          %b", zr);
	     $display("ng=          %b", ng);
	     $display("out_expected=%b", out_expected);
	     errors = errors + 1;
	  end
	{x, y, zx, nx, zy, ny, f, no, out_expected} = testvectors[vectornum];
     end // always @ (posedge clk)

   always @(negedge clk)
     begin
	vectornum = vectornum + 1;
	if (vectornum > 18000-1)
	  begin
	     $display("%d tests completed with %d errors", vectornum, errors);
	     $finish;// End simulation
	  end
     end

endmodule // ALU_tb
