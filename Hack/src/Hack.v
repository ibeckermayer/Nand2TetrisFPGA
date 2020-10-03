module Hack (input clk,
             input reset,
             output [6:0] seg_out,
             output enable0,
             output enable1,
             output enable2,
             output enable3);
    
    // connecting wires
    wire [15:0] data_mem_out_to_cpu_inM;
    wire [15:0] instr_mem_out_to_cpu_instruction;
    wire [15:0] cpu_outM_to_data_mem_in;
    wire        cpu_writeM_to_data_mem_load;
    wire [15:0] cpu_addressM_to_data_mem_address;
    wire [15:0] cpu_pc_to_rom_address;
    wire [6:0]  screen_seg_out_to_Hack_seg_out;
    wire        screen_enable0_to_Hack_enable0;
    wire        screen_enable1_to_Hack_enable1;
    wire        screen_enable2_to_Hack_enable2;
    wire        screen_enable3_to_Hack_enable3;
    
    assign seg_out = screen_seg_out_to_Hack_seg_out;
    assign enable0 = screen_enable0_to_Hack_enable0;
    assign enable1 = screen_enable1_to_Hack_enable1;
    assign enable2 = screen_enable2_to_Hack_enable2;
    assign enable3 = screen_enable3_to_Hack_enable3;
    
    
    // instantiate instruction memory
    RAMROM #(15, "/home/ibeckermayer/Nand2TetrisFPGA/Assembler/RAM_out_attached_to_7_seg_simple.hack") instr_mem
    (
    .clk(clk),
    .address(cpu_pc_to_rom_address),	      // input
    .out(instr_mem_out_to_cpu_instruction), // output
    // ground unused ports
    .screen_address(0),
    .screen_out(0),
    .load(0),
    .in(0)
    );
    
    // instantiate CPU
    CPU cpu
    (
    .clk(clk),				                        // input
    .inM(data_mem_out_to_cpu_inM),		            // input
    .instruction(instr_mem_out_to_cpu_instruction), // input
    .reset(reset),				                    // input
    .outM(cpu_outM_to_data_mem_in),		            // output
    .writeM(cpu_writeM_to_data_mem_load),	        // output
    .addressM(cpu_addressM_to_data_mem_address),    // output
    .pc(cpu_pc_to_rom_address)		                // output
    );
    
    wire [15:0] screen_out_to_7_Seg_hex_in;
    
    // instantiate RAM
    RAMROM #(15) data_mem
    (
    .clk(clk),				                    // input
    .address(cpu_addressM_to_data_mem_address[14:0]), // input, attach lower 15-bits
    .screen_address('h1234),
    .load(cpu_writeM_to_data_mem_load),	        // input
    .in(cpu_outM_to_data_mem_in),		        // input
    .out(data_mem_out_to_cpu_inM),		        // output
    .screen_out(screen_out_to_7_Seg_hex_in)
    );
    
    // instantiate "Screen"
    Hex_to_7_Seg_top screen
    (
    .clk(clk),     // input
    .reset(reset), // input
    .hex_in_0(screen_out_to_7_Seg_hex_in[3:0]),   // input [3:0]
    .hex_in_1(screen_out_to_7_Seg_hex_in[7:4]),   // input [3:0]
    .hex_in_2(screen_out_to_7_Seg_hex_in[11:8]),  // input [3:0]
    .hex_in_3(screen_out_to_7_Seg_hex_in[15:12]), // input [3:0]
    .seg_out(screen_seg_out_to_Hack_seg_out), // output reg [6:0]
    .enable0(screen_enable0_to_Hack_enable0), // output reg
    .enable1(screen_enable1_to_Hack_enable1), // output reg
    .enable2(screen_enable2_to_Hack_enable2), // output reg
    .enable3(screen_enable3_to_Hack_enable3)  // output reg
    );
endmodule // Hack
