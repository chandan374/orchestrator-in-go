package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"cube/task"
)

type ErrResponse struct {
	Error string `json:"error"`
}

func (a *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var te task.TaskEvent
	err := d.Decode(&te)
	if err != nil {
		msg := fmt.Sprintf("Error decoding task event: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrResponse{
			Error: msg,
		})
		return
	}

	a.Worker.AddTask(te.Task)
	log.Printf("Task %s added to worker\n", te.Task.ID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(te)
}

func (a *Api) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a.Worker.GetTasks())
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		log.Printf("No task ID provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tId, err := uuid.Parse(taskID)
	if err != nil {
		log.Printf("Invalid task ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	taskToStop, ok := a.Worker.GetTask(tId)
	if !ok {
		log.Printf("Task %s not found", taskID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	taskCopy := *taskToStop
	taskCopy.State = task.Completed
	a.Worker.AddTask(taskCopy)

	log.Printf("Task %s added to stop the container", taskID)
	w.WriteHeader(http.StatusNoContent)
}

func (a *Api) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a.Worker.Stats)
}

