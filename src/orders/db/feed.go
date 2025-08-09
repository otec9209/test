// Package db
package db

import (
	"log"
	"strings"
)

// GetOrCreateFeed находит ленту по URL или создаёт новую
// Возвращает feedID или ошибку
func GetOrCreateFeed(url string) (int, error) {
	var id int

	// 1. Попробовать найти существующую ленту
	err := DB.QueryRow("SELECT id FROM feeds WHERE url = $1", url).Scan(&id)
	if err == nil {
		// Лента уже есть → возвращаем её id
		return id, nil
	}

	// 2. Ленты нет → добавляем
	name := extractNameFromURL(url)

	err = DB.QueryRow(
		"INSERT INTO feeds (name, url) VALUES ($1, $2) RETURNING id",
		name, url,
	).Scan(&id)

	if err != nil {
		log.Printf("❌ Ошибка при добавлении ленты %s: %v", url, err)
		return 0, err
	}

	log.Printf("✅ Добавлена лента: %s (%s)", name, url)
	return id, nil
}

// extractNameFromURL — упрощённо определяет имя по URL
func extractNameFromURL(url string) string {
	if strings.Contains(url, "sport-express.ru") {
		return "Спорт-Экспресс"
	}
	if strings.Contains(url, "sportbox.ru") {
		return "Sportbox"
	}
	if strings.Contains(url, "sport.ru") {
		return "Sport.ru"
	}
	return "Неизвестный источник"
}
