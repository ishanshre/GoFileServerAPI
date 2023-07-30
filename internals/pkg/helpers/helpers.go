package helpers

import (
	"encoding/json"
	"net/http"
)

type ContextKey string

type Message struct {
	MessageStatus string `json:"status,omitempty"`
	Message       string `json:"message,omitempty"`
	Data          any    `json:"data,omitempty"`
}

func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func StatusInternalServerError(w http.ResponseWriter, message string) {
	WriteJson(w, http.StatusInternalServerError, Message{
		MessageStatus: "error",
		Message:       message,
	})
}

func StatusUnauthorized(w http.ResponseWriter, message string) {
	WriteJson(w, http.StatusUnauthorized, Message{
		MessageStatus: "error",
		Message:       message,
	})
}

func StatusBadRequest(w http.ResponseWriter, message string) {
	WriteJson(w, http.StatusBadRequest, Message{
		MessageStatus: "error",
		Message:       message,
	})
}
func StatusNotFound(w http.ResponseWriter, message string) {
	WriteJson(w, http.StatusNotFound, Message{
		MessageStatus: "error",
		Message:       message,
	})
}

func StatusCreatedData(w http.ResponseWriter, data any) {
	WriteJson(w, http.StatusCreated, Message{
		MessageStatus: "success",
		Data:          data,
	})
}

func StatusOkData(w http.ResponseWriter, data any) {
	WriteJson(w, http.StatusOK, Message{
		MessageStatus: "success",
		Data:          data,
	})
}
func StatusOk(w http.ResponseWriter, message string) {
	WriteJson(w, http.StatusOK, Message{
		MessageStatus: "success",
		Message:       message,
	})
}
