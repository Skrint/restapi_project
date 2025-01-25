package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func CreateMessages(w http.ResponseWriter, r *http.Request) {
	var message Message
	if r.Method == http.MethodPost {
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

func PatchMessages(w http.ResponseWriter, r *http.Request) {
	var res Message
	updates := map[string]interface{}{}
	if r.Method == http.MethodPatch {
		if errDecode := json.NewDecoder(r.Body).Decode(&res); errDecode != nil {
			http.Error(w, errDecode.Error(), http.StatusBadRequest)
			return
		}
		if res.Task != "" {
			updates["task"] = res.Task
		}
		if res.IsDone == false || res.IsDone == true {
			updates["is_done"] = res.IsDone
		}
		if errPatch := DB.Model(&Message{}).Where("id = ?", res.ID).Updates(updates).Error; errPatch != nil {
			http.Error(w, errPatch.Error(), http.StatusBadRequest)
			return
		}
	}
}

func DeleteMessages(w http.ResponseWriter, r *http.Request) {
	var res Message
	if r.Method == http.MethodDelete {
		if errDecode := json.NewDecoder(r.Body).Decode(&res); errDecode != nil {
			http.Error(w, errDecode.Error(), http.StatusBadRequest)
			return
		}
		if errDelete := DB.Delete(&Message{}, res.ID).Error; errDelete != nil {
			http.Error(w, errDelete.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "Удалена запись с id", res.ID)
	}
}

func main() {
	InitDB()
	DB.AutoMigrate(&Message{})
	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessages).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	router.HandleFunc("/api/messages", PatchMessages).Methods("PATCH")
	router.HandleFunc("/api/messages", DeleteMessages).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
