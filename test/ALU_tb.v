module ALU_tb();
   reg signed [15:0]  x;
   reg signed [15:0]  y;
   reg 	 zx;
   reg 	 nx;
   reg 	 zy;
   reg 	 ny;
   reg 	 f;
   reg 	 no;
   wire signed [15:0] out;
   wire 	 zr;
   wire 	 ng;

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
   initial begin
      x = 65535;
      y = 65535;
      zx = 0;
      nx = 0;
      zy = 0;
      ny = 0;
      f  = 0;
      no = 0;
   end // initial begin

   initial #20 $stop;

   initial $monitor("x=%0d, y=%0d, zx=%0d, nx=%0d, zy=%0d, ny=%0d, f=%0d, no=%0d, out=%0d, zr=%0d, ng=%0d,", x, y, zx, nx, zy, ny, f, no, out, zr, ng);

endmodule // ALU_tb
