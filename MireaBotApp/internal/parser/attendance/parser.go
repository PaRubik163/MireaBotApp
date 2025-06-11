package attendance

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"log"

	"github.com/go-resty/resty/v2"
	"google.golang.org/protobuf/proto"
	uuid "mireabot/internal/parser/attendance/proto/GetAvailableVisitingLogsOfStudent"
	respScore "mireabot/internal/parser/attendance/proto/GetLearnRatingScoreReportForStudentInVisitingLog"
)

func callGetLearnRatingScoreReportForStudentInVisiting(client *resty.Client) []byte {

	UUID, err := callGetAvailableVisitingLogsOfStudent(client)

	if err != nil {
		log.Fatalf("Ошибка в вызове gRPC-WEB функции, callGetLearnRatingScoreReportForStudentInVisitingLog()", err)
	}

	req := &respScore.GetScoreAndVisitngRequest{
		Id: UUID,
	}

	// сериализация protobuf
	data, err := proto.Marshal(req)
	if err != nil {
		log.Fatalf("marshal GetMeInfoRequest: %v", err)
	}

	// gRPC-web frame: 1 байт флага и 4 байта длины
	var buf bytes.Buffer
	buf.WriteByte(0x00)
	binary.Write(&buf, binary.BigEndian, uint32(len(data)))
	buf.Write(data)

	// grpc-web-text требует base64
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	resp, err := client.R().
		SetHeader("Content-Type", "application/grpc-web-text").
		SetHeader("X-Grpc-Web", "1").
		SetHeader("Origin", "https://attendance-app.mirea.ru").
		SetHeader("Referer", "https://attendance-app.mirea.ru/").
		SetBody(encoded).
		Post("https://attendance.mirea.ru/rtu_tc.attendance.api.LearnRatingScoreService/GetLearnRatingScoreReportForStudentInVisitingLog")

	if err != nil {
		log.Fatal(err)
	}

	return resp.Body()
}

// gRPC-WEB запрос, чтобы достать id студента для следующих запросов (вспомогательная функция)
func callGetAvailableVisitingLogsOfStudent(client *resty.Client) (string, error) {
	req := &uuid.GetAvailableVisitingLogsOfStudentRequest{}

	// сериализация protobuf
	data, err := proto.Marshal(req)
	if err != nil {
		return "", errors.New("Ошибка кодирования")
	}

	// gRPC-web frame: 1 байт флага и 4 байта длины
	var buf bytes.Buffer
	buf.WriteByte(0x00)
	binary.Write(&buf, binary.BigEndian, uint32(len(data)))
	buf.Write(data)

	// grpc-web-text требует base64
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	// отправка запроса
	resp, err := client.R().
		SetHeader("Content-Type", "application/grpc-web-text").
		SetHeader("X-Grpc-Web", "1").
		SetHeader("Origin", "https://attendance-app.mirea.ru").
		SetHeader("Referer", "https://attendance-app.mirea.ru/").
		SetHeader("User-Agent", "grpc-web-javascript/0.1").
		SetHeader("X-Requested-With", "XMLHttpRequest").
		SetHeader("Accept-Encoding", "gzip, deflate, br, zstd").
		SetHeader("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7").
		SetBody(encoded).
		Post("https://attendance.mirea.ru/rtu_tc.attendance.api.VisitingLogService/GetAvailableVisitingLogsOfStudent")

	if err != nil {
		return "", errors.New("Ошибка gRPC-WEB запроса для взятия ID")
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	UUID := resp.String()[13:49]
	if len(UUID) < 36 {
		return "", errors.New("ID >36 МБ Ошибки в Логине или Пароле")
	}
	return UUID, nil
}
