# Go-websocket

Multi-room real-time messaging server built with **Go** and the **Fiber** framework. This project demonstrates high-concurrency patterns, utilizing Go’s native primitives to manage independent chat rooms, thread-safe state, and asynchronous message broadcasting.

## Features

- **Multi-Room Architecture:** Support for dynamic creation and management of separate chat rooms via URL queries.
- **Concurrent Message Orchestration:** A centralized `select` loop for each room to handle client joins, leaves, and message forwarding without race conditions.
- **Thread-Safe Room Management:** Global room registry protected by `sync.Mutex` to ensure memory safety across concurrent requests.
- **Bi-directional Streaming:** Full-duplex communication using the WebSocket protocol for low-latency updates.
- **JSON Structured Messaging:** Automated serialization of message metadata including sender names and payloads.

## Technical Highlights

- **Synchronization Primitives:** Implemented `sync.Mutex` to safely manage the global `rooms` map, preventing "concurrent map write" panics.
- **Channel-Based Communication:** Utilized unbuffered channels (`join`, `leave`, `forward`) to synchronize state changes and broadcast data within chat rooms.
- **Dedicated Read/Write Loops:** Spawned independent goroutines for each client’s write operations, ensuring that slow network I/O does not block the entire room's message flow.
- **Fiber Framework:** Leveraged Fiber v2 for high-performance HTTP routing and seamless WebSocket upgrading.
- **Resource Management:** Implemented graceful cleanup to close channels and remove client references upon disconnection to prevent memory leaks.

## Usage

To run the server locally, ensure you have Go installed:

```bash
# Clone the repository
git clone github.com/Kazimlyc/go-websocket
cd go-websocket

# Install dependencies
go mod download

# Run the server
go run main.go --addr=":8080"
