module Hex_to_7_Seg (input [3:0] hex,
                     output [6:0] seven_seg);
    
    reg [6:0] seg_out;
    assign seven_seg = ~seg_out;
    
    always @(hex) begin
        case (hex)
            4'h0: seg_out    = 7'b0111111;
            4'h1: seg_out    = 7'b0000110;
            4'h2: seg_out    = 7'b1011011;
            4'h3: seg_out    = 7'b1001111;
            4'h4: seg_out    = 7'b1100110;
            4'h5: seg_out    = 7'b1101101;
            4'h6: seg_out    = 7'b1111101;
            4'h7: seg_out    = 7'b0000111;
            4'h8: seg_out    = 7'b1111111;
            4'h9: seg_out    = 7'b1101111;
            4'hA: seg_out    = 7'b1110111;
            4'hB: seg_out    = 7'b1111100;
            4'hC: seg_out    = 7'b0111001;
            4'hD: seg_out    = 7'b1011110;
            4'hE: seg_out    = 7'b1111001;
            4'hF: seg_out    = 7'b1110001;
            default: seg_out = 7'bxxxxxxx;
        endcase
    end
    
endmodule // Hex_to_7_Seg
