package response

import (
	"encoding/json"
	"net/http"
)

type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func OK(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, Body{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

func Fail(w http.ResponseWriter, httpStatus int, code int, message string) {
	writeJSON(w, httpStatus, Body{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func writeJSON(w http.ResponseWriter, httpStatus int, body Body) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(body)
}
