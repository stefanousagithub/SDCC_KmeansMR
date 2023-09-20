package main

var PORT string = "5000"
var NMAP [5]int = [5]int{2, 5, 10, 15, 20}
var THRSHOLD [3]int = [3]int{1, 5, 10}
var MAXITER [3]int = [3]int{20, 100, 500}
var K [5]int = [5]int{5, 10, 20, 50, 100}
var PATHS [2]string = [2]string{"../points/rand1000.txt", "../points/rand10000.txt"}

func main() {

}

/**
	out, err := exec.Command("../start.sh 32 ./points/rand10000.txt 5000").Output()

	// Open Rpc connection
	client, err := rpc.Dial("tcp", ":"+PORT)
	defer client.Close()
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
start := time.Now()
var out Output
// Call function Partition
err = client.Call("API.MapReduce", Input{k, file}, &out)
if err != nil {
	log.Fatal("Error in API.MapReduce: ", err)
}
elapsed := time.Since(start)
**/
