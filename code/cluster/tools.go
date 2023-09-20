package utils

type KeyValue struct { // Main struct for MapReduce comunications
	Center int         // Center index
	SumObs Coordinates // Partial sum of center
	Npoint int         // Num of points concerned in the sum
}

func NewKV(center int, SumObs Coordinates, Npoint int) KeyValue {
	return KeyValue{center, SumObs, Npoint}
}

type Input struct { // Master input values
	K    int    // Number clusters
	File string // Path file with points
}
type Output struct { // Master output values
	Cc      Clusters
	NPoints int
	NumIters int
}
type InMap Clusters  // Mapper input value
type OutMap struct { // Mapper output value
	Kvs     []KeyValue
	Changes int
}
type InRed []KeyValue // Reducer input value
type OutRed Clusters  // Reducer output value

type TestInput struct { // Master input values for testing: Change kmeans configurations
	NumMap   int     // Number mappers
	MaxIter  int     // Number max iterations
	ThrShold float64 // Number threshold
}
