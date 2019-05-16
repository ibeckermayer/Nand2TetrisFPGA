module Hack
  // The top level module for the Hack computing platform
  (
   input clk,
   input reset
   );

   // connecting wires
   wire data_mem_out_to_cpu_inM;
   wire instr_mem_out_to_cpu_instruction;
   wire [15:0] cpu_outM_to_data_mem_in;
   wire        cpu_writeM_to_data_mem_load;
   wire [15:0] cpu_addressM_to_data_mem_address;
   wire [15:0] cpu_pc_to_rom_address;

   // instantiate instruction memory
   ROM32K instr_mem
     (
      .address(cpu_pc_to_rom_address),	     // input
      .out(instr_mem_out_to_cpu_instruction) // output
      );

   // instantiate CPU
   CPU cpu
     (
      .clk(clk),				      // input
      .inM(data_mem_out_to_cpu_inM),		      // input
      .instruction(instr_mem_out_to_cpu_instruction), // input
      .reset(reset),				      // input
      .outM(cpu_outM_to_data_mem_in),		      // output
      .writeM(cpu_writeM_to_data_mem_load),	      // output
      .addressM(cpu_addressM_to_data_mem_address),    // output
      .pc(cpu_pc_to_rom_address)		      // output
      );

   // instantiate RAM
   RAM #(65536, 16) data_mem
     (
      .clk(clk),				  // input
      .address(cpu_addressM_to_data_mem_address), // input
      .load(cpu_writeM_to_data_mem_load),	  // input
      .in(cpu_outM_to_data_mem_in),		  // input
      .out(data_mem_out_to_cpu_inM)		  // output

      )

endmodule // Hack
