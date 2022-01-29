package main

import (
	log "github.com/sirupsen/logrus"
)

const (
	START_SPREAD                    = 0.2 // as % of mark price
	START_ORDER_SIZE                = 1
	WAIT_TIME_BEFORE_REQUOTING_ONE  = 10 // seconds
	WAIT_TIME_BEFORE_REQUOTING_BOTH = 5  // seconds
)

type MarketMaker struct {
	exchangeAPI  *ExchangeAPI
	timestamp    int64
	markPrice    float64
	positionSize float64
	spread       float64
	askSize      float64
	bidSize      float64
}

func NewMarketMaker(exchangeApi *ExchangeAPI) *MarketMaker {
	return &MarketMaker{
		exchangeAPI: exchangeApi,
		spread:      START_SPREAD,
		askSize:     START_ORDER_SIZE,
		bidSize:     START_ORDER_SIZE,
	}
}

// PUBLIC METHODS
func (m *MarketMaker) ProcessNextTick(tick Tick) {

	m.timestamp = tick.Timestamp
	m.markPrice = tick.Price

	m.positionSize = m.exchangeAPI.GetPosition().Size
	openOrders := m.exchangeAPI.GetOpenOrders()

	if len(openOrders) == 0 {
		m.quoteOrders()

	} else if len(openOrders) == 1 {
		m.handleOneOrder(openOrders[0])

	} else if len(openOrders) == 2 {
		m.handleWaiting(openOrders)
	}
}

// PRIVATE METHODS
func (m *MarketMaker) quoteOrders() {

	log.Debugf("Quote orders")

	// quote bid and ask around mark price
	askPrice := m.markPrice * (1 + m.spread/2/100)
	bidPrice := m.markPrice * (1 - m.spread/2/100)

	m.exchangeAPI.PlaceOrder(*NewLimitOrder("", askPrice, m.askSize, false))
	m.exchangeAPI.PlaceOrder(*NewLimitOrder("", bidPrice, m.bidSize, true))
}

func (m *MarketMaker) handleOneOrder(openOrder LimitOrder) {

	if (m.timestamp - openOrder.Timestamp) >= WAIT_TIME_BEFORE_REQUOTING_ONE {

		log.Debugf("Requote one order")

		price := 0.0
		size := 0.0
		if openOrder.IsBuy {
			price = m.markPrice * (1 - m.spread/2/100)
			size = m.bidSize
		} else {
			price = m.markPrice * (1 + m.spread/2/100)
			size = m.askSize
		}

		m.exchangeAPI.CancelOrder(openOrder.Id)
		m.exchangeAPI.PlaceOrder(*NewLimitOrder("", price, size, openOrder.IsBuy))
	}
}

func (m *MarketMaker) handleWaiting(openOrders []LimitOrder) {

	// find time of last inserted order
	var lastTimestamp int64 = 0
	for _, order := range openOrders {
		if order.Timestamp > lastTimestamp {
			lastTimestamp = order.Timestamp
		}
	}

	// if enough time has passed, cancel orders and quote again
	if (m.timestamp - lastTimestamp) >= WAIT_TIME_BEFORE_REQUOTING_BOTH {
		log.Debugf("Requote both orders")

		m.exchangeAPI.CancelAllOrders()
		// TODO reduce spread?
		m.quoteOrders()
	}
}
