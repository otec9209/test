package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	// 1. Читаем файл
	data, err := os.ReadFile("read.txt")
	if err != nil {
		panic("Ошибка чтения файла: " + err.Error())
	}

	// 2. Преобразуем в строку
	content := string(data)

	// 3. Убираем лишние строки — оставляем только уравнения
	reClean := regexp.MustCompile(`(?m)^\D.*$\n?`)
	content = reClean.ReplaceAllString(content, "")

	// 4. Обрабатываем уравнения (учитываем знак ?)
	reMath := regexp.MustCompile(`(\d+)([+-])(\d+)=\?`)
	content = reMath.ReplaceAllStringFunc(content, func(match string) string {
		parts := reMath.FindStringSubmatch(match)
		a, _ := strconv.Atoi(parts[1])
		op := parts[2]
		b, _ := strconv.Atoi(parts[3])

		var result int
		if op == "+" {
			result = a + b
		} else {
			result = a - b
		}

		return fmt.Sprintf("%d%s%d=%d", a, op, b, result)
	})

	// 5. Записываем результат в файл
	err = os.WriteFile("write.txt", []byte(content), 0644)
	if err != nil {
		panic("Ошибка записи файла: " + err.Error())
	}

	fmt.Println("✅ Уравнения решены и записаны в файл write.txt!")
}
