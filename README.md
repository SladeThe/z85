# Z85

Base-85 encoding and decoding

The library is inspired by https://github.com/tilinna/z85 and partially based on its codebase.

The key improvements compare to the predecessor are:

- automatic padding of input data, source length divisible by 4 is no longer required;
- low level optimizations, now it works about 1.5 times faster.
