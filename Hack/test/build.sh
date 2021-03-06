#!/usr/bin/env bash
# Bash script to build tests

# build test vector files
python3 gen_tv/gen_pc_tv.py
python3 gen_tv/gen_ram_tv.py
python3 gen_tv/gen_rom_tv.py
python3 gen_tv/gen_alu_tv.py
python3 gen_tv/gen_cpu_tv.py

# build test benches
iverilog -o bin/PC_test PC_tb.v ../src/PC.v
iverilog -o bin/RAM_test RAM_tb.v ../src/RAMROM.v
iverilog -o bin/ROM_test ROM_tb.v ../src/RAMROM.v
iverilog -o bin/ALU_test ALU_tb.v ../src/ALU.v
iverilog -o bin/CPU_test CPU_tb.v ../src/CPU.v ../src/ALU.v ../src/PC.v
