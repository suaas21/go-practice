package main

import (
	"fmt"
	"math"
)

type rectangle struct {
	x1, y1, x2, y2 float64
}
type circle struct {
	r float64
}

func (rec rectangle) area() float64 {
	a := rec.x2 - rec.x1
	b := rec.y2 - rec.y1
	return math.Sqrt(a*a + b*b)
}
func (cir circle) area() float64 {
	return math.Pi * cir.r * cir.r
}
func main() {
	rec := rectangle{x1: 2, y1: 2, x2: 3, y2: 3}
	fmt.Println(rec.area())
	cir := circle{r: 3}
	fmt.Println(cir.area())

}
