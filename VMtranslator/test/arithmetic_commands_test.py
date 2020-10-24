import sys
import os.path
from HackAsmSimulator import HackExecutor, AsmParser, CT
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.path.pardir)))
from VMtranslator import VMtranslator


def test_SimpleAdd():
    # Run the VMtranslator
    VMtranslator('test/SimpleAdd.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleAdd.asm').run())
    # Run each instruction
    while True:
        if hack.step().type == CT.END:
            break
    # Should wind up with 15 on the top of the stack
    assert hack.ram[256] == 15


def test_SimpleSub():
    # Run the VMtranslator
    VMtranslator('test/SimpleSub.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleSub.asm').run())
    # Run each instruction
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[256] == -1


def test_SimpleNeg():
    # Run the VMtranslator
    VMtranslator('test/SimpleNeg.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleNeg.asm').run())
    # Run each instruction
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[256] == -8


def test_SimpleEq():
    # Run the VMtranslator
    VMtranslator('test/SimpleEq.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleEq.asm').run())
    # Run each instruction
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[256] == -1
    assert hack.ram[257] == 0
