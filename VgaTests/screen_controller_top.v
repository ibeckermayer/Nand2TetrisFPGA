module screen_controller_top
	(
        input wire clk, reset,
		output wire hsync, vsync,
		output wire [11:0] rgb
	);

    wire [15:0] current_word;
    wire [15:0] current_word_address;
    screen_controller screen_controller_unit 
    (
        .clk(clk), .reset(reset), .current_word(current_word), 
        .hsync(hsync), .vsync(vsync), .rgb(rgb), .current_word_address(current_word_address)
    );

    // instantiate instruction memory
    ROM32K #("/home/ibeckermayer/VGA_test/src/black_white_lines_rom.txt") instr_mem
    (
        .address(current_word_address),	      // input
        .out(current_word) // output
    );
endmodule