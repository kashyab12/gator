package legler

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, status int, payload interface{}) error {

	// Encode the payload
	if encodedPayload, encodingErr := json.Marshal(payload); encodingErr != nil {
		log.Println(encodingErr)
		return encodingErr
	} else {
		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if _, writeErr := w.Write(encodedPayload); writeErr != nil {
			log.Println(writeErr)
			return writeErr
		}
	}
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	errorJson := struct {
		Msg string `json:"error"`
	}{Msg: msg}
	return RespondWithJson(w, code, errorJson)
}
