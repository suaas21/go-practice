package math

//Finds the average of a series of numbers
func Average(xs []float64) float64 {
	total := 0.0
	for i := 0; i < len(xs); i++ {
		total += xs[i]
	}
	return total / float64(len(xs))
}
