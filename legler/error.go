package legler

import "net/http"

func GetErrorLegler(w http.ResponseWriter, r *http.Request) {
	// TODO: what's a good way to handle the error from RespondWithError lol
	_ = RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
