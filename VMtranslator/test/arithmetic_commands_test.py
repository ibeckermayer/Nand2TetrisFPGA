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
    # Simulate program to the end
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
    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[256] == -1


def test_SimpleNeg():
    # Run the VMtranslator
    VMtranslator('test/SimpleNeg.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleNeg.asm').run())
    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[256] == -8


def test_SimpleEq():
    # Run the VMtranslator
    VMtranslator('test/SimpleEq.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleEq.asm').run())
    # Force set the top of the stack to != -1/0, so that we to confirm memory is being set by arithmetic command
    for i in [256, 257, 258]:
        hack.ram[i] = 1
    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[256] == 0
    assert hack.ram[257] == -1
    assert hack.ram[258] == 0


def test_SimpleGt():
    # Run the VMtranslator
    VMtranslator('test/SimpleGt.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleGt.asm').run())
    # Force set the top of the stack to != -1/0, so that we to confirm memory is being set by arithmetic command
    for i in [256, 257, 258]:
        hack.ram[i] = 1
    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[256] == 0
    assert hack.ram[257] == 0
    assert hack.ram[258] == -1


def test_SimpleLt():
    # Run the VMtranslator
    VMtranslator('test/SimpleLt.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleLt.asm').run())
    # Force set the top of the stack to != -1/0, so that we to confirm memory is being set by arithmetic command
    for i in [256, 257, 258]:
        hack.ram[i] = 1
    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[256] == -1
    assert hack.ram[257] == 0
    assert hack.ram[258] == 0


def test_SimpleAnd():
    # Run the VMtranslator
    VMtranslator('test/SimpleAnd.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleAnd.asm').run())
    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[256] == 2


def test_SimpleOr():
    # Run the VMtranslator
    VMtranslator('test/SimpleOr.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleOr.asm').run())
    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[256] == 3


def test_SimpleNot():
    # Run the VMtranslator
    VMtranslator('test/SimpleNot.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/SimpleNot.asm').run())
    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break

    assert hack.ram[256] == -1
