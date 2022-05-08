package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	// 设置最大值
	maxNum := 100
	rand.Seed(time.Now().UnixNano())
	result := rand.Intn(maxNum)
	for {

		// 利用scanf输入
		fmt.Println("Please input your guess:")
		var number int
		_, err := fmt.Scanf("%d", &number)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 判断部分
		if result < number {
			fmt.Println("Your guess is bigger than the result!")
		} else if result > number {
			fmt.Println("Your guess is smaller than the result!")
		} else if result == number {
			fmt.Println("Correct!")
			return
		}
	}
}
