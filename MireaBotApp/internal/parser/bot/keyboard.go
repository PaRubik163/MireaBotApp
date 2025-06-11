package bot
import (
	"strings"
)

func IsGoodLogin(login string) bool {
	if !strings.Contains(login, "@edu.mirea.ru") {
		return false
	}
	return true
}

func IsGoodPassword(password string) bool {
	if len(password) < 8{
		return false
	}
	return true
}