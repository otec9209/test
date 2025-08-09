package rss

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

// Config — структура конфигурации
type Config struct {
	Feeds        []string `json:"feeds"`
	PollInterval int      `json:"poll_interval"`
}

// Load — функция загрузки и валидации конфигурации из файла
func Load(path string) (Config, error) {
	var cfg Config

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	if len(cfg.Feeds) == 0 {
		return cfg, fmt.Errorf("feeds cannot be empty")
	}
	for _, feed := range cfg.Feeds {
		u, err := url.Parse(feed)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			return cfg, fmt.Errorf("invalid feed URL: %s", feed)
		}
	}
	if cfg.PollInterval <= 0 {
		return cfg, fmt.Errorf("poll_interval must be positive")
	}

	return cfg, nil
}
