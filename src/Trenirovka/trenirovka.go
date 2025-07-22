package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func main() {
	// Массив с поговорками
	proverbs := []string{
		"Concurrency is not parallelism.",
		"Make the zero value useful.",
		"Cgo must always be guarded with build tags.",
		"With the unsafe package there are no guarantees.",
		"Don't just check errors, handle them gracefully.",
	}

	// Создаём сервер
	listener, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Сервер запущен на 127.0.0.1:8081")

	// Ждём подключений
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка подключения:", err)
			continue
		}

		// Обработка клиента
		go func(conn net.Conn) {
			defer conn.Close()
			fmt.Println("Клиент подключился:", conn.RemoteAddr())

			// Инициализируем генератор случайных чисел
			rand.Seed(time.Now().UnixNano())

			// Бесконечный цикл: каждые 3 секунды отправляем поговорку
			for {
				// Выбираем случайную поговорку
				quote := proverbs[rand.Intn(len(proverbs))]

				// Отправляем клиенту
				_, err := conn.Write([]byte(quote + "\n"))
				if err != nil {
					fmt.Println("Клиент отключился или ошибка:", err)
					return
				}

				// Ждём 3 секунды
				time.Sleep(3 * time.Second)
			}
		}(conn)
	}
}
