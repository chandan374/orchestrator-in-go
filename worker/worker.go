package worker

import (
	"fmt"
	"time"
	"log"
	"github.com/google/uuid"
	"github.com/golang-collections/collections/queue"

	"cube/task"
)

type Worker struct {
	Name string
	Queue queue.Queue
	Db map[uuid.UUID]*task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	fmt.Println("I will collect stats")
}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(t)
}

func (w *Worker) StartTask(t task.Task) task.DockerResult {
	t.StartTime = time.Now().UTC()
	config := task.NewConfig(&t)
	d := task.NewDocker(config)
	result := d.Run()
	
	if result.Error != nil {
		fmt.Println("ERror in running task")
		t.State = task.Failed
		w.Db[t.ID] = &t
		return result
	}
	
	t.ContainerID = result.ContainerID
	t.State = task.Running
	w.Db[t.ID] = &t
	log.Printf("Task %s started", t.ID)
	return result
}

func (w *Worker) RunTask() task.DockerResult {
	t := w.Queue.Dequeue()
	if t == nil {
		fmt.Println("No task to run")
		return task.DockerResult{Error: fmt.Errorf("no task to run")}
	}

	taskQueued := t.(task.Task)
	taskPersisted := w.Db[taskQueued.ID]
	if taskPersisted == nil {
		taskPersisted = &taskQueued
		w.Db[taskQueued.ID] = taskPersisted
	}

	var result task.DockerResult
	fmt.Printf("Task %s is in state %s\n", taskQueued.State, taskPersisted.State)
	if task.ValidateStateTransition(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case task.Scheduled:
			result = w.StartTask(taskQueued)
		case task.Completed:
			result = w.StopTask(taskQueued)
		default:
			result = task.DockerResult{Error: fmt.Errorf("invalid state transition")}
		}
	} else {
		err := fmt.Errorf("invalid state transition")
		result = task.DockerResult{Error: err}
	}

	return result
}

func (w *Worker) StopTask(t task.Task) task.DockerResult {
	config := task.NewConfig(&t)
	d := task.NewDocker(config)

	result := d.Stop(t.ContainerID)
	if result.Error != nil {
		fmt.Printf("Error stopping container: %v\n", result.Error)
		return result
	}

	t.FinishTime = time.Now().UTC()
	t.State = task.Completed
	w.Db[t.ID] = &t
	log.Printf("Task %s completed", t.ID)
	return result
}

