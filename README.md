
---

# MangaHub - Network Programming Project (Windows Version)

A comprehensive manga management system demonstrating multiple network protocols: HTTP REST API, TCP, UDP, WebSocket, and gRPC.

## 📋 Table of Contents

* [Overview](https://www.google.com/search?q=%23overview)
* [Prerequisites](https://www.google.com/search?q=%23prerequisites)
* [Installation](https://www.google.com/search?q=%23installation)
* [Project Structure](https://www.google.com/search?q=%23project-structure)
* [Building the Project](https://www.google.com/search?q=%23building-the-project)
* [Running the Servers](https://www.google.com/search?q=%23running-the-servers)
* [Testing](https://www.google.com/search?q=%23testing)
* [Database Management](https://www.google.com/search?q=%23database-management)
* [Troubleshooting](https://www.google.com/search?q=%23troubleshooting)

---

## Overview

MangaHub consists of 5 server components and a CLI application:

| Component | Port | Protocol | Description |
| --- | --- | --- | --- |
| **API Server** | 8080 | HTTP/REST | User authentication and manga management |
| **TCP Server** | 9090 | TCP | Real-time progress synchronization |
| **UDP Server** | 9091 | UDP | Broadcast notifications for new releases |
| **gRPC Server** | 9092 | gRPC | Manga operations via Protocol Buffers |
| **WebSocket Server** | 9093 | WebSocket | Real-time chat functionality |
| **CLI App** | - | - | Unified command-line interface |

---

## Prerequisites

### 1. Go (version 1.21 or later)

* Download the installer from [golang.org/dl](https://golang.org/dl/).
* Run the `.msi` file and follow the instructions.
* Verify in PowerShell: `go version`

### 2. Protocol Buffers Compiler (protoc)

* Download `protoc-xx.x-win64.zip` from the [Protobuf Releases page](https://github.com/protocolbuffers/protobuf/releases).
* Extract it to a folder (e.g., `C:\protoc`).
* Add `C:\protoc\bin` to your system **Environment Variables (Path)**.
* Verify in PowerShell: `protoc --version`

### 3. Go Protobuf Plugins

Run these commands in PowerShell:

```powershell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

```

### 4. Setting Up PATH for Go Plugins

The plugins are installed in `%USERPROFILE%\go\bin`. You must add this to your **Path**:

1. Search for "Edit the system environment variables" in the Start menu.
2. Click **Environment Variables**.
3. Under "User variables", select **Path** and click **Edit**.
4. Click **New** and paste: `%USERPROFILE%\go\bin`
5. Click OK and **restart your terminal**.

### 5. Network Tools

* **curl**: Pre-installed on Windows 10 (version 1803+) and Windows 11.
* **ncat**: Included with [Nmap](https://nmap.org/download.html) (recommended for TCP/UDP testing) or use the built-in CLI app.

---

## Installation

1. Open PowerShell and navigate to the project directory:

```powershell
cd C:\Users\YourName\Documents\mangahub

```

2. Install Go dependencies:

```powershell
go mod tidy

```

3. Generate Protocol Buffer code:

```powershell
cd pkg\proto
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       manga.proto
cd ..\..

```

---

## Building the Project

### Build CLI Application

```powershell
go build -o mangahub.exe cmd\cli-app\main.go

```

### Build Individual Servers

```powershell
# API Server
go build -o api-server.exe cmd\api-server\main.go

# TCP Server
go build -o tcp-server.exe cmd\tcp-server\main.go

# UDP Server
go build -o udp-server.exe cmd\udp-server\main.go

# gRPC Server
go build -o grpc-server.exe cmd\grpc-server\main.go

# WebSocket Server
go build -o websocket-server.exe cmd\websocket-server\main.go

```

---

## Running the Servers

**Note:** Run each server in a **separate** PowerShell window or terminal tab.

1. **API Server:** `go run cmd\api-server\main.go` (Port 8080)
2. **TCP Server:** `go run cmd\tcp-server\main.go` (Port 9090)
3. **UDP Server:** `go run cmd\udp-server\main.go` (Port 9091)
4. **gRPC Server:** `go run cmd\grpc-server\main.go` (Port 9092)
5. **WebSocket Server:** `go run cmd\websocket-server\main.go` (Port 9093)

---

## Testing

### 1. Test HTTP API Server

**Register:**

```powershell
curl.exe -X POST http://localhost:8080/auth/register `
  -H "Content-Type: application/json" `
  -d '{\"username\":\"testuser\",\"password\":\"testpass123\"}'

```

**Login:**

```powershell
curl.exe -X POST http://localhost:8080/auth/login `
  -H "Content-Type: application/json" `
  -d '{\"username\":\"testuser\",\"password\":\"testpass123\"}'

```

### 2. Test CLI Application

```powershell
.\mangahub.exe auth login
.\mangahub.exe manga search "One Piece"
.\mangahub.exe chat join

```

### 3. Test TCP Server

Using PowerShell (if you don't have ncat):

```powershell
# Connect to TCP
$client = New-Object System.Net.Sockets.TcpClient("localhost", 9090)
$stream = $client.GetStream()
$writer = New-Object System.IO.StreamWriter($stream)
$writer.AutoFlush = $true
$writer.WriteLine('{"type":"auth","token":"YOUR_TOKEN"}')

```

### 4. Test UDP Server (via PowerShell)

```powershell
$udpclient = New-Object System.Net.Sockets.UdpClient
$udpclient.Connect("localhost", 9091)
$bytes = [System.Text.Encoding]::ASCII.GetBytes("SUBSCRIBE")
$udpclient.Send($bytes, $bytes.Length)

```

---

## Database Management

### Using DBeaver

1. Download from [dbeaver.io](https://dbeaver.io/).
2. Create a new **SQLite** connection.
3. Path: `C:\Users\YourName\Documents\mangahub\mangahub.db`

### Database Location Environment Variable

In PowerShell:

```powershell
$env:DB_PATH="C:\path\to\mangahub.db"
go run cmd\api-server\main.go

```

---

## Troubleshooting

* **Port Already in Use**: Run `netstat -ano | findstr :8080` to find the Process ID (PID), then run `taskkill /F /PID <PID>` to stop it.
* **Script Execution Blocked**: If you cannot run Go commands, run PowerShell as Administrator and type: `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned`.
* **Database Locked**: Ensure you close DBeaver or any other database browser before running the server, as SQLite only allows one write-access process at a time.
* **Path Issues**: Always use backslashes `\` for file paths in Windows CMD/PowerShell, or use double backslashes `\\` inside Go environment variables.

---

**Quick Start Summary (Windows PowerShell):**

```powershell
go mod tidy
cd pkg\proto; protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative manga.proto; cd ..\..
go build -o mangahub.exe cmd\cli-app\main.go
# Start API Server
go run cmd\api-server\main.go

```
