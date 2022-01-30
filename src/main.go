package main

import (
	log "github.com/sirupsen/logrus"
)

const (
	START_SPREAD                    = 0.5 // as % of mark price
	START_ORDER_SIZE                = 1
	WAIT_TIME_BEFORE_REQUOTING_ONE  = 40 // seconds
	WAIT_TIME_BEFORE_REQUOTING_BOTH = 40 // seconds
)

func main() {

	log.SetLevel(log.InfoLevel)

	simulator := NewSimulator("../datasets/DOGE_1s.csv")
	simulator.Start()
	simulator.PrintResults()
	simulator.WriteResults("../results/doge_results.csv")
}
