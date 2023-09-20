package utils

type KeyValue struct {
	Center int
	Obs    Coordinates
}

func NewKV(center int, obs Coordinates) KeyValue {
	return KeyValue{center, obs}
}

type Input struct { // Master input values
	K      int    // Number clusters
	NumMap int    // Number mappers
	File   string // Path file with points
}
type Output struct { // Master output values
	Cc      Clusters
	NPoints int
}
type InMap Clusters  // Mapper input value
type OutMap struct { // Mapper output value
	Kvs     []KeyValue
	Changes int
}
type InRed []KeyValue // Reducer input value
type OutRed Clusters  // Reducer output value

// Variant configurations with combiners
type KeyValueComb struct {
	Center int
	SumObs Coordinates
	Npoint int
}

func NewKVComb(center int, sumObs Coordinates, npoint int) KeyValueComb {
	return KeyValueComb{center, sumObs, npoint}
}

type OutMapComb struct {
	Kvs     []KeyValueComb
	Changes int
}

type InRedComb []KeyValueComb

type TestInput struct { // Master input values for testing: Change kmeans configurations
	NumMap   int     // Number mappers
	MaxIter  int     // Number max iterations
	ThrShold float64 // Number threshold
}
