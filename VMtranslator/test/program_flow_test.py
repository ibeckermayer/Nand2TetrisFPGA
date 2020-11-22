import sys
import os.path
from HackAsmSimulator import HackExecutor, AsmParser, CT
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.path.pardir)))
from VMtranslator import VMtranslator


def test_BasicLoop():
    # Run the VMtranslator
    VMtranslator('test/BasicLoop.vm').run()

    # Load the resulting asm file into the HackExecutor
    asmp = AsmParser('test/BasicLoop.asm')
    hack = HackExecutor(asmp.run())

    # Manual setup for this test
    hack.ram[int(asmp.symbol_table['LCL'])] = 300
    hack.ram[int(asmp.symbol_table['ARG'])] = 400
    hack.ram[400] = 10

    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[256] == 55


def test_FibonacciSeries():
    # Run the VMtranslator
    VMtranslator('test/FibonacciSeries.vm').run()

    # Load the resulting asm file into the HackExecutor
    asmp = AsmParser('test/FibonacciSeries.asm')
    hack = HackExecutor(asmp.run())

    # Manual setup for this test
    hack.ram[int(asmp.symbol_table['ARG'])] = 400
    hack.ram[400] = 6  # 6 elements of the fibonacci series
    hack.ram[401] = 3000  # Starting at address 3000

    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[3000] == 0
    assert hack.ram[3001] == 1
    assert hack.ram[3002] == 1
    assert hack.ram[3003] == 2
    assert hack.ram[3004] == 3
    assert hack.ram[3005] == 5
