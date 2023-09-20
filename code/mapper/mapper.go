package main

import (
	utils "kmeansMR/cluster"
	"log"
	"net"
	"net/rpc"
)

type API int

var obs utils.Observations
var points []int

// obs clusters.Observations, chInMap chan clusters.Clusters, chInRed chan InRed
func (a *API) Mapper(input utils.InMap, reply *utils.OutMap) error {
	// Init struct
	cc := utils.Clusters(input)
	kvs := make([]utils.KeyValue, len(cc))
	for i := 0; i < len(kvs); i++ {
		kvs[i].SumObs = utils.Coordinates{0, 0}
	}

	changes := 0
	// Loop for number of observations. Execute both mapper and combiner work
	for p, point := range obs {
		ci := cc.Nearest(point)
		kvs[ci].Center = ci
		kvs[ci].SumObs.Sum(point)
		kvs[ci].Npoint += 1
		if points[p] != ci {
			points[p] = ci
			changes++
		}
	}

	// 2) Return num changes
	*reply = utils.OutMap{Kvs: kvs, Changes: changes}
	return nil
}

// Init the observations for the mapper
func (a *API) InitObs(input []utils.Coordinates, reply *utils.OutMap) error {
	obs = input
	points = make([]int, len(obs))
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

	// Listen on port 8000
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal("Listener error", err)
	}

	log.Printf("Mapper is listening")
	server.Accept(listener)
}
