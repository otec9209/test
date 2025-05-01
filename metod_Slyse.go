package main

import (
	"fmt"
	"log"
)

func main() {

	var input string
	fmt.Println("Введите целое число")

	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Println("Ошибка ввода:", err)
		return
	}

	log.Println("Вы ввели целое число:", input)
}
