# VMtranslator

The VMtranslator constitutes the backend of Jack compiler. It translates stack machine programs (described in chapters 7 and 8 of the book) into equivalent programs in the assembly language.

## Installation

The VMtranslator was written and tested with `python 3.9`.

Install the requirements by running

```
pip install -r requirements.txt
```

## Usage

The `VMtranslator` is used by running

```
python VMtranslator.py $VM_PROGRAM
```

`$VM_PROGRAM` can be either a single `*.vm` file or a directory containing multiple `*.vm` files. The `VMtranslator` always starts execution by calling a function named `Sys.init`, and so `*.vm` programs should anticipate that. More details regarding this design are found in Chapter 8 of the book.

## Testing

Running tests is as simple as running

```
pytest
```

#### Tests

Test files are stored in the `test/` directory. Sample test files to test the `VMtranslator`'s various features are stored as individual files ending in `*.vm` or in subdirectories of `test/` containing multiple `*.vm` files.

The pytest tests themselves are contained in all the files matching `test/*._test`

#### HackAsmSimulator.py

`HackAsmSimulator.py` is effectively a virtual Hack processor. It's used in all of the tests to check that the assembly code translated from the vm code is actually doing what we want it to do when running on the Hack architecture. It contains an `AsmParser` that parses `*.asm` files into an in-memory representation of the program, and then feeds that into the `HackExecutor` which can simulate the execution of the progam.
