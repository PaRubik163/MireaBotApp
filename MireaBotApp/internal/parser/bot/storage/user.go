package storage
//Сделать шифрование паролей
//Добавить проверку валидности данных
//Добавить функционал с UPDATE
import (
	"log"
)

func IsExists(tgID string) bool {
	var exists int
	stmt, err := DB.Prepare("SELECT 1 FROM user WHERE tg = ? LIMIT 1")

	if err != nil {
		log.Fatalf("Ошибка exists запроса")
	}

	stmt.QueryRow(tgID).Scan(&exists)

	if exists != 1{
		return false
	}

	return true
}

func Insert(tgId, login, password string) {
	stmt, err := DB.Prepare("INSERT INTO user (tg, login, password) VALUES (?, ?, ?)")
	
	if err != nil{
		log.Fatalf("Ошибка INSERT запроса")
	}
	
	_, err = stmt.Exec(tgId, login, password)

	if err != nil{
		log.Fatalf("Ошибка INSERT.exec() запроса")
	}

	log.Println("Успешный INSERT запрос добавлен",tgId, login)
}

func Select(tgID string) (string,string) {
	var login, password string

	err := DB.QueryRow("SELECT login,password FROM user WHERE tg = ?", tgID).Scan(&login, &password)

	if err != nil{
		log.Fatalf("Ошибка Select запроса")
	}

	if login == "" || password == ""{
		return "", ""
	}

	log.Println("Успешный SELECT запрос вытащен", tgID, login)
	return login, password
}

func Update(tgID string, newlogin, newpassword string) bool {
	stmt, err := DB.Prepare("UPDATE user SET login = ?, password = ? WHERE tg = ?")

	if err != nil{
		log.Fatalf("Ошибка UPDATE запроса")
	}

	_, err = stmt.Exec(newlogin, newpassword, tgID)
	if err != nil{
		log.Fatalf("Ошибка UPDATE.EXEC() запроса")
		return false
	}

	log.Print("Успешный UPDATE запрос измене пользователь", tgID)

	return true
}