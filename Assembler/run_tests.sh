echo -e "\033[0mRunning tests for A command syntax checking"
./Test/A_COMMAND_test
echo
echo -e "\033[0mRunning tests for L command syntax checking"
./Test/L_COMMAND_test

# Remove artifacts
rm Test/A_COMMAND_test.hack
rm Test/L_COMMAND_test.hack