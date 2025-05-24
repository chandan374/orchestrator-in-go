package manager

import(
	"cube/task"
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Manager struct {
	Pending queue.Queue
	TaskDb map[string][]*task.Task
	EventDb map[string][]*task.TaskEvent
	Workers []string
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
}

func (m *Manager) SelectWorker() {
	fmt.Println("I will select a worker")
}

func (m *Manager) UpdateTask() {
	fmt.Println("I will update a task")
}

func (m *Manager) SendWork() {
	fmt.Println("I will send work")
}

