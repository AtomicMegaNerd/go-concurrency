package main

import "fmt"

func main() {
	ch1 := make(chan int, 1)
	ch2 := make(chan string, 1)

	ch1 <- 999
	ch2 <- "message"

	// Avoid deterministic behaviours when working with channels... this is
	// undefined behavior, a case will be randomly selected
	select {
	case msg := <-ch1:
		fmt.Println(msg)
	case msg := <-ch2:
		fmt.Println(msg)
	default:
		fmt.Println("default")
	}
}
