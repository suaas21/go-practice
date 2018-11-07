package main

import "fmt"

// func main(){

//   fmt.Println(len("Hello World"))
//   fmt.Println("Hello World"[1])
//   fmt.Println("Hello " + "World")
//   fmt.Println((true && false) || (false && true) || !(false && false))
//   var x string
//   x = "hello "
//   x =x + "world"
//   fmt.Println(x)
//   var y = "hi"
//   fmt.Println(y)
//   var p int = 10
//   fmt.Println(p)

// }

func main() {

	// x := make(map[string]int)
	// x["a"] = 10
	// x["b"] = 11
	// x["c"] = 10
	// x["d"] = 11
	// fmt.Println(len(x))
	// fmt.Println(x)
	// delete(x, "a")
	// fmt.Println(len(x))
	// fmt.Println(x)
	// fmt.Println(x["k"])
	// name, ok := x["k"]
	// fmt.Println(name, ok)
	// for i := 0; i < len(x); i++ {
	// 	fmt.Print(x[""])
	// }
	// if name, ok := x["k"]; ok {
	// 	fmt.Println(name, ok)
	// }
	// elements := map[string]map[string]int16{
	// 	"a": map[string]int16{
	// 		"aa": 1,
	// 		"bb": 2,
	// 	},
	// }
	// if name, ok := elements["b"]; ok {
	// 	fmt.Println(name["aa"], name["bb"])
	//}
	// x := []int16{
	// 	48, 96, 86, 68,
	// 	57, 82, 63, 70,
	// 	37, 34, 83, 27,
	// 	19, 97, 9, 17,
	// }
	// fmt.Println(x)
	// var min int16
	// min = 1 << 8
	// fmt.Println(min)
	// for i := 0; i < len(x); i++ {
	// 	if x[i] < min {
	// 		min = x[i]
	// 	}
	// }
	//fmt.Println(min)
	//*********function
	// var xs []float64
	// xs = make([]float64, 5)
	// for i := 0; i < len(xs); i++ {

	// 	fmt.Scanf("%f", &xs[i])
	// }
	// // xs := []float64{
	// // 	1.1,
	// // 	2.3,
	// // 	3.3,
	// // }
	// for _, value := range xs {
	// 	fmt.Println(value)
	// }
	// //xs := []float64{98, 93, 77, 82, 83}
	// fmt.Println(average(xs))
	//****************
	second()
	defer first()
	third()

}

// func average(arr []float64) float64 {
// 	total := 0.0
// 	for _, value := range arr {
// 		total += value
// 	}
// 	return total / float64(len(arr))
// }
func first() {
	fmt.Println("1st")
}
func second() {
	fmt.Println("2nd")
}
func third() {
	fmt.Println("3rd")
}
