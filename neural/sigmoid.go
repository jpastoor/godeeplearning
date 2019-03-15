package neural

import "math"

// Sigmoid neuron has a sigmoid activation function
type Sigmoid struct{}

func (a *Sigmoid) ActivateFloat(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func (a *Sigmoid) ActivateFloatVector(z []float64) []float64 {
	output := make([]float64, len(z))
	for i := 0; i < len(z); i++ {
		output[i] = a.ActivateFloat(z[i])
	}

	return output
}

func (a *Sigmoid) DActivateDSum(x, output float64) float64 {
	return output * (1 - output)
}

// Linear neuron has a linear activation function
type Linear struct{}

func (a *Linear) Activate(sum float64) float64 {
	return sum
}
func (a *Linear) DActivateDSum(sum, output float64) float64 {
	return 1.0
}

// Tanh

type Tanh struct{}

func (a *Tanh) Activate(sum float64) float64 {
	return 1.7159 * math.Tanh(2.0/3.0*sum)
}
