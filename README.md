sqroot
======

A library to compute square roots to arbitrary precision.

### Examples

Print the first 1000 digits of the square root of 2

```golang
fmt.Printf("%.1000g\n", sqroot.Sqrt(2))
```

Print the 10,000th digit of the square root of 5.

```golang
fmt.Println(sqroot.Sqrt(5).Mantissa().At(9999))
```

Print where the first 4 consecutive 0's start in the cube root of 7.

```golang
fmt.Println(sqroot.FindFirst(
    sqroot.CubeRoot(7).Mantissa(), []int{0, 0, 0, 0}))
```

Print the location of the last 7 in the first 10,000 digits of the cube root of 11.

```golang
fmt.Println(sqroot.FindLast(
    sqroot.CubeRoot(11).WithSignificant(10000).Mantissa(),
    []int{7}))
```

More documentation and examples can be found [here](https://pkg.go.dev/github.com/keep94/sqroot).
