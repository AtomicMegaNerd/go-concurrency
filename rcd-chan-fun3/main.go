package main

import "fmt"

var arr []string

func main() {
	ch := make(chan string, 3)

	for _, v := range [...]arr{"foo", "bar", "baz"} {
		ch <- v
	}
	close(ch)

	for msg := range ch {
		fmt.Println(msg)
	}
}
