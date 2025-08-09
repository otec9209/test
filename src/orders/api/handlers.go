package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"orders/db"
)

// Post — структура для ответа API
// Должна совпадать с db.Post
type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description,omitempty"`
	Published   string `json:"published"`
	FeedID      int    `json:"feed_id"`
	Read        bool   `json:"read"`
}

// GetPosts возвращает все новости из БД
func GetPosts(w http.ResponseWriter, r *http.Request) {
	// 1. Запрос к БД
	rows, err := db.DB.Query(`
		SELECT id, title, link, description, published, feed_id, read 
		FROM items 
		ORDER BY published DESC`)
	if err != nil {
		http.Error(w, "Ошибка БД", 500)
		return
	}
	defer rows.Close()

	// 2. Собираем список новостей
	var posts []Post
	for rows.Next() {
		var p Post
		var description sql.NullString

		err := rows.Scan(&p.ID, &p.Title, &p.Link, &description, &p.Published, &p.FeedID, &p.Read)
		if err != nil {
			continue
		}

		// description может быть NULL
		if description.Valid {
			p.Description = description.String
		}

		posts = append(posts, p)
	}

	// 3. Отправляем как JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GetFeeds возвращает список источников
func GetFeeds(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, url FROM feeds ORDER BY name")
	if err != nil {
		http.Error(w, "Ошибка БД", 500)
		return
	}
	defer rows.Close()

	type Feed struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	var feeds []Feed
	for rows.Next() {
		var f Feed
		rows.Scan(&f.ID, &f.Name, &f.URL)
		feeds = append(feeds, f)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeds)
}
