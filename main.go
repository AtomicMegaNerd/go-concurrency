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
	revOrdCh := make(chan order)
	validOdCh := make(chan order)
	invalidOrdCh := make(chan invalidOrder)

	go receiveOrders(revOrdCh)
	go validateOrders(revOrdCh, validOdCh, invalidOrdCh)

	wg.Add(1)
	go func(valCh <-chan order, invCh <-chan invalidOrder) {
		// This code is really gross... there HAS to be a better way
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
	}(validOdCh, invalidOrdCh)

	wg.Wait()
}

// I added the directional restrictions to the channels myself
func validateOrders(in <-chan order, out chan<- order, errCh chan<- invalidOrder) {
	// We can range over our in channel
	for order := range in {
		if order.Quantity <= 0 {
			errCh <- invalidOrder{
				order: order, err: errors.New("quantity must be greater than zero"),
			}
		} else {
			out <- order
		}
	}
	// Close the channels after we sent the messages to them
	close(out)
	close(errCh)
}

func receiveOrders(out chan<- order) {
	for _, rawOrder := range rawOrders {
		var newOrder order
		err := json.Unmarshal([]byte(rawOrder), &newOrder)
		if err != nil {
			log.Print(err)
			continue
		}
		out <- newOrder
	}
	// Close the channel after we have received the orders
	close(out)
}

var rawOrders = []string{
	`{"productCode": 1111, "quantity": 5, "status": 1}`,
	`{"productCode": 2222, "quantity": 42.3, "status": 1}`,
	`{"productCode": 3333, "quantity": 19, "status": 1}`,
	`{"productCode": 4444, "quantity": 8, "status": 1}`,
}
