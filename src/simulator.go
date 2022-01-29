package main

type Simulator struct {
	exchange    *Exchange
	marketMaker *MarketMaker
	tickData    *TickData
}

func NewSimulator(dataPath string) *Simulator {

	exchange := NewExchange()
	marketMaker := NewMarketMaker()

	// link market maker and exchange through callbacks
	marketMaker.SetExchangeAPI(exchange.GetAPI())
	exchange.NotifyTickUpdate = marketMaker.HandleTickUpdate

	return &Simulator{
		exchange:    exchange,
		marketMaker: marketMaker,
		tickData:    NewSymbolDataFromProcessedFile(dataPath),
	}
}

func (s *Simulator) Start() {

	for _, tick := range s.tickData.Data {
		s.exchange.Next(tick)
	}
}
