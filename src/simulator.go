package main

type Simulator struct {
	marketMaker *MarketMaker
	symbolData  *SymbolData
}

func NewSimulator() *Simulator {
	return &Simulator{}
}

func (s *Simulator) Start(dataPath string, strategy *Strategy) {

	s.symbolData = NewSymbolDataFromProcessedFile(dataPath)
	s.marketMaker = NewMarketMaker(strategy)

	for _, symbolDataItem := range s.symbolData.Data {
		s.marketMaker.ProcessNewData(symbolDataItem)
	}
}
