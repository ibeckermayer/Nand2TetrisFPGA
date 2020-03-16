module Hex_to_7_Seg_top (
	input clk,
	input reset,
	input [3:0] hex_in_0,
	input [3:0] hex_in_1,
	input [3:0] hex_in_2,
	input [3:0] hex_in_3,
	output reg [6:0] seg_out,
	output reg enable0,
	output reg enable1,	
	output reg enable2,
	output reg enable3);

wire [6:0] seg_out_0;
wire [6:0] seg_out_1;
wire [6:0] seg_out_2;
wire [6:0] seg_out_3;
reg [1:0] toggle = 2'b00;
reg [19:0] refresh_counter; 

Hex_to_7_Seg Hex_to_7_Seg_0 (
	.hex(hex_in_0),
	.seven_seg(seg_out_0)
);

Hex_to_7_Seg Hex_to_7_Seg_1 (
	.hex(hex_in_1),
	.seven_seg(seg_out_1)
);

Hex_to_7_Seg Hex_to_7_Seg_2 (
	.hex(hex_in_2),
	.seven_seg(seg_out_2)
);

Hex_to_7_Seg Hex_to_7_Seg_3 (
	.hex(hex_in_3),
	.seven_seg(seg_out_3)
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
		toggle <= toggle + 1;
end

always @(toggle, seg_out_0, seg_out_1)
begin
	case (toggle)
		2'b00: begin
			seg_out = seg_out_0;
			enable0 = 1'b0;
			enable1 = 1'b1;
			enable2 = 1'b1;
			enable3 = 1'b1;
		end
		2'b01: begin
			seg_out = seg_out_1;
			enable0 = 1'b1;
			enable1 = 1'b0;
			enable2 = 1'b1;
			enable3 = 1'b1;
		end
		2'b10: begin
			seg_out = seg_out_2;
			enable0 = 1'b1;
			enable1 = 1'b1;
			enable2 = 1'b0;
			enable3 = 1'b1;
		end
		2'b11: begin
			seg_out = seg_out_3;
			enable0 = 1'b1;
			enable1 = 1'b1;
			enable2 = 1'b1;
			enable3 = 1'b0;
		end
	endcase
end

endmodule // Hex_to_7_Seg_top
