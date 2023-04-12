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
	receivedCh := receiveOrders()
	validCh, errCh := validateOrders(receivedCh)
	reservedCh := reserveInteventory(validCh)

	wg.Add(2)
	go func(errCh <-chan invalidOrder) {
		for order := range errCh {
			fmt.Printf("Invalid order received: %v. Issue: %v\n", order.order, order.err)
		}
		wg.Done()
	}(errCh)

	go func(reservedCh <-chan order) {
		for order := range reservedCh {
			fmt.Printf("Inventory reserved for: %v\n", order)
		}
		wg.Done()
	}(reservedCh)

	wg.Wait()
}

func reserveInteventory(validCh <-chan order) <-chan order {
	reservedCh := make(chan order)
	go func() {
		for order := range validCh {
			order.Status = reserved
			reservedCh <- order
		}
		close(reservedCh)
	}()

	return reservedCh
}

// I added the directional restrictions to the channels myself
// We are using the encapsulation of goroutines pattern here..
func validateOrders(receivedCh <-chan order) (<-chan order, <-chan invalidOrder) {
	validCh := make(chan order)
	errCh := make(chan invalidOrder, 1) // buffered error channel

	// We can range over our in channel
	go func() {
		for order := range receivedCh {
			if order.Quantity <= 0 {
				errCh <- invalidOrder{
					order: order, err: errors.New("quantity must be greater than zero"),
				}
			} else {
				validCh <- order
			}
		}
		close(validCh)
		close(errCh)
	}()
	// Close the channels after we sent the messages to them
	return validCh, errCh
}

// We are using the encapsulation of goroutines pattern here..
func receiveOrders() <-chan order {
	receivedCh := make(chan order)

	go func() {
		for _, rawOrder := range rawOrders {
			var newOrder order
			err := json.Unmarshal([]byte(rawOrder), &newOrder)
			if err != nil {
				log.Print(err)
				continue
			}
			receivedCh <- newOrder
		}
		close(receivedCh)
	}()
	return receivedCh
}

var rawOrders = []string{
	`{"productCode": 1111, "quantity": 5, "status": 1}`,
	`{"productCode": 2222, "quantity": 42.3, "status": 1}`,
	`{"productCode": 3333, "quantity": 19, "status": 1}`,
	`{"productCode": 4444, "quantity": 8, "status": 1}`,
}
