module RAMROM #(parameter ADDR_WIDTH = 16,
                parameter INPUT_FILE = "")      // ram width == 2^(ADDR_WIDTH)
               (input 		 clk,
                input [ADDR_WIDTH-1:0] address,
                input 		 load,
                input [15:0] 	 in,
                output reg [15:0] 	 out);
    
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
        out          <= mem[address];  // synchronous RAM (BRAM)
    end
endmodule // RAM8
