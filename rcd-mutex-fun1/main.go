package main

import (
	"fmt"
	"sync"
)

func main() {
	s := []int{}
	var wg sync.WaitGroup
	var m sync.Mutex

	const iterations = 1000
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func() {
			m.Lock()
			defer m.Unlock()

			s = append(s, 1)
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println(len(s))

}
