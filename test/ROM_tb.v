module ROM_tb();
    reg [15:0]  address;
    wire [15:0] out;
    reg [15:0]  out_expected;
    reg [31:0]  vectornum, errors;   // bookkeeping variables
    reg [31:0]  testvectors[32767:0]; // array of testvectors;
    reg 	       clk;
    
    ROM32K #("tvs/ROM_input.tv") DUT (
    .address(address),
    .out(out)
    );
    
    always
    begin
    #5 clk = ~clk;		// 10ns period
    end
    
    initial
    begin
        clk = 1;
        $readmemb("tvs/ROM.tv", testvectors);
        vectornum               = 0; errors               = 0;
        {address, out_expected} = testvectors[vectornum];
    end
    
    always @(posedge clk)
    begin
        #1;			// wait time for register
        if (out ! = out_expected) // check that output is expected output
            begin			 // if error, display error
            $display("Error at test vector line %d", vectornum+1);
            $display("address      = %d", address);
            $display("out          = %d", out);
            $display("out_expected = %d", out_expected);
            errors                 = errors + 1;
            end
            {address, out_expected} = testvectors[vectornum];
            end // always @ (posedge clk)
        
        // check signals on negedge of clock
        always @(negedge clk)
        begin
            vectornum = vectornum + 1;
            if (vectornum > 32767)
            begin
                $display("%d tests completed with %d errors", vectornum, errors);
                $finish;// End simulation
            end
        end
        
        endmodule // ROM_tb
