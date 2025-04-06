package main

import (
	"fmt"
	"strconv"
)

// Функция для вычисления суммы
func main() {
	var input string
	numbers := []int{}

	fmt.Println("Введите числа:(для завершения введите 'q')")
	for {
		fmt.Scan(&input)
		if input == "q" {
			break
		}
		number, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Введено говно,введите число")
			continue
		}
		numbers = append(numbers, number)
	}
	sum := 0
	for _, n := range numbers {
		sum += n

		// 1

	}
	fmt.Println("текущая сумма равна", sum)
}
