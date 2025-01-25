package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func CreateMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var message Message
		if errDecode := json.NewDecoder(r.Body).Decode(&message); errDecode != nil {
			http.Error(w, errDecode.Error(), http.StatusBadRequest)
			return
		}
		if errCreate := DB.Create(&message).Error; errCreate != nil {
			http.Error(w, errCreate.Error(), http.StatusBadRequest)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(message)
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	var messages []Message
	if r.Method == http.MethodGet {
		if errFind := DB.Find(&messages).Error; errFind != nil {
			http.Error(w, errFind.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}

func main() {
	InitDB()
	DB.AutoMigrate(&Message{})
	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessages).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
