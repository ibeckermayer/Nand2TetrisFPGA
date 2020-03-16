module Hex_to_7_Seg (
           input [3:0] hex,
           output  [6:0] seven_seg);

wire [6:0] seg_out;
assign seven_seg = ~seg_out;

always @(hex) begin
    case (hex)
        4'x0: seg_out = 7'b0111111;
        4'x1: seg_out = 7'b0000110;
        4'x2: seg_out = 7'b1011011;
        4'x3: seg_out = 7'b1001111;
        4'x4: seg_out = 7'b1100110;
        4'x5: seg_out = 7'b1101101;
        4'x6: seg_out = 7'b1111101;
        4'x7: seg_out = 7'b0000111;
        4'x8: seg_out = 7'b1111111;
        4'x9: seg_out = 7'b1101111;
        4'xA: seg_out = 7'b1110111;
        4'xB: seg_out = 7'b1111100;
        4'xC: seg_out = 7'b0111001;
        4'xD: seg_out = 7'b1011110;
        4'xE: seg_out = 7'b1111001;
        4'xF: seg_out = 7'b1110001;
      default: seg_out = 7'bxxxxxxx;
    endcase
end

endmodule // Hex_to_7_Seg
