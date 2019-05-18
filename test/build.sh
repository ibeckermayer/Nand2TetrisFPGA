#!/usr/bin/env bash
# Bash script to build tests

# build test vector files
python3 gen_tv/gen_pc_tv.py

# build test benches
iverilog -o bin/PC_test PC_tb.v ../src/PC.v
