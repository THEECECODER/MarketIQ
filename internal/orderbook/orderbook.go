package orderbook

import "sort"

type PriceLevel struct {
	Price    int64
	Quantity int64
}

type OrderBook struct {
	Bids map[int64]int64
	Asks map[int64]int64
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids: make(map[int64]int64),
		Asks: make(map[int64]int64),
	}
}
func (ob *OrderBook) UpdateBid(price int64, quantity int64) {
	if quantity == 0 {
		delete(ob.Bids, price)
		return
	}
	ob.Bids[price] = quantity
}
func (ob *OrderBook) UpdateAsk(price int64, quantity int64) {
	if quantity == 0 {
		delete(ob.Asks, price)
		return
	}
	ob.Asks[price] = quantity
}

func (ob *OrderBook) BestBid() (int64, bool) {
	var best int64
	found := false
	for price := range ob.Bids {
		if !found || price > best {
			best = price
			found = true
		}
	}
	return best, found
}
func (ob *OrderBook) BestAsk() (int64, bool) {
	var best int64
	found := false

	for price := range ob.Asks {
		if !found || price < best {
			best = price
			found = true
		}
	}

	return best, found
}
func (ob *OrderBook) MidPrice() (int64, bool) {
	bid, bidOK := ob.BestBid()
	ask, askOK := ob.BestAsk()

	if !bidOK || !askOK {
		return 0, false
	}

	return (bid + ask) / 2, true
}
func (ob *OrderBook) Spread() (int64, bool) {
	bid, bidOK := ob.BestBid()
	ask, askOK := ob.BestAsk()
	if !bidOK || !askOK {
		return 0, false
	}

	return ask - bid, true
}

// top bids
func (ob *OrderBook) TopBids(n int) []PriceLevel {
	levels := make([]PriceLevel, 0, len(ob.Bids))

	for price, quantity := range ob.Bids {
		levels = append(levels, PriceLevel{
			Price:    price,
			Quantity: quantity,
		})
	}

	sort.Slice(levels, func(i, j int) bool {
		return levels[i].Price > levels[j].Price
	})

	if n > len(levels) {
		n = len(levels)
	}

	return levels[:n]
}

// top n asks
func (ob *OrderBook) TopAsks(n int) []PriceLevel {
	levels := make([]PriceLevel, 0, len(ob.Asks))

	for price, quantity := range ob.Asks {
		levels = append(levels, PriceLevel{
			Price:    price,
			Quantity: quantity,
		})
	}

	sort.Slice(levels, func(i, j int) bool {
		return levels[i].Price < levels[j].Price
	})

	if n > len(levels) {
		n = len(levels)
	}

	return levels[:n]
}
func (ob *OrderBook) BidDepth() int64 {
	var total int64

	for _, quantity := range ob.Bids {
		total += quantity
	}

	return total
}

func (ob *OrderBook) AskDepth() int64 {
	var total int64

	for _, quantity := range ob.Asks {
		total += quantity
	}

	return total
}
