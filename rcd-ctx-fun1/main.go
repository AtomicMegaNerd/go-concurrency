package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	// WithTimeout with a duration
	// WithDeadline works with a time on the clock
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	// Always make the context the 1st argument by convention
	go func(ctx context.Context) {
		defer wg.Done()
		for range time.Tick(500 * time.Millisecond) {
			// cancelled if err != nil
			if ctx.Err() != nil {
				log.Println(ctx.Err())
				return
			}
			fmt.Println("tick!")
		}
	}(ctx)

	wg.Wait()
}
