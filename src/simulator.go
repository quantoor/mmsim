package main

import "fmt"

type Simulator struct {
	exchange    *Exchange
	marketMaker *MarketMaker
	tickData    *TickData
}

func NewSimulator(dataPath string) *Simulator {

	exchange := NewExchange()
	marketMaker := NewMarketMaker(exchange.GetAPI())

	return &Simulator{
		exchange:    exchange,
		marketMaker: marketMaker,
		tickData:    NewSymbolDataFromProcessedFile(dataPath),
	}
}

func (s *Simulator) Start() {

	for _, tick := range s.tickData.Data {

		s.exchange.ProcessNextTick(tick)
		s.marketMaker.ProcessNextTick(tick)
	}
}

func (s *Simulator) GetResults() {
	fmt.Println("Position size", s.exchange.position.Size)
	fmt.Println(s.exchange.getFilledOrders())
}
