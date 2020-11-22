import sys
import os.path
from HackAsmSimulator import HackExecutor, AsmParser, CT
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.path.pardir)))
from VMtranslator import VMtranslator


def test_FibonacciElement():
    # Run the VMtranslator
    VMtranslator('test/FibonacciElement').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/FibonacciElement.asm').run())
    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break
    assert hack.ram[0] == 257
    assert hack.ram[256] == 21


def test_SimpleFunction():
    # Run the VMtranslator
    VMtranslator('test/SimpleFunction').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleFunction.asm').run())
    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break
    assert hack.ram[0] == 257
    assert hack.ram[256] == 1196


def test_NestedCall():
    # Run the VMtranslator
    VMtranslator('test/NestedCall').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/NestedCall.asm').run())
    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[0] == 256
    assert hack.ram[3] == 4000
    assert hack.ram[4] == 5000
    assert hack.ram[5] == 135
    assert hack.ram[6] == 246


def test_StaticsTest():
    # Run the VMtranslator
    VMtranslator('test/StaticsTest').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/StaticsTest.asm').run())
    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[0] == 258
    assert hack.ram[256] == -2
    assert hack.ram[257] == 8
