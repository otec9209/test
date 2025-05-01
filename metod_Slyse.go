package main

import (
	"fmt"
	"log"
)
func main() {

	var input string
	fmt.Println("Введите данные")

	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Println("Ошибка ввода:", err)
		return
	}

|
	log.Println("Вы ввели следующие даннеы:", input) //я не знаю но что-то надо исправить по заданию, вношу эту строку
}

