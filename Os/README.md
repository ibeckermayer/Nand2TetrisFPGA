Common CLI commands:

```bash
# Set the class name
export CLASSNAME="Math"

# Copy the test files
cp -R ../nand2tetris/projects/12/${CLASSNAME}Test test/
cp ../nand2tetris/tools/OS/*.vm test/${CLASSNAME}Test/ && rm test/${CLASSNAME}Test/${CLASSNAME}.vm
cp ../nand2tetris/projects/12/${CLASSNAME}.jack .
```

After building/editing `./${CLASSNAME}.jack`, compile the test with:

```bash
cp ${CLASSNAME}.jack test/${CLASSNAME}Test && ../Compiler/JackCompiler test/${CLASSNAME}Test
```

Then run the VMEmulator with:

```bash
../nand2tetris/tools/VMEmulator.sh
```

Manually navigate to the `test/${CLASSNAME}Test` directory and run the test.
