module Hex_to_7_Seg_top (
	input clk,
	input reset,
	input [3:0] hex_in_0,
	input [3:0] hex_in_1,
	output reg [6:0] seg_out,
	output enable0,
	output enable1,	
	output enable2,
	output enable3);

wire [6:0] seg_out_0;
wire [6:0] seg_out_1;
reg [1:0] toggle = 2'b10;
assign enable0 = toggle[0];
assign enable1 = toggle[1];
assign enable2 = 1'b1;
assign enable3 = 1'b1;

reg [19:0] refresh_counter; 

Hex_to_7_Seg Hex_to_7_Seg_0 (
	.hex(hex_in_0),
	.seven_seg(seg_out_0)
);

Hex_to_7_Seg Hex_to_7_Seg_1 (
	.hex(hex_in_1),
	.seven_seg(seg_out_1)
);

always @(posedge clk)
begin 
 if(reset)
  refresh_counter <= 0;
 else
  refresh_counter <= refresh_counter + 1;
end

always @(posedge clk)
begin
	if (refresh_counter == {20{1'b1}})
		toggle = ~toggle;
end

always @(toggle, seg_out_0, seg_out_1)
begin
	if (toggle[1] == 1'b1)
		seg_out = seg_out_1;
	else
		seg_out = seg_out_0;
end

endmodule // Hex_to_7_Seg_top
