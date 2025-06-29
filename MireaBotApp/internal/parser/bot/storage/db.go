package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"log"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./mireabot.db")

	if err != nil {
		log.Fatalf("Ошибка открытия БД")
	}
	stmt, err := DB.Prepare("CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY,chatID INTEGER UNIQUE, tg TEXT UNIQUE, login TEXT, password TEXT)")

	if err != nil {
		logrus.Fatal("Ошибка create запроса")
	}

	_, err = stmt.Exec()
	if err != nil {
		logrus.Fatal("Ошибка exec() у create()")
	}

	logrus.Info("БД инициализированна")
}
