package main

import "fmt"

type ExchangeAPI struct {
	PlaceOrder          func(LimitOrder)
	CancelOrder         func(string) bool
	OpenOrders          func() map[string]LimitOrder
	GetMarkPrice        func() float64
	GetBalance          func() float64
	GetPosition         func() *Position
	GetCurrentTimestamp func() int64
}

type Exchange struct {
	timestamp  int64
	markPrice  float64
	balance    float64
	position   *Position
	openOrders map[string]LimitOrder

	NotifyTickUpdate func(Tick) // callback to notify the market maker

	orderIdCounter int64 // used for order ID
}

func NewExchange() *Exchange {
	return &Exchange{
		position:       &Position{},
		openOrders:     make(map[string]LimitOrder),
		orderIdCounter: 0,
	}
}

// PUBLIC METHODS
func (e *Exchange) Init(balance float64, tick Tick) {
	e.timestamp = tick.Timestamp
	e.markPrice = tick.Price
}

func (e *Exchange) Next(tick Tick) {
	// log.Debugf("### Next tick: date %s, mark price %.2f", tick.Time.String(), tick.Price)
	e.timestamp = tick.Timestamp
	e.markPrice = tick.Price

	ordersToExecute := e.getOrdersToExecute()
	for _, order := range ordersToExecute {
		e.executeOrder(order)
	}
}

func (e *Exchange) GetAPI() *ExchangeAPI {
	return &ExchangeAPI{
		PlaceOrder:          e.placeOrder,
		CancelOrder:         e.cancelOrder,
		OpenOrders:          e.getOpenOrders,
		GetMarkPrice:        e.getMarkPrice,
		GetBalance:          e.getBalance,
		GetPosition:         e.getPosition,
		GetCurrentTimestamp: e.getCurrentTimestamp,
	}
}

// PRIVATE METHODS
func (e *Exchange) getOrdersToExecute() map[string]LimitOrder {

	ordersToExecute := make(map[string]LimitOrder)
	for key, order := range e.openOrders {
		if order.Side == SIDE_BUY && order.Price >= e.markPrice {
			ordersToExecute[key] = order
		}
	}

	return ordersToExecute
}

func (e *Exchange) executeOrder(order LimitOrder) {
	// log.Debugf("Exchange: execute order %s", order.String())
	// if order.PositionSide == PositionSideLong {
	// 	if _, ok := e.sessionLong.openOrders[order.ID]; ok {
	// 		delete(e.sessionLong.openOrders, order.ID)
	// 	} else {
	// 		log.Panic("Order id not found in open orders")
	// 	}
	// 	e.sessionLong.orderAmount = order.Amount
	// 	e.sessionLong.orderPrice = order.Price
	// 	realizedProfit := e.sessionLong.position.Update(order)
	// 	e.sessionLong.realizedProfit = realizedProfit
	// 	e.balance += realizedProfit
	// 	if !order.IsTP { // don't update grid reached to 0 if is TP order, this is for the statistics
	// 		e.sessionLong.gridReached = order.GridNumber
	// 	}
	// 	log.Debugf("Exchange: updated position %s", e.sessionLong.position.String())
	// 	e.NotifyPositionUpdateCallback(e.sessionLong.position)
	// } else {
	// 	if _, ok := e.sessionShort.openOrders[order.ID]; ok {
	// 		delete(e.sessionShort.openOrders, order.ID)
	// 	} else {
	// 		log.Panic("Order id not found in open orders")
	// 	}
	// 	e.sessionShort.orderAmount = order.Amount
	// 	e.sessionShort.orderPrice = order.Price
	// 	realizedProfit := e.sessionShort.position.Update(order)
	// 	e.sessionShort.realizedProfit = realizedProfit
	// 	e.balance += realizedProfit
	// 	if !order.IsTP { // don't update grid reached to 0 if is TP order, this is for the statistics
	// 		e.sessionShort.gridReached = order.GridNumber
	// 	}
	// 	log.Debugf("Exchange: updated position %s", e.sessionShort.position.String())
	// 	e.NotifyPositionUpdateCallback(e.sessionShort.position)
	// }
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

func (e *Exchange) getOpenOrders() map[string]LimitOrder {
	return e.openOrders
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

func (e *Exchange) getCurrentTimestamp() int64 {
	return e.timestamp
}
