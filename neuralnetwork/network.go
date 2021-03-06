package neuralnetwork

import (
	"errors"
	"fmt"
)

var nilNetwork = &Network{}
var nilHiddenLayers = []*layer{}

var errNoLayout = errors.New("Error: please input layout information for network layers")

//Network has an input and output layer, and
//a variadic number of hidden layers. The network
//forward-feeds inputs, and also regressively
//feeds backwards to train the weights and biases
type Network struct {
	Config       *Config
	InputLayer   *layer
	HiddenLayers []*layer
	OutputLayer  *layer
	CostFunction costfunction //just like we're doing with activation functions in layers
}

//NewNetwork takes a config struct and a variadic number of inputs
//representing layers. The minimum number of layers is two. Eventually
//this will have more nuanced options for things like activation type
func NewNetwork(config *Config, layout ...int) (*Network, error) {

	res := Network{
		Config: config,
	}

	res.CostFunction = getCost(config.CostFunction)

	if len(layout) < 1 {
		return nilNetwork, errNoLayout
	} else if len(layout) < 2 {
		res.InputLayer = res.newLayer(layout[0], layout[0], config.DefaultActivationType)
		res.OutputLayer = res.newLayer(layout[0], layout[0], config.OutputActivationType)
	} else {
		res.InputLayer = res.newLayer(layout[0], layout[1], config.DefaultActivationType)
		if len(layout) > 2 {
			res.HiddenLayers = make([]*layer, len(layout)-2)
		}
		for i := 0; i < len(layout)-2; i++ { //for all the layers minus the first and last
			res.HiddenLayers[i] = res.newLayer(layout[i+1], layout[i+2], config.DefaultActivationType) //if this isn't clear, we can change the way it's indexed
		}
		//Outputlayer in this case is mostly just a transformation layer, so it will always have
		//the same number of inputs and outputs
		res.OutputLayer = res.newLayer(layout[len(layout)-1], layout[len(layout)-1], config.OutputActivationType)
	}
	return &res, nil
}

//ForwardFeed takes an input slice and feeds it through the network
//to produce a result in the output layer
func (n *Network) ForwardFeed(input []float64) []float64 {
	in := n.InputLayer.fire(input)
	for i := 0; i < len(n.HiddenLayers); i++ {
		in = n.HiddenLayers[i].fire(in)
	}
	return n.OutputLayer.fire(in)
}

//Backpropagate will take a slice of expected values and compare
//it to the network's output to calculate the cost using the
//config's cost function. It will then send this cost up through
//each layer using the backpropagation algorithm
func (n *Network) Backpropagate(prime []float64) {
	prime = n.OutputLayer.stepBack(n.Config.LearningRate, prime)
	for i := len(n.HiddenLayers) - 1; i >= 0; i-- {
		prime = n.HiddenLayers[i].stepBack(n.Config.LearningRate, prime)
	}
	n.InputLayer.stepBack(n.Config.LearningRate, prime)
}

//String is a stringer
func (n *Network) String() string {

	var s string

	s += fmt.Sprintf("Neural Network:,\n\nInput Layer: %v neurons\nOutput Layer: %v neurons\nHidden Layers: %v\n\n", len(n.InputLayer.Inputs), len(n.OutputLayer.Outputs), len(n.HiddenLayers))

	s += "\nInput Layer\n"
	s += n.InputLayer.String(false)

	for i := 0; i < len(n.HiddenLayers); i++ {
		s += fmt.Sprintf("\nLayer %v:\n", i)
		s += n.HiddenLayers[i].String()
	}

	s += "\nOutput Layer\n"
	s += n.OutputLayer.String()

	s += fmt.Sprintf("\nNetwork input [%v inputs]:\n", len(n.InputLayer.Inputs))

	for i := 0; i < len(n.InputLayer.Inputs); i++ {
		s += fmt.Sprintf("%1.4f, ", n.InputLayer.Inputs[i])
	}

	s += fmt.Sprintf("\nNetwork output:\n")

	for i := 0; i < len(n.OutputLayer.Outputs); i++ {
		s += fmt.Sprintf("%1.4f, ", n.OutputLayer.Outputs[i])
	}

	s += "\n\n"

	return s

}
