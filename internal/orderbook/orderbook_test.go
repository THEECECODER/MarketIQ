package orderbook

import "testing"

func TestBestBid(t *testing.T) {
	book := NewOrderBook()

	book.UpdateBid(10000, 200)
	book.UpdateBid(9900, 500)
	book.UpdateBid(9800, 700)

	best, ok := book.BestBid()

	if !ok {
		t.Fatal("expected a best bid")
	}

	if best != 10000 {
		t.Fatalf("expected best bid 10000, got %d", best)
	}
}

func TestBestAsk(t *testing.T) {
	book := NewOrderBook()

	book.UpdateAsk(10100, 200)
	book.UpdateAsk(10200, 500)
	book.UpdateAsk(10300, 700)

	best, ok := book.BestAsk()

	if !ok {
		t.Fatal("expected a best ask")
	}

	if best != 10100 {
		t.Fatalf("expected best ask 10100, got %d", best)
	}
}
func TestMidPrice(t *testing.T) {
	book := NewOrderBook()

	book.UpdateBid(10000, 200)
	book.UpdateAsk(10100, 300)

	mid, ok := book.MidPrice()

	if !ok {
		t.Fatal("expected mid price")
	}

	if mid != 10050 {
		t.Fatalf("expected mid price 10050, got %d", mid)
	}
}
func TestSpread(t *testing.T) {
	book := NewOrderBook()

	book.UpdateBid(10000, 200)
	book.UpdateAsk(10100, 300)

	spread, ok := book.Spread()

	if !ok {
		t.Fatal("expected spread")
	}

	if spread != 100 {
		t.Fatalf("expected spread 100, got %d", spread)
	}
}
func TestDepth(t *testing.T) {
	book := NewOrderBook()

	book.UpdateBid(10000, 200)
	book.UpdateBid(9900, 300)
	book.UpdateBid(9800, 500)

	book.UpdateAsk(10100, 100)
	book.UpdateAsk(10200, 200)

	if book.BidDepth() != 1000 {
		t.Fatalf("expected bid depth 1000, got %d", book.BidDepth())
	}

	if book.AskDepth() != 300 {
		t.Fatalf("expected ask depth 300, got %d", book.AskDepth())
	}
}
