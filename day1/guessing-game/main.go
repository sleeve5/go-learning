package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	maxNum := 100
	rand.Seed(time.Now().UnixNano())
	result := rand.Intn(maxNum)
	fmt.Printf("The result is %v.\n", result)
}
