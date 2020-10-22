import sys
import os.path
from HackAsmSimulator import HackExecutor, AsmParser
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.path.pardir)))
from VMtranslator import VMtranslator


def test_SimpleAdd():
    # Run the VMtranslator
    vmt = VMtranslator('test/SimpleAdd.vm')
    vmt.run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleAdd.asm').run())
    assert len(hack.instructions) > 0
    # Run each instruction
    for i in range(len(hack.instructions)):
        hack.step()
    # Should wind up with 15 on the top of the stack
    assert hack.ram[256] == 15
