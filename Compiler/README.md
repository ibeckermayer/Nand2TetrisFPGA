# Compiler

#### Build

To build the compiler itself:

```
go build -o JackCompiler main.go
```

To build all the test programs:

```
# Loop through each directory in the test/ directory
for dir in test/*/; do
  # Run the JackCompiler command for the current directory
  ./JackCompiler "$dir"
done
```

#### Testing

The Compiler takes advantage of the VMEmulator [provided on the nand2tetris website](https://www.nand2tetris.org/software). Download the software suite and run the `VMEmulator` Java program.

The test programs under the `test/` directory are copied (and optionally modified) from the same software suite. To compile any of the tests simply run `./JackCompiler test/<test_name>`. Each working test contains the "operating system" VM files it needs to be run by the emulator, so to run the test you can just load the given directory in the emulator and press play. Consult chapter 11 of the book for the expected behavior of each test.

**TODO**

The testing is currently done by manually running and checking the test output against the software-suite-provided emulator; however the VM itself isn't all that complex and so it wouldn't be too difficult to build another emulator in go that could be used in an automated test suite.

Alternatively there is likely some way to hook into the provided emulator such that you could load the program, run the simulation, and then check the output programatically.

Tests TODO:

- ~Seven~
- ~ConvertToBin~
- ~Square~
- ~Average~
- ~Pong~
- ~ComplexArrays~
