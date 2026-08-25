package marketdata

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"marketiq/internal/orderbook"

	socketio "github.com/Joaquimborges/go-socket.io"
)

const (
	CoinDCXURL = "wss://stream-spot.coindcx.com"
	Channel    = "B-BTC_USDT@orderbook@20"
)

type Client struct {
	socket    *socketio.Client
	orderBook *orderbook.OrderBook
}

type DepthSnapshot struct {
	Version   int64             `json:"vs"`
	Timestamp int64             `json:"ts"`
	Asks      map[string]string `json:"asks"`
	Bids      map[string]string `json:"bids"`
	Market    string            `json:"pr"`
	Symbol    string            `json:"s"`
}

func NewClient() *Client {
	return &Client{
		orderBook: orderbook.NewOrderBook(),
	}
}

func parsePrice(value string) (int64, error) {
	parts := strings.Split(value, ".")

	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid price %q", value)
	}

	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid price %q: %w", value, err)
	}

	var fraction int64

	if len(parts) == 2 {
		decimal := parts[1]

		if len(decimal) > 2 {
			return 0, fmt.Errorf(
				"price has more than 2 decimal places: %q",
				value,
			)
		}

		if len(decimal) == 1 {
			decimal += "0"
		}

		fraction, err = strconv.ParseInt(decimal, 10, 64)
		if err != nil {
			return 0, fmt.Errorf(
				"invalid price %q: %w",
				value,
				err,
			)
		}
	}

	return whole*100 + fraction, nil
}

func parseQuantity(value string) (int64, error) {
	parts := strings.Split(value, ".")

	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid quantity %q", value)
	}

	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"invalid quantity %q: %w",
			value,
			err,
		)
	}

	var fraction int64

	if len(parts) == 2 {
		decimal := parts[1]

		if len(decimal) > 5 {
			return 0, fmt.Errorf(
				"quantity has more than 5 decimal places: %q",
				value,
			)
		}

		for len(decimal) < 5 {
			decimal += "0"
		}

		fraction, err = strconv.ParseInt(decimal, 10, 64)
		if err != nil {
			return 0, fmt.Errorf(
				"invalid quantity %q: %w",
				value,
				err,
			)
		}
	}

	return whole*100000 + fraction, nil
}

func ConvertSnapshot(
	snapshot DepthSnapshot,
	book *orderbook.OrderBook,
) error {

	for price, quantity := range snapshot.Bids {
		p, err := parsePrice(price)
		if err != nil {
			return err
		}

		q, err := parseQuantity(quantity)
		if err != nil {
			return err
		}

		book.UpdateBid(p, q)
	}

	for price, quantity := range snapshot.Asks {
		p, err := parsePrice(price)
		if err != nil {
			return err
		}

		q, err := parseQuantity(quantity)
		if err != nil {
			return err
		}

		book.UpdateAsk(p, q)
	}

	return nil
}

func formatPrice(value int64) string {
	return fmt.Sprintf("%.2f", float64(value)/100)
}

func (c *Client) Connect() error {
	client, err := socketio.NewClient(CoinDCXURL)
	if err != nil {
		return err
	}

	c.socket = client

	client.OnConnect(func() {
		log.Println("Connected to CoinDCX")
		log.Println("Subscribing to:", Channel)

		err := client.Emit("join", map[string]string{
			"channelName": Channel,
		})

		if err != nil {
			log.Println("Subscription error:", err)
		}
	})

	client.OnDisconnect(func(err error) {
		log.Println("Disconnected from CoinDCX:", err)
	})

	client.On("error", func(err error) {
		log.Println("Socket error:", err)
	})

	client.On("depth-snapshot", func(data interface{}) {
		log.Println("Received depth snapshot")

		payload, ok := data.(map[string]interface{})
		if !ok {
			log.Printf(
				"Unexpected snapshot format: %T\n",
				data,
			)
			return
		}

		rawData, ok := payload["data"].(string)
		if !ok {
			log.Println("Snapshot does not contain data field")
			return
		}

		var snapshot DepthSnapshot

		if err := json.Unmarshal(
			[]byte(rawData),
			&snapshot,
		); err != nil {
			log.Println(
				"Failed to parse snapshot:",
				err,
			)
			return
		}

		if err := ConvertSnapshot(
			snapshot,
			c.orderBook,
		); err != nil {
			log.Println(
				"Failed to update order book:",
				err,
			)
			return
		}

		bid, bidOK := c.orderBook.BestBid()
		ask, askOK := c.orderBook.BestAsk()

		if bidOK {
			log.Printf(
				"Best Bid: %s\n",
				formatPrice(bid),
			)
		}

		if askOK {
			log.Printf(
				"Best Ask: %s\n",
				formatPrice(ask),
			)
		}

		mid, midOK := c.orderBook.MidPrice()
		if midOK {
			log.Printf(
				"Mid Price: %s\n",
				formatPrice(mid),
			)
		}

		spread, spreadOK := c.orderBook.Spread()
		if spreadOK {
			log.Printf(
				"Spread: %s\n",
				formatPrice(spread),
			)
		}
	})

	return client.Connect()
}
