package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run genPoints.go [nPoints] [name]\n")
	}
	nPoints, _ := strconv.Atoi(os.Args[1])
	name := os.Args[2]

	f, err := os.Create("./" + name + ".txt")
	check(err)
	datawriter := bufio.NewWriter(f)

	defer f.Close()
	defer datawriter.Flush()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < nPoints; i++ {
		n1 := fmt.Sprintf("%f", rand.Float64()*100)
		n2 := fmt.Sprintf("%f", rand.Float64()*100)
		datawriter.WriteString(n1 + " " + n2 + "\n")
	}

}
