package neural

/**
def forward(self):
        raise NotImplementedError

    def get_forward_input(self):
        if self.previous != None:
            return self.previous.output_data
        else:
            return self.input_data

    def backward(self):
        raise NotImplementedError

    def get_backward_input(self):
        if self.next != None:
            return self.next.output_delta
        else:
            return self.input_delta

    def clear_deltas(self):
        pass

    def update_params(self, learning_rate):
        pass

    def describe(self):
        raise NotImplementedError
 */

type Layer interface {
	Forward()
	ForwardInput() []float32
	ForwardOutput() []float32
	Backward()
	BackwardInput() []float32
	BackwardOutput() []float32
	ClearDeltas()
	UpdateParams(learningRate float32)
	Describe() string
	SetNext(layer Layer)
}

type AbstractLayer struct {
	Params map[string]interface{}

	Previous Layer
	Next     Layer

	InputData  []float32
	OutputData []float32

	// For the backward pass
	InputDelta   []float32
	OutputDelta []float32
}

// Connecting layers through successors and predecessors in a sequential network
func (l AbstractLayer) Connect(layer Layer) {
	l.Previous = layer
	layer.SetNext(l)
}

func (l AbstractLayer) Forward() {

}

func (l AbstractLayer) ForwardOutput() []float32 {
	return l.OutputData
}

func (l AbstractLayer) ForwardInput() []float32 {
	// Our input is the output of the predecessor
	if l.Previous != nil {
		return l.Previous.ForwardOutput()
	}
	// Whenever we are first in line, use our input
	return l.InputData
}
func (l AbstractLayer) Backward() {

}
func (l AbstractLayer) BackwardInput() []float32 {
	if l.Next != nil {
		return l.Next.BackwardOutput()
	}

	return l.InputDelta
}
func (l AbstractLayer) ClearDeltas() {

}
func (l AbstractLayer) UpdateParams(learningRate float32) {

}
func (l AbstractLayer) Describe() string {
	return ""
}
func (l AbstractLayer) SetNext(layer Layer) {

}

func (l AbstractLayer) BackwardOutput() []float32 {
	return l.OutputDelta
}
