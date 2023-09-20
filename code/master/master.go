// Package mapReduce implements the k-means clustering algorithm
package main

import (
	"bufio"
	"fmt"
	utils "kmeansMR/cluster"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

// Kmeans configuration/option struct
type Kmeans struct {
	port               int     // Master port: Default to 8000
	mappers            []int   // index of available mappers
	plotter            Plotter // when a plotter is set, Plot gets called after each iteration. NOT IMPLEMENTED
	deltaThreshold     float64 // deltaThreshold (in percent between 0.0 and 0.1) aborts processing ifless than n%
	iterationThreshold int     // iterationThreshold aborts processing when the specified amount of algorithm iterations
}

// The Plotter interface lets you implement your own plotters
type Plotter interface {
	Plot(cc utils.Clusters, iteration int) error
}

// NewWithOptions returns a Kmeans configuration struct with custom settings
func NewWithOptions(numMap int, maxIters int, deltaThreshold float64, plotter Plotter) (Kmeans, error) {
	if deltaThreshold <= 0.0 || deltaThreshold >= 1.0 {
		return Kmeans{}, fmt.Errorf("threshold is out of bounds (must be >0.0 and <1.0, in percent)")
	}
	if numMap < 0 || numMap > 99 {
		return Kmeans{}, fmt.Errorf("number mappers is out of range (must be > 0 and < 100)")
	}
	if maxIters < 0 {
		return Kmeans{}, fmt.Errorf("number maximum iterations is out of range (must be > 0)")
	}
	return Kmeans{
		port:               8000,
		mappers:            makeMappers(numMap),
		plotter:            plotter,
		deltaThreshold:     deltaThreshold,
		iterationThreshold: maxIters,
	}, nil
}

type API int

var km Kmeans
var mapName string = "code-mapper-"

// For testing: Configure kmeans struct
func (a *API) SetKmeans(input utils.TestInput, reply *int) error {
	var e error
	km, e = NewWithOptions(input.NumMap, input.MaxIter, input.ThrShold, nil)
	if e != nil {
		*reply = -1
		log.Fatal("Error in setKmeans: ", e)
		return nil
	}
	*reply = 0
	return nil
}

func makeMappers(numMap int) []int {
	a := make([]int, numMap)
	for i := range a {
		a[i] = i
	}
	return a
}

// Remove mapper from kmeans structure
func removeMapper(mapper int) {
	var pos int
	for i, val := range km.mappers {
		if val == mapper {
			pos = i
			break
		}
	}
	km.mappers = append(km.mappers[:pos], km.mappers[pos+1:]...)
}

func emptyChannels(chInMap chan utils.InMap, chOutMap chan utils.OutMap) {
	for len(chInMap) > 0 {
		<-chInMap
	}
	for len(chOutMap) > 0 {
		<-chOutMap
	}
}

// Open file and fill array of all iterations
func fillDataset(path string) []utils.Coordinates {
	var splitting []string
	var dataset []utils.Coordinates
	file, err := os.Open(path)

	defer file.Close()
	if err != nil {
		log.Fatal("Error in fillDataset: ", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		splitting = append(splitting, scanner.Text()+"\n")
	}
	for _, elem := range splitting {
		words := strings.Fields(elem)
		if len(words) != 2 {
			log.Fatal("Error in fillDataset. File format is not correct: ")
		}
		p1, e1 := strconv.ParseFloat(words[0], 64)
		p2, e2 := strconv.ParseFloat(words[1], 64)
		if e1 != nil || e2 != nil {
			log.Fatal("Error in fillDataset. File format is not correct: ")
		}
		dataset = append(dataset, utils.Coordinates{p1, p2})
	}
	return dataset
}

// Initialize mappers: Send observations and restart if occure some error
func initMappers(dataset []utils.Coordinates) {
	restart := true
	chInit := make(chan int)
	var ret int

	// Loop until all mappers execute properly
	for restart {
		restart = false
		// Stop if all mappers have failed
		numMapper := len(km.mappers)
		if numMapper < 1 {
			log.Fatal("All mappers have failed!")
		}

		// Assign partitions of dataset (observations)
		nPoints := len(dataset)
		pos, lenChunk, resChunk := 0, nPoints/numMapper, nPoints%numMapper
		pos_fin := lenChunk + resChunk
		for _, i := range km.mappers {
			go initObs(i, dataset[pos:pos_fin], chInit)
			pos = pos_fin
			pos_fin += lenChunk
		}
		// Verify correct execution
		for i := 0; i < numMapper; i++ {
			ret = <-chInit
			if ret != -1 {
				removeMapper(ret)
				restart = true
			}
		}
	}
}

// Initialize mapper's observations
func initObs(idxMap int, obs []utils.Coordinates, chInit chan int) {
	// Open RPC connection with Mapper
	mapper, err := rpc.Dial("tcp", mapName+strconv.Itoa(idxMap+1)+":8000")
	if err != nil {
		chInit <- idxMap
		return
	}
	defer mapper.Close()

	// RPC call
	err = mapper.Call("API.InitObs", obs, nil)
	if err != nil {
		chInit <- idxMap
		return
	}
	chInit <- -1
}

// Initialize clusters reducer
func initReducer(cc utils.Clusters) {
	// Open RPC connection with Reducer
	reducer, err := rpc.Dial("tcp", "reducer:8000")
	if err != nil {
		log.Fatal("Connection reducer error: ", err)
	}
	defer reducer.Close()

	// RPC call
	err = reducer.Call("API.InitReducer", cc, nil)
	if err != nil {
		log.Fatal("Error in API.InitReducer: ", err)
	}
}

// Manage reducer connection
func threadRed(cc utils.Clusters, nPoints int, deltaThreshold float64, chInRed chan utils.InRed, chOutRed chan utils.OutRed) {
	// Open RPC connection with Reducer
	reducer, err := rpc.Dial("tcp", "reducer:8000")
	if err != nil {
		log.Fatal("Connection reducer error: ", err)
	}
	defer reducer.Close()

	for {
		var retVal utils.OutRed

		// Get input by main process
		input := <-chInRed

		// RPC call
		err = reducer.Call("API.Reducer", input, &retVal)
		if err != nil {
			log.Fatal("Error in API.Reduce: ", err)
		}
		chOutRed <- retVal
	}
}

// Manage mapper connection
func threadMap(idxMap int, chInMap chan utils.InMap, chOutMap chan utils.OutMap) {
	// Open RPC connection with Mapper
	mapper, err := rpc.Dial("tcp", mapName+strconv.Itoa(idxMap+1)+":8000")
	if err != nil {
		return
	}
	defer mapper.Close()

	// Iterations
	for {
		var retVal utils.OutMap

		// Get input by main process
		input := <-chInMap

		// RPC call
		err = mapper.Call("API.Mapper", input, &retVal)
		if err != nil {
			chOutMap <- utils.OutMap{Kvs: nil, Changes: 0}
			return
		}

		// Return status operation
		chOutMap <- retVal
	}
}

// MapReduce executes the k-means algorithm on the given dataset and
// partitions it into k clusters
func (a *API) MapReduce(input utils.Input, reply *utils.Output) error {
	// 1.0) handle input values
	dataset := fillDataset(input.File)
	k := input.K
	nPoints := len(dataset)
	if k > len(dataset) {
		return fmt.Errorf("the size of the data set must at least equal k")
	}

	// 1.1) Initialize clusters, associate punti a clusters random e var changes
	cc, err := utils.New(k, dataset)
	if err != nil {
		return err
	}

	// 1.2) Create channels
	chInMap := make(chan utils.InMap)
	chOutMap := make(chan utils.OutMap)
	chInRed := make(chan utils.InRed)
	chOutRed := make(chan utils.OutRed)

	// 1.3) Start Threads: Initialize workers (initMappers and initReducer)
	// 					   and manage threads comunication (threadMap and ThreadRed)
	initMappers(dataset)
	for _, i := range km.mappers {
		go threadMap(i, chInMap, chOutMap)
	}
	initReducer(cc)
	go threadRed(cc, nPoints, km.deltaThreshold, chInRed, chOutRed)

	// Start iterations
	count := 0
	for iter := 0; ; iter++ {
		fmt.Println(km.mappers)

		// 2.0) Init variables and empty th channels
		changes := 0
		kvs := make(utils.InRed, 0, nPoints)
		numMapper := len(km.mappers)
		emptyChannels(chInMap, chOutMap)

		// 2.1) send clusters to mapper
		for i := 0; i < numMapper; i++ {
			chInMap <- utils.InMap(cc)
		}

		// 2.2) get output Mappers and concatenate in key values array
		failMapper := false
		for i := 0; i < numMapper; i++ {
			outMap := <-chOutMap
			if outMap.Kvs == nil {
				failMapper = true
				continue
			}
			kvs = append(kvs, outMap.Kvs...)
			changes += outMap.Changes
		}
		// If some mappers fail: Reinitialize mappers with range observations
		// and restart the current iteration
		if failMapper {
			initMappers(dataset)
			continue
		}

		// 2,3) Shuffle and sort: Not necessary with one reducer

		// 2.4) Send keyValues to reducer
		chInRed <- kvs

		// 2.5) wait cluster results (Channel2)
		cc = utils.Clusters(<-chOutRed)

		if cc == nil {
			return fmt.Errorf("Error reducer; %s", err)
		}
		// 2.6) Verify conditions to continue
		if iter >= km.iterationThreshold || changes == 0 || changes <= int(float64(nPoints)*km.deltaThreshold) {
			break
		}
		count++

	}
	fmt.Println("num mappers : ", len(km.mappers), ", Iterarions: ", count)

	// Return clusters and number observations
	*reply = utils.Output{Cc: cc, NPoints: nPoints, NumIters: count}
	return nil
}

func main() {
	// Analize input data
	if len(os.Args) != 4 {
		log.Fatal("Usage: go run master.go [numMappers] [maxIters] [threshold]\n")
	}
	numMappers, _ := strconv.Atoi(os.Args[1])
	maxIters, _ := strconv.Atoi(os.Args[2])
	deltaThreshold, _ := strconv.ParseFloat(os.Args[3], 64)

	if numMappers < 1 || numMappers >= 99 {
		log.Fatal("The number of mappers isn't correct or is too high!")
	}

	// Create kmeans struct: To save main configurations
	var err error
	km, err = NewWithOptions(numMappers, maxIters, deltaThreshold/1000, nil)
	if err != nil {
		log.Fatal("error kmeans parameters: ", err)
	}

	// Open RPC connection
	api := new(API)
	server := rpc.NewServer()
	err = server.RegisterName("API", api)
	if err != nil {
		log.Fatal("error registering API", err)
	}
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(km.port))
	if err != nil {
		log.Fatal("Listener error", err)
	}
	log.Printf("\nMapReduce started: Master is listening")

	server.Accept(listener)
}
