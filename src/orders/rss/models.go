package rss

// Post — структура для новости, которую мы будем сохранять в БД.
// Поля соответствуют колонкам в таблице items.
type Post struct {
	Title       string
	Link        string
	Description string
	Published   string
	FeedID      int
}

// RSSFeed — временная структура для парсинга XML.
// Она "копирует" структуру RSS-ленты, чтобы Go мог заполнить её из XML.
type RSSFeed struct {
	Items []struct {
		Title       string `xml:"title"`       // <title>Месси</title>
		Link        string `xml:"link"`        // <link>https://...</link>
		Description string `xml:"description"` // <description>...</description>
		PubDate     string `xml:"pubDate"`     // <pubDate>Mon, 01 Apr 2025</pubDate>
	} `xml:"channel>item"` // ищи <item> внутри <channel>
}
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}
