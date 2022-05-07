package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// 生成随机数
	maxNum := 100
	rand.Seed(time.Now().UnixNano())
	result := rand.Intn(maxNum)
	// fmt.Printf("The result is %v.\n", result)

	// 读取用户输入
	fmt.Println("Please input your guess:")
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occurred while reading input. Please try again.", err)
			continue
		}

		// windows 与 linux不同，末尾加上\r\n作为换行符
		input = strings.TrimSuffix(input, "\r\n")

		// 输入转化为数字
		number, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please input an integer value.", err)
			continue
		}
		fmt.Printf("Your guess is %v.\n", number)

		// 实现判断
		switch {
		case result < number:
			{
				fmt.Println("Your guess is bigger than the result!")
			}
		case result > number:
			{
				fmt.Println("Your guess is smaller than the result!")
			}
		case result == number:
			{
				fmt.Println("Correct!")
				break
			}
		}
	}
}
