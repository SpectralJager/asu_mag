package main

import "fmt"

type Perceptron struct {
	weights      []float64
	threshold    float64
	learningRage float64
}

func NewPerceptron(numOfInputs int, threshold float64, leaningRate float64) Perceptron {
	return Perceptron{
		weights:      make([]float64, numOfInputs),
		threshold:    threshold,
		learningRage: leaningRate,
	}
}

func (per Perceptron) Predict(inputs []float64) float64 {
	acc := 0.0
	if len(inputs) != len(per.weights) {
		panic("mismatched number of inputs and weights")
	}
	for i, input := range inputs {
		acc += input * per.weights[i]
	}

	if acc < per.threshold {
		return 0.0
	}
	return 1.0
}

func (per Perceptron) Update(updates []float64) {
	if len(updates) != len(per.weights) {
		panic("mismatched number of updates and weights")
	}
	for i := range per.weights {
		per.weights[i] += updates[i]
	}
}

func (per Perceptron) Train(inputs [][]float64, targets []float64) {
	for i, input := range inputs {
		out := per.Predict(input)
		err := targets[i] - out
		// if err == 0 {
		// 	continue
		// }
		updates := make([]float64, len(input))
		for j := range input {
			updates[j] = per.learningRage * err * input[j]
		}
		per.Update(updates)
	}
}

func main() {
	perceptron := NewPerceptron(4, 0.5, 1)
	for range 10 {
		perceptron.Train(
			[][]float64{
				{0, 0, 0, 0},
				{0, 0, 0, 1},
				{0, 0, 1, 0},
				{0, 0, 1, 1},
				{0, 1, 0, 0},
				{0, 1, 0, 1},
				{0, 1, 1, 0},
				{0, 1, 1, 1},
				{1, 0, 0, 0},
				{1, 0, 0, 1},
			},
			[]float64{0, 1, 0, 1, 0, 1, 0, 1, 0, 1},
		)
	}
	fmt.Printf("%+v\n", perceptron)
	fmt.Println(perceptron.Predict([]float64{1, 1, 0, 1}))
}
