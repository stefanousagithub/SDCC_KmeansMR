package main

import (
	"fmt"
	"image/color"
	utils "kmeansMR/cluster"
	"log"
	"net/rpc"
	"os"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

var PORTSERVER string = "5000"

func main() {
	// Analize input data
	if len(os.Args) != 3  {
		fmt.Printf("Usage: go run client.go [k] [fileObj] \n")
		os.Exit(1)
	}

	k, e := strconv.Atoi(os.Args[1])
	file := os.Args[2]
	if e != nil {
		fmt.Printf("k values is not correct!")
		os.Exit(1)
	}

	printBefore(k, file)

	// Open Rpc connection
	client, err := rpc.Dial("tcp", ":"+PORTSERVER)
	defer client.Close()
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	var out utils.Output
	// Call function Partition
	err = client.Call("API.MapReduce", utils.Input{K: k, File: file}, &out)
	if err != nil {
		log.Fatal("Error in API.MapReduce: ", err)
	}

	// print results console
	printAfter(out.Cc, out.NPoints ,out.NumIters)

	// plot chart centers
	plotKmeans(out.Cc, k, out.NPoints)
}

func printBefore(k int, file string){
	fmt.Println("*****************************************************************")
	fmt.Println("*****************************************************************")
	fmt.Println("***********************                      ********************")
	fmt.Println("***********************        Client        ********************")
	fmt.Println("***********************                      ********************")
	fmt.Println("*****************************************************************")
	fmt.Println("*****************************************************************")

	fmt.Println("")
	fmt.Println("k selected:", k, " path:", file)
	fmt.Println("----------------------  Start execution   -----------------------")
	fmt.Println("")
}

func printAfter(cc utils.Clusters, nPoints int, numIters int){
	fmt.Println("dataset analized")
	fmt.Println("-- number points: ", nPoints, ", number iterations: ", numIters)
	fmt.Println("")

	// Analize output data
	for i, c := range cc {
		fmt.Printf("Cluster: %d\n", i)
		fmt.Printf("-> Centered at x: %.2f y: %.2f\n", c.Center[0], c.Center[1])
	}
	fmt.Println("")
}

func plotKmeans(cc utils.Clusters, k int, nPoints int) {
	/** Plot centers in Scatter plot **/
	p := plot.New()

	xysCenters := make(plotter.XYs, k)
	for i, c := range cc {
		xysCenters[i].X = c.Center[0]
		xysCenters[i].Y = c.Center[1]
	}

	s, err := plotter.NewScatter(xysCenters)

	if err != nil {
		log.Fatalf("could not create scatter: %v", err)
		return
	}
	s.GlyphStyle.Shape = draw.CrossGlyph{}
	s.Color = color.RGBA{R: 255, A: 255}
	p.Add(s)

	wt, err := p.WriterTo(500, 500, "png")
	if err != nil {
		log.Fatalf("Could not create writer: %v", err)
		return
	}

	f, err := os.Create("out.png")
	if err != nil {
		log.Fatalf("could not create out.png: %v", err)
		return
	}
	_, err = wt.WriteTo(f)
	if err != nil {
		log.Fatalf("could not write to out.png: %v", err)
		return
	}

	if err := f.Close(); err != nil {
		log.Fatalf("could not close out.png: %v", err)
		return
	}
}
