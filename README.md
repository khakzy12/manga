Here is the full **README.md** content, optimized for **Windows** and written entirely in **English**.

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
* [Testing on Windows](https://www.google.com/search?q=%23testing-on-windows)
* [Database Management](https://www.google.com/search?q=%23database-management)
* [Windows-Specific Notes](https://www.google.com/search?q=%23windows-specific-notes)

---

## Overview

MangaHub consists of 5 server components and a CLI application:

| Component | Port | Protocol | Description |
| --- | --- | --- | --- |
| **API Server** | 8080 | HTTP/REST | User authentication and library management |
| **TCP Server** | 9090 | TCP | Real-time reading progress synchronization |
| **UDP Server** | 9091 | UDP | Broadcast notifications for new releases |
| **gRPC Server** | 9092 | gRPC | High-performance internal data operations |
| **WebSocket Server** | 9093 | WebSocket | Real-time community chat functionality |

---

## Prerequisites (Windows)

### 1. Install Go

* Download the `.msi` installer from [golang.org/dl](https://golang.org/dl/).
* Verify installation in PowerShell:
```powershell
go version

```



### 2. Install Protocol Buffers (protoc)

* Download `protoc-xx.x-win64.zip` from [GitHub Releases](https://github.com/protocolbuffers/protobuf/releases).
* Extract it to a permanent folder (e.g., `C:\proto`).
* Add the `bin` folder (e.g., `C:\proto\bin`) to your **Environment Variables (Path)**.

### 3. Install Go Protobuf Plugins

Open PowerShell and run:

```powershell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

```

**Note:** Ensure `%USERPROFILE%\go\bin` is added to your Windows **Path** to execute `protoc-gen-go`.

---

## Installation

1. Navigate to the project directory:
```powershell
cd "C:\Path\To\Your\Project\mangahub"

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

To create `.exe` executables for Windows:

```powershell
# API Server
go build -o api-server.exe cmd\api-server\main.go

# TCP Server
go build -o tcp-server.exe cmd\tcp-server\main.go

# UDP Server
go build -o udp-server.exe cmd\udp-server\main.go

# CLI App
go build -o mangahub.exe cmd\cli-app\main.go

```

---

## Running the Servers

Open each server in a separate **PowerShell** or **Command Prompt** window:

* **Terminal 1 (API):** `go run cmd\api-server\main.go`
* **Terminal 2 (TCP):** `go run cmd\tcp-server\main.go`
* **Terminal 3 (UDP):** `go run cmd\udp-server\main.go`

---

## Testing on Windows

### 1. HTTP API (Using native Windows curl.exe)

**Register:**

```powershell
curl.exe -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d '{\"username\":\"testuser\",\"password\":\"password123\"}'

```

**Login (To get JWT Token):**

```powershell
curl.exe -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d '{\"username\":\"testuser\",\"password\":\"password123\"}'

```

### 2. TCP Sync (Using PowerShell Socket)

Open a new PowerShell window to act as a "Listening Device":

```powershell
$client = New-Object System.Net.Sockets.TCPClient("localhost", 9090); `
$stream = $client.GetStream(); $writer = New-Object System.IO.StreamWriter($stream); `
$writer.AutoFlush = $true; $writer.WriteLine('{"user_id":"YOUR_USER_ID_HERE"}'); `
$reader = New-Object System.IO.StreamReader($stream); `
while($client.Connected) { $line = $reader.ReadLine(); if($line) { Write-Host "🔔 SYNC RECEIVED: $line" -ForegroundColor Cyan } }

```

### 3. UDP Notifications

Send a manual notification via PowerShell:

```powershell
$udp = New-Object System.Net.Sockets.UdpClient; `
$udp.Connect("127.0.0.1", 9091); `
$data = [System.Text.Encoding]::ASCII.GetBytes("New Release: One Piece Chapter 1111 is now available!"); `
$udp.Send($data, $data.Length); $udp.Close()

```

---

## Database Management

The project uses SQLite. On Windows, you can manage the database via:

1. **DBeaver:** Download the Windows version from [dbeaver.io](https://dbeaver.io/). Point it to the `mangahub.db` file.
2. **SQLite CLI:**
* Download `sqlite-tools-win32-x86.zip` from [sqlite.org](https://www.sqlite.org/download.html).
* Run: `.\sqlite3.exe mangahub.db`



---

## Windows-Specific Notes ⚠️

* **File Paths:** Always use backslashes `\` in CMD/PowerShell for directory navigation.
* **Execution Policy:** If script execution is blocked, run PowerShell as Administrator and execute: `Set-ExecutionPolicy RemoteSigned`.
* **Windows Firewall:** When running the servers for the first time, click **"Allow Access"** for both Private and Public networks to enable TCP/UDP traffic.
* **Environment Variables:** After adding `protoc` or `go/bin` to the Path, you **must restart** your Terminal/IDE for the changes to take effect.

---

*Project by Hung Khang - Network Programming 2025*
