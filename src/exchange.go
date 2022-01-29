package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type ExchangeAPI struct {
	PlaceOrder      func(LimitOrder)
	CancelOrder     func(string) bool
	CancelAllOrders func()
	GetOpenOrders   func() []LimitOrder
	GetFilledOrders func() []LimitOrder
	GetMarkPrice    func() float64
	GetBalance      func() float64
	GetPosition     func() *Position
	GetTimestamp    func() int64
}

type Exchange struct {
	timestamp    int64
	markPrice    float64
	balance      float64
	position     *Position
	openOrders   map[string]LimitOrder
	filledOrders map[string]LimitOrder

	orderIdCounter int64 // used for order ID
}

func NewExchange() *Exchange {
	return &Exchange{
		position:       &Position{},
		openOrders:     make(map[string]LimitOrder),
		filledOrders:   make(map[string]LimitOrder),
		orderIdCounter: 0,
	}
}

// PUBLIC METHODS
func (e *Exchange) ProcessNextTick(tick Tick) {

	e.timestamp = tick.Timestamp
	e.markPrice = tick.Price

	// look for open orders to execute
	for _, order := range e.openOrders {

		if order.IsBuy && order.Price >= e.markPrice {
			e.executeOrder(order)

		} else if !order.IsBuy && order.Price <= e.markPrice {
			e.executeOrder(order)
		}
	}
}

func (e *Exchange) GetAPI() *ExchangeAPI {
	return &ExchangeAPI{
		PlaceOrder:      e.placeOrder,
		CancelOrder:     e.cancelOrder,
		CancelAllOrders: e.cancelAllOrders,
		GetOpenOrders:   e.getOpenOrders,
		GetFilledOrders: e.getFilledOrders,
		GetMarkPrice:    e.getMarkPrice,
		GetBalance:      e.getBalance,
		GetPosition:     e.getPosition,
		GetTimestamp:    e.getTimestamp,
	}
}

// PRIVATE METHODS
func (e *Exchange) executeOrder(order LimitOrder) {

	log.Debugf("Execute order %s", order.String())

	// update order
	order.FilledTimestamp = e.timestamp
	order.ComputeRealizedPnl(e.position)

	// update open orders
	delete(e.openOrders, order.Id)

	// update filled orders
	e.filledOrders[order.Id] = order

	// update position
	if order.IsBuy {
		e.position.Size += order.Size
	} else {
		e.position.Size -= order.Size
	}
}

func (e *Exchange) placeOrder(order LimitOrder) {

	order.Id = fmt.Sprint(e.orderIdCounter)
	e.openOrders[order.Id] = order
	e.orderIdCounter++
}

func (e *Exchange) cancelOrder(id string) bool {

	if _, ok := e.openOrders[id]; ok {
		delete(e.openOrders, id)
		return true
	}

	return false
}

func (e *Exchange) cancelAllOrders() {
	e.openOrders = make(map[string]LimitOrder)
}

func (e *Exchange) getOpenOrders() []LimitOrder {

	orders := make([]LimitOrder, 0)
	for _, order := range e.openOrders {
		orders = append(orders, order)
	}

	return orders
}

func (e *Exchange) getFilledOrders() []LimitOrder {

	orders := make([]LimitOrder, 0)
	for _, order := range e.filledOrders {
		orders = append(orders, order)
	}

	return orders
}

func (e *Exchange) getMarkPrice() float64 {
	return e.markPrice
}

func (e *Exchange) getBalance() float64 {
	return e.balance
}

func (e *Exchange) getPosition() *Position {
	return e.position
}

func (e *Exchange) getTimestamp() int64 {
	return e.timestamp
}
