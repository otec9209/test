package rss

import (
	"encoding/xml" // для парсинга XML
	"fmt"          // для форматирования ошибок
	"io"           // для чтения тела ответа
	"net/http"     // для HTTP-запросов
	"time"         // для работы с датами
)

// FetchFeed — загружает и парсит одну RSS-ленту по URL.
// Возвращает список новостей или ошибку.
func FetchFeed(url string, feedID int) ([]Post, error) {
	// 1. Делаем HTTP-запрос по ссылке
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе к %s: %w", url, err)
	}
	// Обязательно закрываем тело ответа
	defer resp.Body.Close()

	// 2. Читаем весь XML
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	// 3. Создаём переменную для результата парсинга
	var feed RSSFeed

	// 4. Распарсиваем XML в структуру
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга XML: %w", err)
	}

	// 5. Преобразуем результат в []Post
	var posts []Post
	for _, item := range feed.Items {
		// Парсим дату, если нужно, но оставляем как строку
		pubTime, err := parseRSSDate(item.PubDate)
		var publishedStr string
		if err != nil {
			// Если не получилось — оставляем как есть
			publishedStr = item.PubDate
		} else {
			// Форматируем в понятный вид: "2025-04-01 12:00:00"
			publishedStr = pubTime.Format("2006-01-02 15:04:05")
		}

		// Добавляем новость в список
		posts = append(posts, Post{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Published:   publishedStr,
			FeedID:      feedID,
		})
	}

	// 6. Возвращаем список новостей
	return posts, nil
}

// parseRSSDate — пробует распарсить строку даты в нескольких форматах
func parseRSSDate(dateStr string) (time.Time, error) {
	formats := []string{
		"Mon, 02 Jan 2006 15:04:05 GMT",
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 -0700 (MST)",
		time.RFC1123Z,
		time.RFC1123,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
	}

	for _, layout := range formats {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("не удалось распарсить дату: %s", dateStr)
}
