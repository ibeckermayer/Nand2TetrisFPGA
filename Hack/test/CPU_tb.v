module CPU_tb();
reg clk, reset;
reg [15:0] inM, instruction;
wire writeM;
wire [15:0] outM, addressM, pc;
reg [31:0] 	 vectornum, errors, outM_errors, writeM_errors, addressM_errors, pc_errors; // bookkeeping variables
reg writeM_expected;
reg [15:0] outM_expected, addressM_expected, pc_expected, alu_out_expected, A_expected, D_expected;
reg [130-1:0] 	 testvectors[100000-1:0]; // array of testvectors; size determined from `wc -l tvs/CPU.tv`

localparam DEBUG = 0;

CPU DUT
    (
        .clk(clk),
        .inM(inM),
        .instruction(instruction),
        .reset(reset),
        .outM(outM),
        .writeM(writeM),
        .addressM(addressM),
        .pc(pc)
    );

// generate clock signal
always
begin
    #5 clk = ~clk;		// 10ns period
end

initial
begin
    clk = 1;
    $readmemb("/home/ibeckermayer/Nand2TetrisFPGA/Hack/test/tvs/CPU.tv", testvectors);
    vectornum= 0; errors = 0; outM_errors = 0; writeM_errors = 0; addressM_errors = 0; pc_errors = 0;
    {inM, instruction, reset, outM_expected, writeM_expected, addressM_expected, pc_expected, alu_out_expected, A_expected, D_expected} = testvectors[vectornum];
end // initial begin

always @(posedge clk)
begin
    if (DEBUG)
    begin
        $display("inM=%d, instruction=%b (%d), reset=%b", inM, instruction, instruction, reset);
        // Need to check combinational logic before register tick, or we will be observing incorrect values.
        // Consider if the instruction means "A + D". That means we want out whatever was in the A register last-tick
        // plus whatever was in the D register last-tick. If we wait for the tick to register, the combinational logic
        // (which is effectively updated immediately) will show us an unwanted output.
        $display("CPU.alu_out=         %d", DUT.alu_out);
        $display("CPU.alu_out_expected=%d", $signed(alu_out_expected));
        $display("outM=                %d", $signed(outM));
        $display("outM_expected=       %d", $signed(outM_expected));
    end
    if (outM != outM_expected)
    begin
        errors = errors + 1;
        outM_errors = outM_errors + 1;
    end
    #1;			 // wait time for tick to register
    if (DEBUG)
    begin
        $display("CPU.D=               %d", DUT.D);
        $display("CPU.D_expected=      %d", D_expected);
        $display("CPU.A=               %d", DUT.A);
        $display("CPU.A_expected=      %d", A_expected);
        $display("pc=                  %d", pc);
        $display("pc_expected=         %d", pc_expected);
        $display("writeM=              %b", writeM);
        $display("writeM_expected=     %b", writeM_expected);
        $display("addressM=            %d", addressM);
        $display("addressM_expected=   %d", addressM_expected);
        $display("");
    end


    if ({writeM, addressM, pc} != {writeM_expected, addressM_expected, pc_expected})
    begin
        errors = errors + 1;
    end
    if (writeM != writeM_expected)
    begin
        writeM_errors = writeM_errors + 1;
    end
    if (addressM != addressM_expected)
    begin
        addressM_errors = addressM_errors + 1;
    end
    if (pc != pc_expected)
    begin
        pc_errors = pc_errors + 1;
    end
    {inM, instruction, reset, outM_expected, writeM_expected, addressM_expected, pc_expected, alu_out_expected, A_expected, D_expected} = testvectors[vectornum];
end // always @ (posedge clk)

always @(negedge clk)
begin
    vectornum = vectornum + 1;
    if (vectornum > 100000-1 || errors > 9)
    begin
        $display("%d tests completed with %d errors", vectornum, errors);
        $display("%d outM_errors", outM_errors);
        $display("%d writeM_errors", writeM_errors);
        $display("%d addressM_errors", addressM_errors);
        $display("%d pc_errors", pc_errors);
        $finish;// End simulation
    end
end

endmodule // CPU_tb
