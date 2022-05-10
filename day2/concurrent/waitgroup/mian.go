package main

import "sync"

func hello(i int) {
	println("goroutine:", i)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		j := i
		go func() {
			defer wg.Done()
			hello(j)
		}()
	}
	wg.Wait()
}
