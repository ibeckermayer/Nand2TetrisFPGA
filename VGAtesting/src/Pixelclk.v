module Pixelclk(input wire clk,
                input wire reset,
                output wire pixel_clk);
    // generate a 25 MHz pixel clock
    reg [15:0] cnt;
    reg pixel_clk_internal;
    always @(posedge clk)
    begin
        if (reset)
        begin
            cnt                <= 0;
            pixel_clk_internal <= 0;
        end
        else
            {pixel_clk_internal, cnt} <= cnt + 16'h4000;  // divide by 4: (2^16)/4 = 0x4000
    end
    
    assign pixel_clk = pixel_clk_internal;
endmodule
