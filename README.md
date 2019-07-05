# Nand2TetrisFPGA
This project is an FPGA implementation of the Hack computer architecture described in the book [The Elements of Computing Systems](https://www.nand2tetris.org/), By Noam Nisan and Shimon Schocken. Components are described in Verilog and can be found in the `src` directory.

## Testing
Test benches are simulated using Icarus Verilog version 10.2. Test benches are stored in the `test` directory.

#### Build
Build all tests by running the `build.sh` script from the `test` directory. This script runs all of the 
python (3.7.2) scripts located in `test/gen_tv`. These scripts describe logic to build the test vector files 
stored under `test/tvs`. 

#### Run
The test vector files are subsequently loaded by their respective testbenches and 
used to check devices are functioning as expected. Binaries are saved in `test/bin` and can be run to 
confirm proper functioning of each hardware device.

Note that tests must be run from the `test/` directory, owing to the hardcoded paths in the testbenches. 

#### Example
```sh
$ cd test/
$ ./build.sh
$ ./bin/ALU_test
WARNING: ALU_tb.v:44: $readmemb: Standard inconsistency, following 1364-2005.
     18000 tests completed with          0 errors
```
