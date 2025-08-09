// cmd/main.go
package main

import (
	"log"
	"net/http"
	"orders/api"
	"orders/db"
	"orders/rss"
	"time"
)

func main() {
	// 1. Подключиться к БД
	if err := db.Connect(); err != nil {
		log.Fatal("❌ Ошибка подключения:", err)
	}
	defer db.Close()

	// 2. Создать таблицы
	if err := db.CreateTables(); err != nil {
		log.Fatal("❌ Ошибка таблиц:", err)
	}

	// 3. Загрузить конфиг
	cfg, err := rss.Load("config.json")
	if err != nil {
		log.Fatal("❌ Ошибка конфига:", err)
	}

	// Запустить API ДО цикла
	api.RegisterRoutes()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})

	// Запускаем HTTP-сервер в отдельной горутине
	go func() {
		log.Println("🌍 API запущено на :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("❌ Ошибка запуска HTTP-сервера:", err)
		}
	}()

	// 5. Бесконечный цикл с таймером
	ticker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer ticker.Stop()

	log.Printf("✅ Агрегатор запущен. Опрос каждые %d секунд", cfg.PollInterval)

	for range ticker.C {
		log.Println("🔄 Начинаю опрос всех лент...")

		// 5. Для каждой ленты — запустить в отдельной горутине
		for _, feedURL := range cfg.Feeds {
			// Захватываем переменную
			feedURL := feedURL

			go func(feedURL string) {
				// 5.1. Получить или создать feedID
				feedID, err := db.GetOrCreateFeed(feedURL)
				if err != nil {
					log.Printf("❌ Ошибка с лентой %s: %v", feedURL, err)
					return
				}

				// 5.2. Парсим ленту
				posts, err := rss.FetchFeed(feedURL, feedID)
				if err != nil {
					log.Printf("❌ Ошибка парсинга %s: %v", feedURL, err)
					return
				}

				// 5.3. Присвоить feedID
				for i := range posts {
					posts[i].FeedID = feedID
				}

				// 5.4. Сохранить
				db.SavePosts(posts)
			}(feedURL)
		}

		// Не ждём завершения — пусть работают в фоне
		log.Println("💤 Опрос запущен. Жду следующего таймера...")
	}
}
