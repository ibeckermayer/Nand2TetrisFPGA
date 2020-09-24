set_property PACKAGE_PIN W5 [get_ports clk]
set_property IOSTANDARD LVCMOS33 [get_ports clk]
create_clock -period 10.000 -name sys_clk_pin -waveform {0.000 5.000} -add [get_ports clk]

set_property PACKAGE_PIN W7 [get_ports {seg_out[0]}]
set_property IOSTANDARD LVCMOS33 [get_ports {seg_out[0]}]
set_property PACKAGE_PIN W6 [get_ports {seg_out[1]}]
set_property IOSTANDARD LVCMOS33 [get_ports {seg_out[1]}]
set_property PACKAGE_PIN U8 [get_ports {seg_out[2]}]
set_property IOSTANDARD LVCMOS33 [get_ports {seg_out[2]}]
set_property PACKAGE_PIN V8 [get_ports {seg_out[3]}]
set_property IOSTANDARD LVCMOS33 [get_ports {seg_out[3]}]
set_property PACKAGE_PIN U5 [get_ports {seg_out[4]}]
set_property IOSTANDARD LVCMOS33 [get_ports {seg_out[4]}]
set_property PACKAGE_PIN V5 [get_ports {seg_out[5]}]
set_property IOSTANDARD LVCMOS33 [get_ports {seg_out[5]}]
set_property PACKAGE_PIN U7 [get_ports {seg_out[6]}]
set_property IOSTANDARD LVCMOS33 [get_ports {seg_out[6]}]
set_property PACKAGE_PIN U2 [get_ports enable0]
set_property IOSTANDARD LVCMOS33 [get_ports enable0]
set_property PACKAGE_PIN U4 [get_ports enable1]
set_property IOSTANDARD LVCMOS33 [get_ports enable1]
set_property PACKAGE_PIN V4 [get_ports enable2]
set_property IOSTANDARD LVCMOS33 [get_ports enable2]
set_property PACKAGE_PIN W4 [get_ports enable3]
set_property IOSTANDARD LVCMOS33 [get_ports enable3]
set_property PACKAGE_PIN T18 [get_ports reset]
set_property IOSTANDARD LVCMOS33 [get_ports reset]

set_input_delay -clock [get_clocks sys_clk_pin] -min -add_delay 0.000 [get_ports reset]
set_input_delay -clock [get_clocks sys_clk_pin] -max -add_delay 0.000 [get_ports reset]
