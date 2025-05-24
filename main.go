package main

import (
	"log"
	"os"
	"strconv"
	"time"
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"

	"cube/task"
	"cube/worker"
)

func main() {
	host := os.Getenv("CUBE_HOST")
	port, _ := strconv.Atoi(os.Getenv("CUBE_PORT"))

	fmt.Println("Starting worker")

	w := worker.Worker{
		Queue: *queue.New(),
		Db: make(map[uuid.UUID]*task.Task),
	}
	api := worker.Api{
		Address: host,
		Port: port, 
		Worker: &w,
	}

	go runTasks(&w)
	api.Start()
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() > 0 {
			result := w.RunTask() 
			if result.Error != nil {
				log.Printf("Error running task: %v", result.Error)
			}
		} else {
			log.Printf("No tasks to run")
		}
		log.Printf("Sleeping for 10 seconds")
		time.Sleep(10 * time.Second)
	}
}
