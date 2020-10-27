import sys
import os.path
from HackAsmSimulator import HackExecutor, AsmParser, CT
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.path.pardir)))
from VMtranslator import VMtranslator


def test_BasicTest():
    # Run the VMtranslator
    VMtranslator('test/BasicTest.vm').run()
    # Load the resulting asm file into the HackExecutor
    hack = HackExecutor(AsmParser('test/BasicTest.asm').run())
    # Simulate program to the end
    while True:
        if hack.step().type == CT.END:
            break
    assert hack.ram[11] == 510
    assert hack.ram[3015] == 45
    assert hack.ram[3012] == 42
    assert hack.ram[3006] == 36
    assert hack.ram[402] == 22
    assert hack.ram[401] == 21
    assert hack.ram[300] == 10
    assert hack.ram[256] == 472
    assert hack.ram[3000] == 8
    assert hack.ram[3010] == 7