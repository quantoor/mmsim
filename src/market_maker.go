package main

type Status int

const (
	SPREAD Status = iota
	QUOTING
	WAITING
)

type MarketMaker struct {
	timestamp    int64
	positionSize float64
	openOrders   map[string]LimitOrder
	filledOrders map[string]LimitOrder
	pnl          float64
	status       Status
	exchangeAPI  *ExchangeAPI
}

func NewMarketMaker() *MarketMaker {
	return &MarketMaker{
		timestamp:    0,
		positionSize: 0,
		openOrders:   make(map[string]LimitOrder),
		filledOrders: make(map[string]LimitOrder),
		pnl:          0,
		status:       SPREAD,
	}
}

func (m *MarketMaker) SetExchangeAPI(api *ExchangeAPI) {
	m.exchangeAPI = api
}

func (m *MarketMaker) HandleTickUpdate(tick Tick) {
	m.timestamp = tick.Timestamp

	switch m.status {

	case SPREAD:
		//

	case QUOTING:
		//

	case WAITING:
		//
	}
}
