package main

type MarketMaker struct {
	timestamp    int64
	strategy     *Strategy
	positionSize float64
	openOrders   map[string]LimitOrder
	filledOrders map[string]LimitOrder
	pnl          float64
}

func NewMarketMaker(strategy *Strategy) *MarketMaker {
	return &MarketMaker{
		timestamp:    0,
		strategy:     strategy,
		positionSize: 0,
		openOrders:   make(map[string]LimitOrder),
		filledOrders: make(map[string]LimitOrder),
		pnl:          0,
	}
}

func (m *MarketMaker) ProcessNewData(data SymbolDataItem) {
	m.timestamp = data.Timestamp

}
