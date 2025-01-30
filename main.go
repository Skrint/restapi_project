package main

import (
	"encoding/json"
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
	vars := mux.Vars(r)
	id := vars["id"]
	updates := map[string]interface{}{}
	if r.Method == http.MethodPatch {
		if errDecode := json.NewDecoder(r.Body).Decode(&res); errDecode != nil {
			http.Error(w, errDecode.Error(), http.StatusBadRequest)
			return
		}
		updates["task"] = res.Task
		updates["is_done"] = res.IsDone
		if errPatch := DB.Model(&Message{}).Where("id = ?", id).Updates(updates).Error; errPatch != nil {
			http.Error(w, errPatch.Error(), http.StatusBadRequest)
			return
		}
		if errFind := DB.Find(&res).Error; errFind != nil {
			http.Error(w, errFind.Error(), http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func DeleteMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if r.Method == http.MethodDelete {
		if errDelete := DB.Delete(&Message{}, id).Error; errDelete != nil {
			http.Error(w, errDelete.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {
	InitDB()
	DB.AutoMigrate(&Message{})
	router := mux.NewRouter()
	router.HandleFunc("/api/messages", CreateMessages).Methods("POST")
	router.HandleFunc("/api/messages", GetMessages).Methods("GET")
	router.HandleFunc("/api/messages/{id}", PatchMessages).Methods("PATCH")
	router.HandleFunc("/api/messages/{id}", DeleteMessages).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
