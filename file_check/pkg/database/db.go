package database

import (
	"database/sql"
	"log"

	// SQLite (без CGO)
	_ "modernc.org/sqlite"
)

// глобальная перменная SQL-запросов
var DB *sql.DB

func InitDB() { // функция создания таблицы и данных
	var err error

	DB, err = sql.Open("sqlite", "app.db")
	if err != nil {
		log.Fatal("Ошибка создания БД: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("БД недоступна: ", err)
	}

	// наша таблциа в sql с данными
	// хэш - как id(ключ)
	// а также данные модулей(позже сделаю итог проверки а не вывод теста)
	query := `
	CREATE TABLE IF NOT EXISTS scan_results (
		file_hash TEXT PRIMARY KEY,
		file_name TEXT,
		status TEXT,
		module_outputs TEXT
	);`

	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal("Ошибка создания SQL-таблицы: ", err)
	}

	log.Println("Создана таблица данных о сканировании")
}
