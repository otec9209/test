package db

import (
	"log"
	"orders/rss"
)

// Post — структура для новости продублировал чтобы не было зависимости бд от рсс
type Post struct {
	Title       string
	Link        string
	Description string
	Published   string
	FeedID      int
}

// SavePosts сохраняет список новостей в таблицу items.
func SavePosts(posts []rss.Post) error {
	for _, post := range posts {
		_, err := DB.Exec(
			`INSERT INTO items (title, link, description, published, feed_id) 
			  VALUES ($1, $2, $3, $4, $5) 
			  ON CONFLICT (link) DO NOTHING`,
			post.Title,
			post.Link,
			post.Description,
			post.Published,
			post.FeedID, // будет установлен при вызове
		)
		if err != nil {
			log.Printf("❌ Ошибка сохранения новости %s: %v", post.Title, err)
			continue
		}
		log.Printf("✅ Сохранено: %s", post.Title)
	}
	return nil
}
