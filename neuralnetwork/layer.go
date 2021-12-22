package neuralnetwork

import (
	"math/rand"
	"time"
)

//Layer is a two-dimensional array of neurons. I guess. How does this work?
type Layer struct {
	Neurons []*Neuron
	Weights [][]float64 //Matrix of connectins to next layer
}

//NewLayer takes two inputs, nsize and wsize, where nsize is the number of neurons
//in the layer, and wsize is the number of weights for each neuron, aka the number
//of neurons in the next layer. The function returns a new layer
func NewLayer(nsize, wsize int) *Layer {

	//activation type should be handled here, not at the neuron level

	rand.Seed(time.Now().UnixNano())

	n := make([]*Neuron, nsize)
	for i := 0; i < nsize; i++ {
		n[i] = NewNeuron()
	}

	//each row of the weights matrix will represent a neuron in the next layer
	w := make([][]float64, wsize)
	for i := 0; i < wsize; i++ {
		//and each column in each row will represent a neuron in the current layer's
		//relationship to the output neuron
		w[i] = make([]float64, nsize)
		for j := 0; j < nsize; j++ {
			//rando-fill those weights
			w[i][j] = rand.Float64()
		}
	}

	return &Layer{
		Neurons: n,
		Weights: w,
	}

}

//Activate does a nasty matrix multiplication thing
//Evemtially I'll bring in an optimized matrix library,
//but for now I'll hand-roll it
func (l *Layer) Activate(inputs []float64) []float64 {
	res := make([]float64, len(inputs))
	//Send each input through its respective neuron
	for i := 0; i < len(inputs); i++ {
		a := l.Neurons[i].Fire(inputs[i])
		//Feed the output from each neuron through the weights matrix and =+ the
		//dot product to the jth index of res
		for j := 0; j < len(res); j++ {
			res[j] += l.Weights[i][j] * a
		}
	}
	//return that vector
	return res
}

//Descend does gradient descent at a layer level
func (l *Layer) Descend() {}
