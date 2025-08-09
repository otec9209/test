package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // драйвер PostgreSQL
)

var DB *sql.DB

// Connect — подключается к PostgreSQL
func Connect() error {

	connStr := "host=localhost user=postgres password=1234 dbname=orders sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("✅ Подключено к PostgreSQL: orders")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("🔌 PostgreSQL соединение закрыто")
	}
}
