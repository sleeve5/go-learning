package main

import "fmt"

func main() {
	var age = 19
	var birthday = "20020813"
	var name = "Yin Xinyu"
	var result = fmt.Sprintf("\nname: %s, \nbirthday: %s, \nage %d. \n", name, birthday, age)
	fmt.Println(result)
}
