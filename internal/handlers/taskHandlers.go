package handlers

import (
	"context"
	"restapi_project/internal/taskService"
	"restapi_project/internal/web/tasks"
)

type Handler struct {
	Service *taskService.TaskService
}

func (h *Handler) GetUsersUserIdTasks(ctx context.Context, request tasks.GetUsersUserIdTasksRequestObject) (tasks.GetUsersUserIdTasksResponseObject, error) {
	userID := request.UserId
	tasksByUserID, err := h.Service.GetTasksByUserID(uint(userID))
	if err != nil {
		return nil, err
	}

	response := tasks.GetUsersUserIdTasks200JSONResponse{}

	for _, tsk := range tasksByUserID {
		task := tasks.Task{
			Id:        &tsk.ID,
			IsDone:    &tsk.IsDone,
			Task:      &tsk.Task,
			CreatedAt: &tsk.CreatedAt,
			UpdatedAt: &tsk.UpdatedAt,
		}
		response = append(response, task)
	}

	return response, nil
}

func (h *Handler) DeleteTasksId(ctx context.Context, request tasks.DeleteTasksIdRequestObject) (tasks.DeleteTasksIdResponseObject, error) {
	taskID := request.Id
	err := h.Service.DeleteTaskByID(uint(taskID))
	if err != nil {
		return nil, err
	}
	response := tasks.DeleteTasksId204Response{}
	return response, nil
}

func (h *Handler) PatchTasksId(ctx context.Context, request tasks.PatchTasksIdRequestObject) (tasks.PatchTasksIdResponseObject, error) {
	taskID := request.Id
	body := request.Body

	taskToUpdate := taskService.Task{
		Task:   *body.Task,
		IsDone: *body.IsDone,
	}

	updatedTask, err := h.Service.UpdateTaskByID(uint(taskID), taskToUpdate)
	if err != nil {
		return nil, err
	}

	response := tasks.PatchTasksId200JSONResponse{
		Id:        &updatedTask.ID,
		Task:      &updatedTask.Task,
		IsDone:    &updatedTask.IsDone,
		CreatedAt: &updatedTask.CreatedAt,
		UpdatedAt: &updatedTask.UpdatedAt,
		UserId:    &updatedTask.UserID,
	}
	return response, nil
}

func (h *Handler) GetTasks(ctx context.Context, request tasks.GetTasksRequestObject) (tasks.GetTasksResponseObject, error) {
	allTasks, err := h.Service.GetAllTasks()
	if err != nil {
		return nil, err
	}

	response := tasks.GetTasks200JSONResponse{}

	for _, tsk := range allTasks {
		task := tasks.Task{
			Id:        &tsk.ID,
			IsDone:    &tsk.IsDone,
			Task:      &tsk.Task,
			CreatedAt: &tsk.CreatedAt,
			UpdatedAt: &tsk.UpdatedAt,
			UserId:    &tsk.UserID,
		}
		response = append(response, task)
	}
	return response, nil
}

func (h *Handler) PostTasks(ctx context.Context, request tasks.PostTasksRequestObject) (tasks.PostTasksResponseObject, error) {
	taskRequest := request.Body
	taskToCreate := taskService.Task{
		Task:   *taskRequest.Task,
		IsDone: *taskRequest.IsDone,
		UserID: *taskRequest.UserId,
	}
	createdTask, err := h.Service.CreateTask(taskToCreate)

	if err != nil {
		return nil, err
	}

	response := tasks.PostTasks201JSONResponse{
		Id:     &createdTask.ID,
		Task:   &createdTask.Task,
		IsDone: &createdTask.IsDone,
		UserId: &createdTask.UserID,
	}

	return response, nil
}

func NewHandler(service *taskService.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}
