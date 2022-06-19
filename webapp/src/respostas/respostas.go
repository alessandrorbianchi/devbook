package respostas

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorAPI struct {
	Error string `json:"error"`
}

func JSON(w http.ResponseWriter, statusCode int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if statusCode != http.StatusNoContent {
		if err := json.NewEncoder(w).Encode(dados); err != nil {
			log.Fatal(err)
		}
	}
}

func TratarStatusCodeDeErro(w http.ResponseWriter, r *http.Response) {
	var erro ErrorAPI
	json.NewDecoder(r.Body).Decode(&erro)
	JSON(w, r.StatusCode, erro)
}
