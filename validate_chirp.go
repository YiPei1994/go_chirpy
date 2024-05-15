package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type valid struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	param := parameters{}
	err := decoder.Decode(&param)

	if err != nil {

		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return

	}
	const maxLength = 140
	if len(param.Body) > maxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned := handleClean(param.Body)

	respondWithJSON(w, http.StatusOK, valid{
		Cleaned_body: cleaned,
	})
}

func respondWithError(w http.ResponseWriter, statuscode int, val string) {
	if statuscode > 499 {
		log.Printf("Responding with 5XX error: %s", val)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, statuscode, errorResponse{
		Error: val,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)

	w.Write(dat)
}

func handleClean(msg string) string {

	splited := strings.Split(msg, " ")
	for i, v := range splited {

		if strings.ToLower(v) == "kerfuffle" || strings.ToLower(v) == "sharbert" || strings.ToLower(v) == "fornax" {
			splited[i] = "****"
		}
	}
	fmt.Println(splited)
	joined := strings.Join(splited, " ")
	return joined
}