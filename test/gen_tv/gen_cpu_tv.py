'''
script for generating the testvector file for use in ALU_tb.v

Format will be:
(inM)(instruction)(reset)_(outM_exp)(writeM_exp)(addressM_exp)(pc_exp)_(A_exp)(D_exp)
NOTE that clk and reset should be set internally in the testbench

Algorithm:
- begin with a reset to set PC to zero (NOTE: do this in the testbench and set up checker logic s.t. it doesn't check for proper outputs on reset)
'''
