package main

import (
	"fmt"
	m "golang-book/packages/math"
)

func main() {
	xs := make([]float64, 5)
	for i := 0; i < len(xs); i++ {
		fmt.Scanf("%f", &xs[i])
	}

	avg := m.Average(xs)
	fmt.Println(avg)

}
