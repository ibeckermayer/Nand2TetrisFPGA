import sys
import os.path
from HackAsmSimulator import HackExecutor, AsmParser
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.path.pardir)))
from VMtranslator import VMtranslator


def test_SimpleAdd():
    # Run the VMtranslator
    VMtranslator('test/SimpleAdd.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleAdd.asm').run())
    assert len(hack.instructions) > 0
    # Run each instruction
    for i in range(len(hack.instructions)):
        hack.step()
    # Should wind up with 15 on the top of the stack
    assert hack.ram[256] == 15


def test_SimpleSub():
    # Run the VMtranslator
    VMtranslator('test/SimpleSub.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleSub.asm').run())
    # Run each instruction
    for i in range(len(hack.instructions)):
        hack.step()
    # Should wind up with 15 on the top of the stack
    assert hack.ram[256] == -1


def test_SimpleNeg():
    # Run the VMtranslator
    VMtranslator('test/SimpleNeg.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleNeg.asm').run())
    assert len(hack.instructions) > 0
    # Run each instruction
    for i in range(len(hack.instructions)):
        hack.step()
    # Should wind up with 15 on the top of the stack
    assert hack.ram[256] == -8
