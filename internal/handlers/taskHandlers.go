package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"restapi_project/internal/database"
	"restapi_project/internal/taskService"
	"strconv"
)

type Handler struct {
	Service *taskService.TaskService
}

func NewHandler(service *taskService.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Service.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task taskService.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdTask, err := h.Service.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdTask)
}

func (h *Handler) UpdateTaskByID(w http.ResponseWriter, r *http.Request) {
	var res taskService.Task
	vars := mux.Vars(r)
	id := vars["id"]
	u64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	wd := uint(u64)
	updates := map[string]interface{}{}
	if errDecode := json.NewDecoder(r.Body).Decode(&res); errDecode != nil {
		http.Error(w, errDecode.Error(), http.StatusBadRequest)
		return
	}
	if res.Name != "" {
		updates["name"] = res.Name
	}
	updates["is_done"] = res.IsDone
	updatedTask, err := h.Service.UpdateTaskByID(wd, res)
	if errFind := database.DB.First(&updatedTask, wd).Error; errFind != nil {
		http.Error(w, errFind.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

func (h *Handler) DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	u64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Fprintln(w, err)
	}
	wd := uint(u64)
	errDeleteTask := h.Service.DeleteTaskByID(wd)
	if errDeleteTask != nil {
		http.Error(w, errDeleteTask.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
