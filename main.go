package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}


func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/api/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/validate_chirp",handlerChirpsValidate)
	mux.HandleFunc("POST /api/chirps",handleChirps)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func handleChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	fmt.Println(r.Body)
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {

		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return

	}
	const maxLength = 140
	if len(params.Body) > maxLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	handleSave(w,http.StatusOK,params.Body)

}

func handleSave(w http.ResponseWriter, code int, payload string) {
	type newResponse struct {
		Id int `json:"id"`
		Body string `json:"body"`
	}

	res := newResponse{
		Id: 1,
		Body: payload,
	}
	encode,err := json.Marshal(res)

	if err != nil {
		fmt.Println("something went wrong")
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(code)
	w.Write(encode)
}

