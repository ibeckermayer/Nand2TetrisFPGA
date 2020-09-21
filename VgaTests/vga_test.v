module vga_test
	(
        input wire clk, reset,
		input wire [11:0] sw,
		output wire hsync, vsync,
		output wire [11:0] rgb
	);
	
	// register for Basys 2 8-bit RGB DAC 
	reg [11:0] rgb_reg;
	
	// video status output from vga_sync to tell when to route out rgb signal to DAC
	wire video_on;

        // instantiate vga_sync
        vga_sync vga_sync_unit (.clk(clk), .reset(reset), .hsync(hsync), .vsync(vsync),
                                .video_on(video_on), .p_tick());
   
        // rgb buffer
        always @(posedge clk, posedge reset)
        if (reset)
            rgb_reg <= 0;
        else
            rgb_reg <= sw;
        
        // output
        assign rgb = (video_on) ? rgb_reg : 12'b0;
endmodule