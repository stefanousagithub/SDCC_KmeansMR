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
	cc := utils.Clusters(input)
	// 1) Execute mapper work and save number changes (points that change cluster in currect iteration)
	kvs := make([]utils.KeyValue, 0, len(obs))
	changes := 0
	for p, point := range obs {
		ci := cc.Nearest(point)
		kvs = append(kvs, utils.NewKV(ci, point))
		if points[p] != ci {
			points[p] = ci
			changes++
		}
	}
	// 2) Return num changes
	*reply = utils.OutMap{Kvs: kvs, Changes: changes}
	return nil
}

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
