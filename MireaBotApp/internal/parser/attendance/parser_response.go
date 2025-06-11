package attendance

import (
	"bytes"
	"encoding/binary"
	"github.com/go-resty/resty/v2"
	"math"
)

// Так как я не понимаю почему я не могу декодировать бинарный формат в строку и анмаршалить в структуру, приходится костылить
func ParseGrpcResponse(client *resty.Client) ([]map[string]interface{}, bool) {

	data := callGetLearnRatingScoreReportForStudentInVisiting(client)
	pos := 5 // Пропускаем gRPC заголовок
	var subjects []map[string]interface{}

	for pos < len(data) {
		// Ищем начало названия предмета (русские буквы)
		if data[pos] == 0xD0 || data[pos] == 0xD1 {
			start := pos
			// Ищем конец названия (маркер 0x12)
			for pos < len(data) && data[pos] != 0x12 {
				pos++
			}
			nameBytes := data[start:pos]
			// Удаляем некорректные UTF-8 последовательности
			name := string(bytes.ToValidUTF8(nameBytes, []byte("?")))

			// Ищем баллы (каждый блок чисел начинается с 0x12)
			var scores []float64
			for pos < len(data)-8 {
				if data[pos] == 0x12 && data[pos+1] == 0x09 { // Маркер числа
					pos += 2
					if pos+8 <= len(data) {
						scoreBytes := data[pos : pos+8]
						score := binary.LittleEndian.Uint64(scoreBytes)
						scores = append(scores, roundFloat(float64(math.Float64frombits(score)), 1))
						pos += 8
					}
				} else if data[pos] == 0xD0 || data[pos] == 0xD1 { // Сдедующий предмет
					break
				} else {
					pos++
				}
			}

			// Первые два балла: текущий контроль и посещения
			if len(scores) >= 2 {
				subject := map[string]interface{}{
					"name":            name,
					"current_control": scores[0],
					"attendance":      scores[1],
				}
				subjects = append(subjects, subject)
			}
		} else {
			pos++
		}
	}

	return subjects, true
}

// Вспомогательная функция для округления чисел
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
