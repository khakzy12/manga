
---

# MangaHub - Multi-Protocol Manga Tracking System

MangaHub is a distributed manga management system designed to demonstrate the implementation of various network protocols using the **Go** programming language. The system allows users to track reading progress, receive notifications, and interact via chat across multiple communication layers.

## 📋 Table of Contents

* [Overview](https://www.google.com/search?q=%23overview)
* [System Architecture](https://www.google.com/search?q=%23system-architecture)
* [Prerequisites (Windows)](https://www.google.com/search?q=%23prerequisites-windows)
* [Installation](https://www.google.com/search?q=%23installation)
* [Project Structure](https://www.google.com/search?q=%23project-structure)
* [Building and Running](https://www.google.com/search?q=%23building-and-running)
* [Testing Protocols](https://www.google.com/search?q=%23testing-protocols)
* [Database Management](https://www.google.com/search?q=%23database-management)
* [Troubleshooting](https://www.google.com/search?q=%23troubleshooting)

---

## Overview

MangaHub consists of five core server components and a unified CLI application:

| Component | Port | Protocol | Description |
| --- | --- | --- | --- |
| **API Server** | 8080 | **HTTP/REST** | Handles authentication, library management, and progress updates. |
| **TCP Sync Server** | 9090 | **TCP** | Provides real-time progress synchronization across devices. |
| **UDP Server** | 9091 | **UDP** | Broadcasts system-wide notifications for new manga releases. |
| **gRPC Server** | 9092 | **gRPC** | High-performance internal service for manga data operations. |
| **WebSocket Server** | 9093 | **WS** | Facilitates real-time community chat functionality. |

---

## Prerequisites (Windows)

### 1. Go Programming Language

* Download the Windows installer from [golang.org](https://golang.org/dl/).
* Ensure Go version **1.21 or later** is installed.
* Verify by running `go version` in PowerShell.

### 2. Protocol Buffers (protoc)

* Download the Windows zip (e.g., `protoc-xx.x-win64.zip`) from [Protobuf Releases](https://github.com/protocolbuffers/protobuf/releases).
* Extract to a folder (e.g., `C:\protoc`) and add the `bin` directory to your System **Environment Variables (Path)**.

### 3. Go Protobuf Plugins

Run the following commands in PowerShell to install necessary plugins:

```powershell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

```

*Note: Ensure `%USERPROFILE%\go\bin` is also added to your Windows **Path**.*

---

## Installation

1. **Clone the repository** and navigate to the project root:
```powershell
cd mangahub

```


2. **Install dependencies**:
```powershell
go mod tidy

```


3. **Generate Protobuf code**:
```powershell
cd pkg\proto
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       manga.proto
cd ..\..

```



---

## Project Structure

```text
mangahub/
├── cmd/
│   ├── api-server/         # HTTP REST API server
│   ├── tcp-server/         # TCP sync server
│   ├── udp-server/         # UDP notification server
│   ├── grpc-server/        # gRPC server
│   └── cli-app/            # Command-line interface
├── internal/
│   ├── api/                # REST handlers (Auth, Library, Progress)
│   ├── database/           # SQLite initialization & Seeding
│   ├── model/              # Shared data structures
│   ├── tcp/                # TCP core logic & broadcasting
│   └── udp/                # UDP packet handling
└── mangahub.db             # Persistent SQLite database

```

---

## Building and Running

### Build Executables

To create Windows `.exe` files for all components:

```powershell
go build -o api-server.exe cmd\api-server\main.go
go build -o tcp-server.exe cmd\tcp-server\main.go
go build -o udp-server.exe cmd\udp-server\main.go

```

### Run Servers

Open separate PowerShell windows for each:

* **API:** `go run cmd\api-server\main.go`
* **TCP:** `go run cmd\tcp-server\main.go`
* **UDP:** `go run cmd\udp-server\main.go`

---

## Testing Protocols

### 1. HTTP API (REST)

Use the native Windows `curl.exe` to register or login:

```powershell
# Login to receive a JWT Token
curl.exe -X POST http://localhost:8080/auth/login `
  -H "Content-Type: application/json" `
  -d '{\"username\":\"your_user\",\"password\":\"your_pass\"}'

```

### 2. TCP Real-time Sync

Simulate a listening device using a PowerShell script:

```powershell
$client = New-Object System.Net.Sockets.TCPClient("localhost", 9090); `
$stream = $client.GetStream(); $writer = New-Object System.IO.StreamWriter($stream); `
$writer.AutoFlush = $true; $writer.WriteLine('{"user_id":"your_user_id"}'); `
$reader = New-Object System.IO.StreamReader($stream); `
while($client.Connected) { $line = $reader.ReadLine(); if($line) { Write-Host "SYNC RECEIVED: $line" -ForegroundColor Cyan } }

```

### 3. UDP Notifications

Blast a global notification to the server:

```powershell
$udp = New-Object System.Net.Sockets.UdpClient; `
$udp.Connect("127.0.0.1", 9091); `
$data = [System.Text.Encoding]::ASCII.GetBytes("New Release: One Piece 1111!"); `
$udp.Send($data, $data.Length); $udp.Close()

```

---

## Database Management

The system uses **SQLite** for persistence.

* **File:** `mangahub.db`
* **Seeding:** The database is automatically seeded with sample manga (e.g., One Piece, Naruto) upon the first run of the API server.
* **Recommended Tool:** Use [DBeaver](https://dbeaver.io/) for viewing and managing tables.

---

## Troubleshooting

* **Port Conflict:** If a port is occupied, use `netstat -ano | findstr :8080` to find the Process ID (PID), then run `taskkill /F /PID <PID>` to stop it.
* **DB Locked:** Ensure no other application (like DBeaver) has an active write-lock on `mangahub.db` when starting the servers.
* **Protoc PATH:** If `protoc` is not recognized, restart your terminal after updating Environment Variables.

---

*Developed by: Nguyễn Trầm Gia Hưng - ITITIU23007*
