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
    hack.ram[int(asmp.symbol_table['ARG'])] = 400
    hack.ram[400] = 10

    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[0] == 257
    assert hack.ram[256] == 55
