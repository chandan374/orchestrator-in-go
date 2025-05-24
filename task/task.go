package task

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int 

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

var StateTransitionMap = map[State][]State{
	Pending: []State{Scheduled},
	Scheduled: []State{Running, Failed, Scheduled},
	Running: []State{Running, Completed, Failed},
	Completed: []State{},
	Failed: []State{},
}

func Contains(states []State, state State) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}
	return false
}

func ValidateStateTransition(src State, dst State) bool {
	states := StateTransitionMap[src]
	return Contains(states, dst)
}

type Task struct {
	ID uuid.UUID
	Name string
	State State
	Image string
	Memory int
	Disk int
	ExposedPorts nat.PortSet
	PortBindings map[string]string
	RestartPolicy string
	StartTime time.Time
	EndTime time.Time
	ContainerID string
	FinishTime time.Time
}

type TaskEvent struct {
	ID uuid.UUID
	State State
	Timestamp time.Time
	Task Task
}

type Config struct {
	Name string
	AttachStdin bool
	AttachStdout bool
	AttachStderr bool
	ExposedPorts nat.PortSet
	Cmd []string 
	Image string
	Cpu float64
	Memory int64
	Disk int64
	Env []string
	RestartPolicy string
	Runtime struct {
		ContainerID string
	}
}

func NewDocker(cfg Config) *Docker {
	client, _ := client.NewClientWithOpts(client.FromEnv)
	return &Docker{
		Client: client,
		Config: cfg,
	}
}

func NewConfig(t *Task) Config {
	return Config{
		Name: t.Name,
		Image: t.Image,
		Cmd: []string{"/bin/sh", "-c", "echo 'Hello, World!'"},
	}
}

type Docker struct {
	Client *client.Client
	Config Config
}

type DockerResult struct {
	Error error
	Action string
	ContainerID string
	Result string
}

func (d *Docker) Run() DockerResult {
	ctx := context.Background()
	reader, err := d.Client.ImagePull(
		ctx,
		d.Config.Image,
		types.ImagePullOptions{},
	)
	if err != nil {
		log.Printf("Error pulling image: %v", err)
		return DockerResult{Error: err}
	}
	io.Copy(os.Stdout, reader)

	rp := container.RestartPolicy{
		Name: d.Config.RestartPolicy,
	}

	r := container.Resources{
		Memory: d.Config.Memory,
		NanoCPUs: int64(d.Config.Cpu * 1000000000),
	}

	cc := container.Config{
		Image: d.Config.Image,
		Tty: false,
		Env: d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}

	hc := container.HostConfig{
		RestartPolicy: rp,
		Resources: r,
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container: %v", err)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Printf("Error starting container: %v", err)
		return DockerResult{Error: err}
	}

	d.Config.Runtime.ContainerID = resp.ID

	out, err := d.Client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
	})

	if err != nil {
		log.Printf("Error getting container logs: %v", err)
		return DockerResult{Error: err}
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return DockerResult{
		ContainerID: resp.ID,
		Action: "start",
		Result: "success",
	}
}

func (d *Docker) Stop(id string) DockerResult {
	log.Printf("Stopping container: %s", d.Config.Runtime.ContainerID)
	ctx := context.Background()
	err := d.Client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		log.Printf("Error stopping container: %v", err)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks: false,
		Force: true,
	})

	if err != nil {
		log.Printf("Error removing container: %v", err)
		return DockerResult{Error: err}
	}

	return DockerResult{
		ContainerID: id,
		Action: "stop",
		Result: "success",
	}
}
