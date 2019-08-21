module CPU (input 	 clk,
            input [15:0] inM,
            input [15:0] instruction,
            input 	 reset,
            output [15:0] outM,
            output 	 writeM,
            output [15:0] addressM,
            output [15:0] pc);

// initialize A and D registers
reg [15:0] 	     A;
reg [15:0] 	     D;

// initialize wires
wire [15:0] 	     D_to_alu_x;
wire [15:0] 	     A_or_inM;
wire 	     a;
wire 	     c1;
wire 	     c2;
wire 	     c3;
wire 	     c4;
wire 	     c5;
wire 	     c6;
wire 	     d1;
wire 	     d2;
wire 	     d3;
wire 	     j1;
wire 	     j2;
wire 	     j3;
wire 	     is_j1;	// intermediate wire for determining j1 logic
wire 	     is_j2;	// intermediate wire for determining j2 logic
wire 	     is_j3;	// intermediate wire for determining j3 logic
wire 	     jump;	// final wire determining whether to jump or not
wire signed [15:0] alu_out;
wire 	      alu_zr;
wire 	      alu_ng;

// initialize ALU
ALU _alu
    (
        .x(D),			// input
        .y(A_or_inM),		// input
        .zx(c1),			// input
        .nx(c2),			// input
        .zy(c3),			// input
        .ny(c4),			// input
        .f(c5),			// input
        .no(c6),			// input
        .out(alu_out),		// output
        .zr(alu_zr),		// output
        .ng(alu_ng)		// output
    );

// initialize PC
PC _pc
   (
       .clk(clk),
       .in(A),
       .inc(1'b1),		// over-ridden by load, can be set to 1
       .load(jump),		// over-ridden by reset
       .reset(reset),		// highest priority
       .out(pc)
   );

// The Hack assembly language specification denotes two distinct types of instructions:
// A-instruction: MSB == 0, loads the A register with value V = 0vvvvvvvvvvvvvvv
// C-instruction: MSB == 1, instruction = xx a cccccc ddd jjj
//                                             123456 123 123

// A-instruction logic
// A-instruction logic is simple: if MSB == 0, load the A register
always @(posedge clk) begin
    if (!(instruction[15])) begin // if A-instruction
        A < = instruction; // A gets the value V = 0vvvvvvvvvvvvvvv
    end
end

// C-instruction logic
// C-instruction logic is considerably more involved. See chapter 4 in the book for reference

// ALU connections
// Connections to the ALU are self explanatory, and can derived from the specification in the book
// These connections are instruction-type-agnostic, by which I mean they can be specified as continuous
// assignments without regard for whether we're recieving an A-instruction or a C-instruction.
// They only affect the input and therefore output of the ALU, which on A-instructions can be junk values.

// assign wires to analogous instruction bits
// 1010101010101010
assign a        = instruction[12];
assign A_or_inM = (a) ? inM : A; //  // a-bit determines whether alu's y input is inM or A-register
assign c1       = instruction[11];
assign c2       = instruction[10];
assign c3       = instruction[9];
assign c4       = instruction[8];
assign c5       = instruction[7];
assign c6       = instruction[6];

// Destination logic
// Unlike ALU connections, destination logic is not instruction-type-agnostic.
// If we are dealing with an A-instruction, we don't want our CPU writing to any registers
// due to the binary representation of the number we are instructing it to write into the A
// register. For this reason, our d wires will only be set to their corresponding d values
// if we are dealing with a C-instruction, otherwise they will all be set to 0 which means
// "the value is not stored anywhere" according to Hack language semantics.
assign d1 = (instruction[15]) ? instruction[5] : 0;
assign d2 = (instruction[15]) ? instruction[4] : 0;
assign d3 = (instruction[15]) ? instruction[3] : 0;

// Memory logic can be combinational assignment, registers are clocked RAM module
assign outM = alu_out;	// outM connects to ALU out
assign addressM = A;		// A specifies address
assign writeM = d3;		// d3 determines whether M[A] = alu_out

// register logic
always @(posedge clk) begin
    if (d1) begin
        A <= alu_out;
    end
    if (d2) begin
        D <= alu_out;
    end
end

// Jump logic
// Similar to destination logic, jump logic is not instruction-type-agnostic.
assign j1 = (instruction[15]) ? instruction[2] : 0;
assign j2 = (instruction[15]) ? instruction[1] : 0;
assign j3 = (instruction[15]) ? instruction[0] : 0;

// logic specific to each j value, according to the specification
assign is_j1 = (j1) ? ((alu_ng) ? 1 : 0) : 0; // is j1 true and alu_out < 0?
assign is_j2 = (j2) ? ((alu_zr) ? 1 : 0) : 0; // is j2 true and alu_out = 0?
assign is_j3 = (j3) ? ((!(alu_ng) && !(alu_zr)) ? 1 : 0) : 0; // is j3 true and alu_out > 0?

// final jump determination, fed into PC
assign jump = (is_j1 || is_j2 || is_j3);

endmodule // CPU
