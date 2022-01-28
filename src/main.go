package main

func main() {

	simulator := NewSimulator()

	strategy := NewStrategy()
	simulator.Start("../datasets/DOGE_1s.csv", strategy)
	// simulator.Start("../datasets/test_doge")
}
