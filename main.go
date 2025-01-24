package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type requestBody struct {
	Message string `json:"message"`
}

var task string

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var reqBody requestBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Ошибка декодирования:", http.StatusBadRequest)
			return
		}
		task = reqBody.Message
		fmt.Fprintf(w, "Получено: %s\n", task)
	} else {
		fmt.Fprintf(w, "Разрешается только метод POST\n")
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, "hello, %s\n", task)
	} else {
		fmt.Fprintf(w, "Разрешается только метод GET\n")
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/hello", PostHandler).Methods("POST")
	router.HandleFunc("/api/hello", GetHandler).Methods("GET")
	http.ListenAndServe(":8080", router)
}
