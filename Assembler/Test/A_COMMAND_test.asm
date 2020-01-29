// Purposefully pedantic as the first test being written in order to
// weed out all potential bugs
@TEST
@TEST // with comment
@TEST / with broken comment
@TEST1
@TEST1 // with number and comment
@TEST1 / with broken comment
@1
@1 // with comment
@1 / with broken comment
@10
@10 // with comment
@10 / with broken comment
@1A
@1A // with comment
@1A / with broken comment
@_.$:
@_.$: // with comment
@_.$: / with broken comment
@_.$:*
*@_.$: // won't be registered as A_COMMAND but tests for invalid beginning syntax
@_.*$: