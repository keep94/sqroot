package sqroot_test

import (
	"fmt"
	"slices"

	"github.com/keep94/sqroot/v3"
)

func ExampleSqrt() {

	// Print the square root of 13 with 100 significant digits.
	fmt.Printf("%.100g\n", sqroot.Sqrt(13))
	// Output:
	// 3.605551275463989293119221267470495946251296573845246212710453056227166948293010445204619082018490717
}

func ExampleCubeRoot() {

	// Print the cube root of 3 with 100 significant digits.
	fmt.Printf("%.100g\n", sqroot.CubeRoot(3))
	// Output:
	// 1.442249570307408382321638310780109588391869253499350577546416194541687596829997339854755479705645256
}

func ExampleNewNumberForTesting() {

	// n = 10.2003400340034...
	n, _ := sqroot.NewNumberForTesting([]int{1, 0, 2}, []int{0, 0, 3, 4}, 2)

	fmt.Println(n)
	// Output:
	// 10.20034003400340
}

func ExampleNewFiniteNumber() {

	// n = 563.5
	n, _ := sqroot.NewFiniteNumber([]int{5, 6, 3, 5}, 3)

	fmt.Println(n.Exact())
	// Output:
	// 563.5
}

func ExampleAsString() {

	// sqrt(3) = 0.1732050807... * 10^1
	n := sqroot.Sqrt(3)

	fmt.Println(sqroot.AsString(n.WithStart(2).WithEnd(10)))
	// Output:
	// 32050807
}

func ExampleMatches() {

	// sqrt(2) = 0.14142135... * 10^1
	n := sqroot.Sqrt(2)

	// '14' matches at index 0, 2, 144, ...
	count := 0
	for index := range sqroot.Matches(n, []int{1, 4}) {
		fmt.Println(index)
		count++
		if count == 3 {
			break
		}
	}
	// Output:
	// 0
	// 2
	// 144
}

func ExampleFindFirst() {

	// sqrt(3) = 0.1732050807... * 10^1
	n := sqroot.Sqrt(3)

	fmt.Println(sqroot.FindFirst(n, []int{0, 5, 0, 8}))
	// Output:
	// 4
}

func ExampleFindLast() {
	n := sqroot.Sqrt(2)
	fmt.Println(sqroot.FindLast(n.WithEnd(1000), []int{1, 4}))
	// Output:
	// 945
}

func ExampleBackwardMatches() {
	n := sqroot.Sqrt(2)
	count := 0
	iterator := sqroot.BackwardMatches(n.WithEnd(1000), []int{1, 4})
	for index := range iterator {
		fmt.Println(index)
		count++
		if count == 3 {
			break
		}
	}
	// Output:
	// 945
	// 916
	// 631
}

func ExampleFiniteNumber_Exponent() {

	// sqrt(50176) = 0.224 * 10^3
	n := sqroot.Sqrt(50176)

	fmt.Println(n.Exponent())
	// Output:
	// 3
}

func ExampleFiniteNumber_All() {

	// sqrt(7) = 0.26457513110... * 10^1
	n := sqroot.Sqrt(7)

	for index, value := range n.All() {
		fmt.Println(index, value)
		if index == 5 {
			break
		}
	}
	// Output:
	// 0 2
	// 1 6
	// 2 4
	// 3 5
	// 4 7
	// 5 5
}

func ExampleFiniteNumber_Values() {
	// sqrt(7) = 0.26457513110... * 10^1
	n := sqroot.Sqrt(7)

	for value := range n.WithEnd(6).Values() {
		fmt.Println(value)
	}
	// Output:
	// 2
	// 6
	// 4
	// 5
	// 7
	// 5
}

func ExampleFiniteNumber_Backward() {

	// sqrt(7) = 0.26457513110... * 10^1
	n := sqroot.Sqrt(7).WithSignificant(6)

	for index, value := range n.Backward() {
		fmt.Println(index, value)
	}
	// Output:
	// 5 5
	// 4 7
	// 3 5
	// 2 4
	// 1 6
	// 0 2
}

func ExampleFiniteNumber_WithStart() {

	// sqrt(29) = 5.3851648...
	n := sqroot.Sqrt(29)

	// Find all occurrences of '85' in the first 1000 digits of sqrt(29)
	matches := sqroot.Matches(n.WithEnd(1000), []int{8, 5})
	fmt.Println(slices.Collect(matches))

	// Find all occurrences of '85' in the first 1000 digits of sqrt(29)
	// on or after position 800
	matches = sqroot.Matches(n.WithStart(800).WithEnd(1000), []int{8, 5})
	fmt.Println(slices.Collect(matches))
	// Output:
	// [2 167 444 507 511 767 853 917 935 958]
	// [853 917 935 958]
}

func ExampleFiniteNumber_Exact() {
	n := sqroot.Sqrt(2).WithSignificant(60)
	fmt.Println(n.Exact())
	// Output:
	// 1.41421356237309504880168872420969807856967187537694807317667
}

func ExampleWrite() {
	n := sqroot.Sqrt(2)
	sqroot.Write(n.WithEnd(1000))
	// Output:
	//   0  14142 13562 37309 50488 01688 72420 96980 78569 67187 53769
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

func ExamplePrint() {

	// Find the square root of 2.
	n := sqroot.Sqrt(2)

	fmt.Printf("10^%d *\n", n.Exponent())
	sqroot.Print(n, sqroot.UpTo(1000))
	fmt.Println()
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

func ExamplePrint_positions() {
	var pb sqroot.PositionsBuilder
	sqroot.Print(
		sqroot.Sqrt(2),
		pb.AddRange(110, 120).AddRange(200, 220).Build(),
		sqroot.DigitsPerRow(20),
	)
	fmt.Println()
	// Output:
	// 100  ..... ..... 30912 29702
	// 200  70109 55997 16059 70274
}

func ExampleFiniteNumber_At() {

	// sqrt(7) = 0.264575131106459...*10^1
	n := sqroot.Sqrt(7)

	fmt.Println(n.At(0))
	fmt.Println(n.At(1))
	fmt.Println(n.At(2))
	// Output:
	// 2
	// 6
	// 4
}

func ExampleFiniteNumber_WithSignificant() {

	// n is 1.42857142857... but truncated to 10000 significant digits
	n := sqroot.SqrtRat(100, 49).WithSignificant(10000)

	// Instead of running forever, FindFirst returns -1 because n is truncated.
	fmt.Println(sqroot.FindFirst(n, []int{1, 1, 2}))
	// Output:
	// -1
}

func ExamplePositions() {
	var builder sqroot.PositionsBuilder
	builder.AddRange(0, 7).AddRange(40, 50)
	positions := builder.AddRange(5, 10).Build()
	fmt.Printf("End: %d\n", positions.End())
	for pr := range positions.All() {
		fmt.Printf("%+v\n", pr)
	}
	// Output:
	// End: 50
	// {Start:0 End:10}
	// {Start:40 End:50}
}
