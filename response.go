package main

import (
	"net/http"
	"encoding/json"

)

func respondWithJSON (w http.ResponseWriter, code int, payload interface{}) error {
	res, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "apllication/json")
	w.WriteHeader(code)
	w.Write(res)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
    return respondWithJSON(w, code, map[string]string{"error": msg})
}

