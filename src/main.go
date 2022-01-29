package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.InfoLevel)
	simulator := NewSimulator("../datasets/DOGE_1s.csv")
	simulator.Start()
	simulator.GetResults()
}
