module ROM32K
  (
   input [14:0]  address,
   output [15:0] out
   );

   reg [15:0]    ram [32767: 0]; // 32k ram

   assign out = ram[address];
   
endmodule // ROM32K
