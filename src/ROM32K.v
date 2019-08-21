module ROM32K #(parameter inputfile = "") // binary program to load into memory
               (input [15:0] address,
                output [15:0] out);
    
    reg [15:0]    rom [32767: 0]; // 32k ram at 16-bit address width
    
    initial
    begin
        $readmemb(inputfile, rom);
    end
    
    assign out = rom[address];
    
endmodule // ROM32K
