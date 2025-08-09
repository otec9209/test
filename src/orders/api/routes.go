package api

import "net/http"

// RegisterRoutes регистрирует все маршруты API
func RegisterRoutes() {
	http.HandleFunc("/api/posts", GetPosts)
	http.HandleFunc("/api/feeds", GetFeeds)
	// Можно добавить: /api/posts?feed_id=1
}
