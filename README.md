sqroot
======

A package to compute square roots and cube roots to arbitrary precision.

## How this package differs from big.Float.Sqrt

big.Float.Sqrt requires a finite precision to be set ahead of time. The answer it gives is only accurate to that precision. This package does not require a precision to be set because it computes exact square roots almost instantly. These exact square roots compute their digits lazily on an as needed basis. Also this package features cube roots which the big.Float in the standard library does not offer as of this writing.

For doing math operations with square roots, big.Float in the standard library is the best and by far the fastest. However, when printing out square roots in human readable format, this package can be almost twice as fast as big.Float.

## Examples

```golang
package main

import (
    "fmt"

    "github.com/keep94/sqroot/v3"
)

func main() {

    // Print the first 1000 digits of the square root of 2.
    fmt.Printf("%.1000g\n", sqroot.Sqrt(2))

    // Print the 10,000th digit of the square root of 5.
    fmt.Println(sqroot.Sqrt(5).At(9999))

    // Print the location of the first 4 consecutive 0's in the cube root of 7.
    fmt.Println(sqroot.FindFirst(sqroot.CubeRoot(7), []int{0, 0, 0, 0}))

    // Print the location of the last 7 in the first 10,000 digits of the
    // cube root of 11.
    fmt.Println(sqroot.FindLast(sqroot.CubeRoot(11).WithEnd(10000), []int{7}))
}
```

More documentation and examples can be found [here](https://pkg.go.dev/github.com/keep94/sqroot/v3).
