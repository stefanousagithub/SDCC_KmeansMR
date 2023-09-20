package main

import (
	utils "kmeansMR/cluster"
	"log"
	"net"
	"net/rpc"
)

type API int

var Cc utils.Clusters
var k int

// func Reducer(cc clusters.Clusters, numPoints int, deltaThreshold float64, chInRed chan InRed, chOutRed chan clusters.Clusters, chMastRed chan bool) {
func (a *API) Reducer(input utils.InRed, reply *utils.OutRed) error {
	// Init temporary struct
	arrSums := make([]utils.Coordinates, k)
	weights := make([]int, k)
	for i := 0; i < len(arrSums); i++ {
		arrSums[i] = utils.Coordinates{0, 0}
	}

	// Join all partial sums
	for _, kv := range input {
		arrSums[kv.Center].Sum(kv.SumObs)
		weights[kv.Center] += kv.Npoint
	}

	// Divide each center for the number of points
	for i := 0; i < k; i++ {
		// If cluster is empty. Steal a point to another cluster
		if weights[i] == 0 {
			ri := 0
			for ; ri < len(weights); ri++ {
				// find a cluster with at least two data points, otherwise
				// we're just emptying one cluster to fill another
				if weights[ri] > 1 {
					break
				}
			}
			utils.MoveObs(&arrSums[i], &arrSums[ri], float64(weights[ri]))
			weights[i] += 1
			weights[ri] -= 1
		}
		arrSums[i].Divide(float64(weights[i]))
	}

	// Set inside the cluster struct
	Cc.SetCc(arrSums)

	*reply = utils.OutRed(Cc) // Send cluster result
	return nil
}

// Set initial cluster
func (a *API) InitReducer(input utils.Clusters, reply *utils.OutRed) error {
	Cc = input
	k = len(Cc)
	return nil
}

func main() {
	// Open Rpc connection
	api := new(API)
	server := rpc.NewServer()
	err := server.RegisterName("API", api)
	if err != nil {
		log.Fatal("error registering API", err)
	}

	// Rpc lister
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal("Listener error", err)
	}

	log.Printf("Reducer is listening")
	server.Accept(listener)
}
