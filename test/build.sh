#!/usr/bin/env bash
# Bash script to build tests

# build test vector files
python3 gen_tv/gen_pc_tv.py
python3 gen_tv/gen_ram_tv.py

# build test benches
iverilog -o bin/PC_test PC_tb.v ../src/PC.v
iverilog -o bin/RAM_test RAM_tb.v ../src/RAM.v
