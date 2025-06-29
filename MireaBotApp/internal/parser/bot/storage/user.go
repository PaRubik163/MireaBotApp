package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"io"
	"log"
)

func IsExists(tgID string) bool {
	var exists int
	stmt, err := DB.Prepare("SELECT 1 FROM user WHERE tg = ? LIMIT 1")

	if err != nil {
		logrus.Fatalf("Ошибка exists запроса")
	}

	stmt.QueryRow(tgID).Scan(&exists)

	if exists != 1 {
		return false
	}

	return true
}

func Insert(chatID int, tgId, login, password string, key []byte) {
	stmt, err := DB.Prepare("INSERT INTO user (chatID, tg, login, password) VALUES (?, ?, ?, ?)")

	if err != nil {
		logrus.Fatalf("Ошибка INSERT запроса")
	}

	cipherPassword := encrypt(password, key)

	_, err = stmt.Exec(chatID, tgId, login, cipherPassword)

	if err != nil {
		logrus.Fatalf("Ошибка INSERT.exec() запроса")
	}

	log.Println("Успешный INSERT запрос добавлен", tgId, login)
}

func SelectLoginandPassword(tgID string, key []byte) (string, string) {
	var login, encPassword string

	err := DB.QueryRow("SELECT login,password FROM user WHERE tg = ?", tgID).Scan(&login, &encPassword)

	if err != nil {
		log.Fatalf("Ошибка Select запроса")
	}

	password := decrypt(encPassword, key)

	if login == "" || password == "" {
		return "", ""
	}

	logrus.Info("Успешный SELECT запрос вытащен", tgID, login)

	return login, password
}

func Update(tgID string, newlogin, newpassword string, key []byte) bool {
	stmt, err := DB.Prepare("UPDATE user SET login = ?, password = ? WHERE tg = ?")

	if err != nil {
		logrus.Fatalf("Ошибка UPDATE запроса")
	}

	cipherNewPassword := encrypt(newpassword, key)

	_, err = stmt.Exec(newlogin, cipherNewPassword, tgID)
	if err != nil {
		logrus.Fatalf("Ошибка UPDATE.EXEC() запроса")
		return false
	}

	logrus.Info("Успешный UPDATE запрос измене пользователь ", tgID)

	return true
}

func SelectAllForBroadcast() []int {
	usersID := make([]int, 0, 10)

	rows, err := DB.Query("SELECT chatID FROM user")

	if err != nil {
		tgbotapi.NewMessage(-4801118127, "🚫Ошибка при чтении из бд данных для рассылки!")
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			continue
		}
		usersID = append(usersID, userID)
	}

	return usersID
}

func encrypt(text string, key []byte) string {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		logrus.Fatalf("Ошибка NewCipher encript")
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	// Заполняем IV случайными байтами
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		logrus.Fatal("Ошибка io.Readfull")
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext)
}

func decrypt(cryptoText string, key []byte) string {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		logrus.Fatal("Ошибка base64.DecodetoString()")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		logrus.Fatal("Ошибка newChiper decode")
	}

	if len(ciphertext) < aes.BlockSize {
		logrus.Info("ciphertext очень короткий")
		return ""
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext)
}
