package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	in := receiveOrders()
	out, errCh := validateOrders(in)

	wg.Add(1)
	go func(valCh <-chan order, invCh <-chan invalidOrder) {
	loop: // Label for our break
		for {
			select {
			case order, ok := <-valCh:
				if ok {
					fmt.Printf("Valid order received: %v\n", order)
				} else {
					break loop
				}
			case order, ok := <-invCh:
				if ok {
					fmt.Printf("Invalid order received: %v, Issue %v\n", order.order, order.err)
				} else {
					break loop
				}
			}
		}
		wg.Done()
	}(out, errCh)

	wg.Wait()
}

// I added the directional restrictions to the channels myself
// We are using the encapsulation of goroutines pattern here..
func validateOrders(in <-chan order) (<-chan order, <-chan invalidOrder) {
	out := make(chan order)
	errCh := make(chan invalidOrder, 1) // buffered error channel

	// We can range over our in channel
	go func() {
		for order := range in {
			if order.Quantity <= 0 {
				errCh <- invalidOrder{
					order: order, err: errors.New("quantity must be greater than zero"),
				}
			} else {
				out <- order
			}
		}
		close(out)
		close(errCh)
	}()
	// Close the channels after we sent the messages to them
	return out, errCh
}

// We are using the encapsulation of goroutines pattern here..
func receiveOrders() <-chan order {
	out := make(chan order)

	go func() {
		for _, rawOrder := range rawOrders {
			var newOrder order
			err := json.Unmarshal([]byte(rawOrder), &newOrder)
			if err != nil {
				log.Print(err)
				continue
			}
			out <- newOrder
		}
		close(out)
	}()
	return out
}

var rawOrders = []string{
	`{"productCode": 1111, "quantity": 5, "status": 1}`,
	`{"productCode": 2222, "quantity": 42.3, "status": 1}`,
	`{"productCode": 3333, "quantity": 19, "status": 1}`,
	`{"productCode": 4444, "quantity": 8, "status": 1}`,
}
