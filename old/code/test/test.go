package main

import (
	"fmt"
	utils "kmeansMR/cluster"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

var PORT string = "8001"
var NMAPS [8]int = [8]int{1, 2, 3, 4, 5, 6, 7, 8}
var THRSHOLD float64 = 0.01
var MAXITER int = 50

var KS [3]int = [3]int{10, 50, 100}
var PATHS [1]string = [1]string{"./points/rand10000.txt"}

type result struct {
	path string
	k    int
	nmap int
	time time.Duration
}

func main() {
	// Analize input data
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run test.go [IP_SERVER]\n")
	}
	ip := os.Args[1]
	client, err := rpc.Dial("tcp", ip+":"+PORT)
	defer client.Close()
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	var out1 int
	var out2 utils.Output
	// default test

	var results []result

	printBefore()

	
	// Loop for path
	for _, p := range PATHS {
		fmt.Println("Path: ", p)
		// loop for k
		for _, k := range KS {
			fmt.Println("-- Number k: ", k)
			// loop for nmap
			for _, nmap := range NMAPS {
				err = client.Call("API.SetKmeans", utils.TestInput{NumMap: nmap, MaxIter: MAXITER, ThrShold: THRSHOLD}, &out1)
				if err != nil || out1 != 0 {
					log.Fatal("[PATHS: ", p, " nmap: ", nmap, " ks: ", k, "] Error in API.SetKmeans: ", err)
				}
				start := time.Now()
				err = client.Call("API.MapReduce", utils.Input{K: k, File: p}, &out2)
				if err != nil {
					log.Fatal("[PATHS: ", p, " nmap: ", nmap, " ks: ", k, "] Error in API.MapReduce: ", err)
				}
				elapsed := time.Since(start)
				results = append(results, result{p, k, nmap, elapsed})
				fmt.Println("---- nmap = ", nmap, " -> time: ", elapsed)
			}
			fmt.Println("")
		}
	}

	graphResult(PATHS[0], results)
}

func printBefore(){
	fmt.Println("*****************************************************************")
	fmt.Println("*****************************************************************")
	fmt.Println("*****************                               *****************")
	fmt.Println("*****************             Test              *****************")
	fmt.Println("*****************                               *****************")
	fmt.Println("*****************************************************************")
	fmt.Println("*****************************************************************")
	fmt.Println("")
	fmt.Println("\ndefault config: THRSHOLD = ", THRSHOLD, ", MAXITER = ", MAXITER)
	fmt.Println()
}

func graphResult(path string, results []result) {
	testChart := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	testChart.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
		charts.WithTitleOpts(opts.Title{Title: "Test result", Subtitle: "10 000 random points. X,Y = [num.mappers, time(ms)]. K = 20(Y), 10(B), 5(R)"}))

	// Put data into instance
	for _, k := range KS {
		items := make([]opts.LineData, 0)
		for _, res := range results {
			if path == res.path && k == res.k {
				items = append(items, opts.LineData{Value: res.time / 1000000})

			}
		}
		testChart.SetXAxis([]string{"1", "2", "3", "4", "5", "6", "7", "8"}).
			AddSeries("Category "+strconv.Itoa(k), items, charts.WithLabelOpts(opts.Label{
				Show:      true,
				Formatter: "{c}",
			}))
	}

	f, _ := os.Create("line.html")
	_ = testChart.Render(f)
}
