package legler

import "net/http"

func GetReadinessLegler(w http.ResponseWriter, _ *http.Request) {
	statusResponse := struct {
		Status string `json:"status"`
	}{Status: "ok"}
	if responseError := RespondWithJson(w, http.StatusOK, statusResponse); responseError != nil {
		// TODO: what's a good way to handle the error from RespondWithError lol
		_ = RespondWithError(w, http.StatusInternalServerError, responseError.Error())
	}
}
