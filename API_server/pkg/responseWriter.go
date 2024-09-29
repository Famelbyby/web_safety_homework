package pkg

import (
	"encoding/json"
	"main/domain"
	"net/http"
	"strconv"
)

func WriteHTTPRecord(w http.ResponseWriter, result domain.HTTPEntity) error {
	w.Write([]byte("ID запроса: " + strconv.Itoa(result.ID) + "\n"))
	w.Write([]byte("Запрос:\n"))

	err := json.NewEncoder(w).Encode(result.Request)

	if err != nil {
		return err
	}

	w.Write([]byte("Ответ:\n"))

	err = json.NewEncoder(w).Encode(result.Answer)

	if err != nil {
		return err
	}

	w.Write([]byte("\n"))

	return nil
}

func WriteHTTPSRecord(w http.ResponseWriter, result domain.HTTPSEntity) error {
	w.Write([]byte("ID запроса: " + strconv.Itoa(result.ID) + "\nЗапрос:\n"))
	w.Write([]byte(result.ClientRequest + "\n"))
	w.Write([]byte("Ответ:\n" + result.AnswerData + "\n"))

	return nil
}
