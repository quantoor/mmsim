package main

func main() {

	simulator := NewSimulator("../datasets/DOGE_1s.csv")
	simulator.Start()
	// simulator.Start("../datasets/test_doge")
}
