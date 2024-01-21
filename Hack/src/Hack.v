module Hack (input wire clk,
             input wire reset,
             output wire vga_hs,       // horizontal sync output
             output wire vga_vs,       // vertical sync output
             output wire [3:0] vga_r,  // 4-bit VGA red output
             output wire [3:0] vga_g,  // 4-bit VGA green output
             output wire [3:0] vga_b); // 4-bit VGA blue

    // connecting wires
    wire [15:0] data_mem_out_to_cpu_inM;
    wire [15:0] instr_mem_out_to_cpu_instruction;
    wire [15:0] cpu_outM_to_data_mem_in;
    wire        cpu_writeM_to_data_mem_load;
    wire [15:0] cpu_addressM_to_data_mem_address;
    wire [14:0] cpu_pc_to_rom_address;
    // Driven by vga_ctrl
    wire [14:0] screen_addr;
    wire [14:0] screen_out;

    // instantiate ROM
    RAMROM #(16, "../../Assembler/write_every_other_pixel.hack") instr_mem
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

    // instantiate RAM
    RAMROM #(15) data_mem
    (
    .clk(clk),				                          // input
    .address(cpu_addressM_to_data_mem_address[14:0]), // input, attach lower 15-bits
    .screen_address(screen_addr),                     // input: the memory address to read from for the screen
    .load(cpu_writeM_to_data_mem_load),	              // input
    .in(cpu_outM_to_data_mem_in),		              // input
    .out(data_mem_out_to_cpu_inM),		              // output
    .screen_out(screen_out)                           // output: the data at `screen_address`
    );

    // instantiate VGA controller
    VGA320x240_Controller vga_ctrl
    (
    .clk(clk),                  // input
    .reset(reset),              // input
    // output: this component determines which screen address to read from
    // in the RAMROM above, and outputs it here.
    .screen_addr(screen_addr),
    // input: the data at `screen_addr` read in here to determine the color
    // of the pixel to draw.
    .screen_in(screen_out),
    .vga_hs(vga_hs),
    .vga_vs(vga_vs),
    .vga_r(vga_r),
    .vga_b(vga_b),
    .vga_g(vga_g)
    );

endmodule // Hack
