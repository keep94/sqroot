package sqroot_test

import (
	"fmt"

	"github.com/keep94/consume2"
	"github.com/keep94/sqroot"
)

func ExampleSqrt() {

	// Print the square root of 13 with 50 significant digits.
	fmt.Printf("%.50g", sqroot.Sqrt(13))
	// Output:
	// 3.6055512754639892931192212674704959462512965738452
}

func ExampleMantissa_Find() {

	// sqrt(2) = 0.14142135... * 10^1
	n := sqroot.Sqrt(2)

	// '14' matches at index 0, 2, 144, ...
	matches := n.Mantissa().Find([]int{1, 4})

	fmt.Println(matches())
	fmt.Println(matches())
	fmt.Println(matches())
	// Output:
	// 0
	// 2
	// 144
}

func ExampleMantissa_FindAll() {

	// sqrt(2) = 0.14142135... * 10^1
	n := sqroot.Sqrt(2)

	var matches []int
	n.Mantissa().FindAll(
		[]int{1, 4}, consume2.Slice(consume2.AppendTo(&matches), 0, 3))
	fmt.Println(matches)
	// Output:
	// [0 2 144]
}

func ExampleMantissa_FindFirst() {

	// sqrt(3) = 0.1732050807... * 10^1
	n := sqroot.Sqrt(3)

	fmt.Println(n.Mantissa().FindFirst([]int{0, 5, 0, 8}))
	// Output:
	// 4
}

func ExampleMantissa_FindFirstN() {

	// sqrt(2) = 0.14142135... * 10^1
	n := sqroot.Sqrt(2)

	fmt.Println(n.Mantissa().FindFirstN([]int{1, 4}, 3))
	// Output:
	// [0 2 144]
}

func ExampleMantissa_Iterator() {

	// sqrt(7) = 0.264575... * 10^1
	n := sqroot.Sqrt(7)

	iter := n.Mantissa().Iterator()

	fmt.Println(iter())
	fmt.Println(iter())
	fmt.Println(iter())
	fmt.Println(iter())
	fmt.Println(iter())
	fmt.Println(iter())
	// Output:
	// 2
	// 6
	// 4
	// 5
	// 7
	// 5
}

func ExampleMantissa_Print() {

	// Find the square root of 2.
	n := sqroot.Sqrt(2)

	fmt.Printf("10^%d *\n", n.Exponent())
	n.Mantissa().Print(
		1000,
		sqroot.DigitsPerRow(50),
		sqroot.DigitsPerColumn(5),
		sqroot.ShowCount(true))
	// Output:
	// 10^1 *
	//    0.14142 13562 37309 50488 01688 72420 96980 78569 67187 53769
	//  50  48073 17667 97379 90732 47846 21070 38850 38753 43276 41572
	// 100  73501 38462 30912 29702 49248 36055 85073 72126 44121 49709
	// 150  99358 31413 22266 59275 05592 75579 99505 01152 78206 05714
	// 200  70109 55997 16059 70274 53459 68620 14728 51741 86408 89198
	// 250  60955 23292 30484 30871 43214 50839 76260 36279 95251 40798
	// 300  96872 53396 54633 18088 29640 62061 52583 52395 05474 57502
	// 350  87759 96172 98355 75220 33753 18570 11354 37460 34084 98847
	// 400  16038 68999 70699 00481 50305 44027 79031 64542 47823 06849
	// 450  29369 18621 58057 84631 11596 66871 30130 15618 56898 72372
	// 500  35288 50926 48612 49497 71542 18334 20428 56860 60146 82472
	// 550  07714 35854 87415 56570 69677 65372 02264 85447 01585 88016
	// 600  20758 47492 26572 26002 08558 44665 21458 39889 39443 70926
	// 650  59180 03113 88246 46815 70826 30100 59485 87040 03186 48034
	// 700  21948 97278 29064 10450 72636 88131 37398 55256 11732 20402
	// 750  45091 22770 02269 41127 57362 72804 95738 10896 75040 18369
	// 800  86836 84507 25799 36472 90607 62996 94138 04756 54823 72899
	// 850  71803 26802 47442 06292 69124 85905 21810 04459 84215 05911
	// 900  20249 44134 17285 31478 10580 36033 71077 30918 28693 14710
	// 950  17111 16839 16581 72688 94197 58716 58215 21282 29518 48847
}

func ExampleMantissa_Send() {
	var mantissaDigits []int

	// sqrt(1000 / 3) = 0.1825741858...*10^2
	n := sqroot.SqrtRat(1000, 3)

	n.Mantissa().Send(consume2.Slice(consume2.AppendTo(&mantissaDigits), 0, 10))
	fmt.Println(mantissaDigits)
	fmt.Println(n.Exponent())
	// Output:
	// [1 8 2 5 7 4 1 8 5 8]
	// 2
}

func ExampleNumber_WithSignificant() {

	// n is 1.42857142857... but truncated to 10000 significant digits
	n := sqroot.SqrtRat(100, 49).WithSignificant(10000)

	// Instead of running forever, FindFirst returns -1 because n is truncated.
	fmt.Println(n.Mantissa().FindFirst([]int{1, 1, 2}))
	// Output:
	// -1
}
