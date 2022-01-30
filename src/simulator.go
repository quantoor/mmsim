package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Simulator struct {
	exchange       *Exchange
	marketMaker    *MarketMaker
	tickData       *TickData
	resultsHandler *ResultsHandler
}

func NewSimulator(dataPath string) *Simulator {

	exchange := NewExchange()
	marketMaker := NewMarketMaker(exchange.GetAPI())

	return &Simulator{
		exchange:       exchange,
		marketMaker:    marketMaker,
		tickData:       NewSymbolDataFromProcessedFile(dataPath),
		resultsHandler: NewResultsHandler(),
	}
}

func (s *Simulator) Start() {

	log.Infof("Start simulation from %s to %s", s.tickData.StartDate, s.tickData.EndDate)

	for _, tick := range s.tickData.Data {

		s.exchange.ProcessNextTick(tick)
		s.marketMaker.ProcessNextTick(tick)

		s.addResults()
	}
}

func (s *Simulator) PrintResults() {
	fmt.Println("Position: ", s.exchange.getPosition().Size)
	fmt.Println("Balance: ", s.exchange.getBalance())
	fmt.Println("Executed orders: ", len(s.exchange.getFilledOrders()))
}

func (s *Simulator) addResults() {
	bids := make([]float64, 0)
	asks := make([]float64, 0)
	for _, o := range s.exchange.openOrders {
		if o.IsBuy {
			bids = append(bids, o.Price)
		} else {
			asks = append(asks, o.Price)
		}
	}

	data := ResultsData{
		Timestamp: s.exchange.timestamp,
		Price:     s.exchange.markPrice,
		Bids:      bids,
		Asks:      asks,
		Balance:   s.exchange.balance,
		Position:  s.exchange.position.Size,
	}
	s.resultsHandler.Data = append(s.resultsHandler.Data, data)
}

func (s *Simulator) WriteResults(path string) {
	s.resultsHandler.writeToFile(path)
}
