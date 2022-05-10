package main

import "fmt"

func main() {
	src := make(chan int)
	dest := make(chan int)
	go func() {
		defer close(src)
		for i := 0; i < 5; i++ {
			src <- i
		}
	}()
	go func() {
		defer close(dest)
		for i := range src {
			dest <- i * i
		}
	}()
	for i := range dest {
		fmt.Printf("%v\n", i)
	}
}
