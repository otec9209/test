package db

import "log"

func CreateTables() error {
	query := `
    CREATE TABLE IF NOT EXISTS feeds (
        id         SERIAL PRIMARY KEY,
        name       TEXT NOT NULL,
        url        TEXT UNIQUE NOT NULL,
        active     BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS items (
        id          SERIAL PRIMARY KEY,
        title       TEXT NOT NULL,
        link        TEXT UNIQUE NOT NULL,
        description TEXT,
        published   TEXT,
        feed_id     INTEGER NOT NULL,
        read        BOOLEAN DEFAULT FALSE,
        created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

        FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE
    );`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	log.Println("✅ Таблицы 'feeds' и 'items' готовы к работе")
	return nil
}
