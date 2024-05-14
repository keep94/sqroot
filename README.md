sqroot
======

A package to compute square roots and cube roots to arbitrary precision.

This package is dedicated to my mother, who taught me how to calculate square roots by hand as a child.

## How this package differs from big.Float.Sqrt

big.Float.Sqrt requires a finite precision to be set ahead of time. The answer it gives is only accurate to that precision. This package does not require a precision to be set in advance. Square root values in this package compute their digits lazily on an as needed basis just as one would compute square roots by hand. Also, this package features cube roots which big.Float in the standard library does not offer as of this writing. Cube root values in this package also compute their digits lazily on an as needed basis just as one would compute cube roots by hand.

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
