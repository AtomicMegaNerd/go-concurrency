package main

import "strconv"

var (
	in = make(chan string)
)

func worker(in <-chan string) (out chan<- int, errCh chan<- error) {
	out = make(chan int)
	errCh = make(chan error, 1)
	defer close(out)
	defer close(errCh)

	go func() {
		for msg := range in {
			i, err := strconv.Atoi(msg)
			if err != nil {
				errCh <- err
				return
			}
			out <- i
		}
	}()
	return out, errCh
}

func main() {

}
