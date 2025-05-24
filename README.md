# Cube - A Container Orchestrator in Go

This project is a learning implementation of a container orchestrator, similar to Kubernetes but simplified. It's built in Go and helps understand the fundamentals of container orchestration.

## Project Overview

Cube manages containers across a distributed system with the following components:

- **Task Management**: Handles container lifecycle (create, run, stop)
- **Worker**: Manages tasks and containers on individual nodes
- **Stats Collection**: Monitors system resources (CPU, Memory, Disk)
- **API**: RESTful interface for task management

## Architecture

### Components

1. **Task**
   - Represents a container workload
   - Manages state transitions (Pending → Scheduled → Running → Completed)
   - Handles Docker operations

2. **Worker**
   - Executes tasks
   - Manages local container lifecycle
   - Collects system statistics
   - Provides REST API endpoints

### API Endpoints

- POST `/tasks` - Create a new task
- GET `/tasks` - List all tasks
- GET `/tasks/{id}` - Get task details
- DELETE `/tasks/{id}` - Stop a task
- GET `/stats` - Get system statistics

## Getting Started

### Prerequisites

- Go 1.19+
- Docker
- Linux-based system (for stats collection)

### Environment Variables

```bash
CUBE_HOST=localhost  # API host
CUBE_PORT=5555      # API port
```

### Running the Project

```bash
# Start the worker
go run main.go
```

## Learning Goals

This project demonstrates:
- Container orchestration concepts
- Go concurrency patterns
- System resource monitoring
- RESTful API design
- Docker API integration
- State management in distributed systems

## Current Status

This is a work in progress, being built as a learning project to understand container orchestration concepts.

## Next Steps

- [ ] Add cluster management
- [ ] Implement scheduling algorithms
- [ ] Add service discovery
- [ ] Implement health checks
- [ ] Add networking features
- [ ] Implement volume management

## License

MIT
