package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/time", GetTime)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func GetTime(w http.ResponseWriter, r *http.Request) {
	log.Printf("Выполнен запрос от пользователя: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "<h1>%s</h1><br>", "Запуск приложения в контейнере Docker")
	fmt.Fprintf(w, "Текущие дата и время %s\n",
		time.Now().Format("02-01-2006 15:04"))
}
