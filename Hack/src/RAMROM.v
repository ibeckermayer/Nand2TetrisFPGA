module RAMROM #(parameter ADDR_WIDTH = 16,
                parameter INPUT_FILE = "")             // ram width == 2^(ADDR_WIDTH)
               (input clk,
                input [ADDR_WIDTH-1:0] address,
                input [ADDR_WIDTH-1:0] screen_address,
                input load,
                input [15:0] in,
                output reg [15:0] out,
                output reg [15:0] screen_out);
    
    reg [15:0] 		   mem [2**ADDR_WIDTH-1:0]; // 2^ADDR_WIDTH array of 16-bit wide reg
    
    initial
    begin
        if (INPUT_FILE != "")
        begin
            $readmemb(INPUT_FILE, mem);
        end
    end
    
    always @(posedge clk)
    begin
        if (load)
        begin
            mem[address] <= in;
        end
        out        <= mem[address];  // synchronous RAM (BRAM)
        screen_out <= mem[screen_address]; // second port for screen output
    end
endmodule // RAM8
