package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func InitDB(){
	var err error
	DB, err = sql.Open("sqlite3", "./mireabot.db")

	if err != nil{
		log.Fatalf("Ошибка открытия БД")
	}
	stmt, err := DB.Prepare("CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, tg TEXT UNIQUE, login TEXT, password TEXT)")

	if err != nil{
		log.Fatalf("Ошибка create запроса")
	}

	_, err = stmt.Exec()
	if err != nil{
		log.Fatalf("Ошибка exec() у create()")
	}

	log.Print("БД инициализированна")
}
