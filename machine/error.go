// Copyright 2015, Stève Sfartz
// Licensed under the MIT License

// Utility to send json errors over HTTP
package machine

import (
	"net/http"
	"encoding/json"
)

// Error structure ala google (see Vision API)
type errorWrapper struct {
	Error `json:"error"`
}

type Error struct {
	Code int `json:"code"`
	Status string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Details *[]errorDetails `json:"details,omitempty"`
}

type errorDetails struct {
	Type string `json:"@type"`
	Links []struct {
		Description string `json:"description"`
	} `json:"links"`
}

func sendBadRequest(w http.ResponseWriter, message string) {
	sendError(w, http.StatusBadRequest, "BAD REQUEST", message)
}

func sendInternalError(w http.ResponseWriter, reason string, message string) {
	sendError(w, http.StatusInternalServerError, reason, message)
}

func sendError(w http.ResponseWriter, code int, status string, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	error := Error{
		Code: code,
		Status: status,
		Message: message,
		Details: nil,
	}
	ew := errorWrapper{error}
	enc := json.NewEncoder(w)
	enc.Encode(ew)
}


