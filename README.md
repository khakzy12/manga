MangaHub - Network Programming Project
A comprehensive manga management system demonstrating multiple network protocols: HTTP REST API, TCP, UDP, WebSocket, and gRPC.

📋 Table of Contents
Overview
Prerequisites
Installation
Project Structure
Building the Project
Running the Servers
Testing
Database Management
Troubleshooting
Windows Notes
Overview
MangaHub consists of 5 server components and a CLI application:

Component	Port	Protocol	Description
API Server	8080	HTTP/REST	User authentication and manga management
TCP Server	9090	TCP	Real-time progress synchronization
UDP Server	9091	UDP	Broadcast notifications for new releases
gRPC Server	9092	gRPC	Manga operations via Protocol Buffers
WebSocket Server	9093	WebSocket	Real-time chat functionality
CLI App	-	-	Unified command-line interface
Prerequisites
macOS Requirements
Go (version 1.21 or later)

brew install go
# Verify installation
go version
Protocol Buffers Compiler (protoc)

brew install protobuf
# Verify installation
protoc --version
Go Protobuf Plugins

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
Database Tools (Optional but recommended)

DBeaver: Download from dbeaver.io
SQLite: Usually pre-installed on macOS, or install via Homebrew:
brew install sqlite
Network Tools (for testing)

netcat (nc) - Usually pre-installed on macOS
curl - Pre-installed on macOS
Setting Up Protobuf PATH
The Go protobuf plugins are installed to ~/go/bin. Add this to your PATH:

For zsh (default on macOS):

# Add to ~/.zshrc
echo 'export PATH=$PATH:~/go/bin' >> ~/.zshrc
source ~/.zshrc
For bash:

# Add to ~/.bash_profile
echo 'export PATH=$PATH:~/go/bin' >> ~/.bash_profile
source ~/.bash_profile
Verify PATH setup:

which protoc-gen-go
which protoc-gen-go-grpc
Both commands should return paths like /Users/yourusername/go/bin/protoc-gen-go

Installation
Clone or navigate to the project directory:

cd /Users/huynhngocanhthu/mangahub
Install Go dependencies:

go mod tidy
Generate Protocol Buffer code:

cd pkg/proto
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       manga.proto
cd ../..
This generates:

pkg/proto/manga.pb.go - Message types
pkg/proto/manga_grpc.pb.go - Service stubs
Project Structure
mangahub/
├── cmd/
│   ├── api-server/         # HTTP REST API server
│   ├── tcp-server/          # TCP progress sync server
│   ├── udp-server/          # UDP notification server
│   ├── grpc-server/         # gRPC server
│   ├── websocket-server/    # WebSocket chat server
│   └── cli-app/             # CLI application
├── internal/
│   ├── models/              # Data models (User, Manga, UserProgress)
│   ├── database/            # SQLite database initialization
│   ├── auth/                # JWT authentication
│   └── grpc/                # gRPC server implementation
├── pkg/
│   └── proto/               # Protocol Buffer definitions
├── go.mod                   # Go module dependencies
└── mangahub.db              # SQLite database (created on first run)
Building the Project
Build CLI Application
go build -o mangahub cmd/cli-app/main.go
This creates an executable mangahub in the current directory.

Build Individual Servers
# API Server
go build -o api-server cmd/api-server/main.go

# TCP Server
go build -o tcp-server cmd/tcp-server/main.go

# UDP Server
go build -o udp-server cmd/udp-server/main.go

# gRPC Server
go build -o grpc-server cmd/grpc-server/main.go

# WebSocket Server
go build -o websocket-server cmd/websocket-server/main.go
Running the Servers
Important: You need to run all servers in separate terminal windows/tabs for the full system to work.

Terminal 1: API Server (HTTP)
go run cmd/api-server/main.go
Expected output: Server starting on port 8080

Terminal 2: TCP Server
go run cmd/tcp-server/main.go
Expected output: TCP Progress Sync Server listening on port 9090

Terminal 3: UDP Server
go run cmd/udp-server/main.go
Expected output: UDP Notification Server listening on port 9091

Terminal 4: gRPC Server
go run cmd/grpc-server/main.go
Expected output: gRPC MangaService server listening on port 9092

Terminal 5: WebSocket Server
go run cmd/websocket-server/main.go
Expected output: WebSocket Chat Server listening on port 9093

Terminal 6: CLI Testing
Keep this terminal for running CLI commands and testing.

Testing
1. Test HTTP API Server
Register a new user:

curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}'
Login:

curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass123"}'
Save the returned token for later use.

Get manga list:

curl http://localhost:8080/manga
Search manga:

curl "http://localhost:8080/manga?search=One"
2. Test CLI Application
Login via CLI:

./mangahub auth login
# Enter username and password when prompted
Search manga:

./mangahub manga search
# Or with query
./mangahub manga search "One Piece"
Join chat:

./mangahub chat join
# Type messages and press Enter. Type 'exit' to quit.
Start sync listener:

./mangahub sync start
3. Test TCP Server
Connect and authenticate:

nc localhost 9090
Then send (replace YOUR_TOKEN with actual JWT token):

{"type":"auth","token":"YOUR_TOKEN"}
Send progress update:

{"type":"ProgressUpdate","manga_id":"manga_1","chapter":100}
4. Test UDP Server
Subscribe to notifications:

echo "SUBSCRIBE" | nc -u localhost 9091
Send release notification (in another terminal):

echo "RELEASE:One Piece" | nc -u localhost 9091
All subscribers will receive: New Chapter Released: One Piece

5. Test WebSocket Server
Using CLI:

./mangahub chat join
Using Browser Console:

const ws = new WebSocket('ws://localhost:9093/ws');
ws.onmessage = (event) => console.log('Received:', event.data);
ws.send('Hello from browser!');
6. Test gRPC Server
Using grpcurl (install via brew install grpcurl):

# List services
grpcurl -plaintext localhost:9092 list

# Get manga by ID
grpcurl -plaintext -d '{"id":"manga_1"}' \
  localhost:9092 proto.MangaService/GetManga

# Update progress
grpcurl -plaintext -d '{"user_id":"user_123","manga_id":"manga_1","chapter":50}' \
  localhost:9092 proto.MangaService/UpdateProgress
Database Management
Using DBeaver
Open DBeaver and create a new connection:

Database type: SQLite
Path: /Users/huynhngocanhthu/mangahub/mangahub.db
View tables:

users - User accounts
manga - Manga catalog (seeded with 5 sample manga)
user_progress - User reading progress
Query examples:

-- View all manga
SELECT * FROM manga;

-- View user progress
SELECT * FROM user_progress;

-- View users
SELECT id, username, created_at FROM users;
Using SQLite CLI
# Open database
sqlite3 mangahub.db

# View tables
.tables

# Query manga
SELECT * FROM manga;

# Exit
.quit
Database Location
The database file mangahub.db is created in the project root directory on first server startup. You can change the location by setting the DB_PATH environment variable:

export DB_PATH=/path/to/your/database.db
go run cmd/api-server/main.go
Environment Variables
Variable	Default	Description
PORT	8080 (API), 9093 (WebSocket)	Server port
DB_PATH	mangahub.db	SQLite database path
JWT_SECRET	mangahub-secret-key-change-in-production	JWT signing secret
Troubleshooting
Port Already in Use
# Find process using the port
lsof -i :8080  # or :9090, :9091, :9092, :9093

# Kill the process
kill -9 <PID>
Database Locked
Make sure only one process is accessing mangahub.db at a time. Close DBeaver connections if you're running servers.

Protobuf Generation Fails
# Verify protoc is installed
protoc --version

# Verify plugins are in PATH
which protoc-gen-go
which protoc-gen-go-grpc

# If not found, add to PATH
export PATH=$PATH:~/go/bin
JWT Token Expired
Tokens expire after 24 hours. Re-login:

./mangahub auth login
Connection Refused
Ensure all servers are running before testing. Check server logs for errors.

Go Module Issues
# Clean and re-download dependencies
go clean -modcache
go mod tidy
Windows Notes
⚠️ This project is primarily tested on macOS. Running on Windows may require additional configuration:

Install Go for Windows:

Download from golang.org/dl
Install and add to PATH
Install Protocol Buffers:

Download from github.com/protocolbuffers/protobuf/releases
Extract and add bin directory to PATH
Install Go Protobuf Plugins:

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
Add %USERPROFILE%\go\bin to Windows PATH.

Network Tools:

netcat may not be available. Use alternatives like:
PowerShell Test-NetConnection
Third-party tools like ncat (Nmap)
Or use the CLI application instead
Database Path:

Use forward slashes or double backslashes in paths
Example: DB_PATH=C:\\Users\\YourName\\mangahub.db
Terminal:

Use PowerShell or Git Bash instead of Command Prompt for better compatibility
Quick Start Summary
# 1. Install dependencies
go mod tidy

# 2. Generate protobuf code
cd pkg/proto && protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative manga.proto && cd ../..

# 3. Build CLI
go build -o mangahub cmd/cli-app/main.go

# 4. Start all servers (in separate terminals)
go run cmd/api-server/main.go      # Terminal 1
go run cmd/tcp-server/main.go      # Terminal 2
go run cmd/udp-server/main.go      # Terminal 3
go run cmd/grpc-server/main.go    # Terminal 4
go run cmd/websocket-server/main.go # Terminal 5

# 5. Test with CLI
./mangahub auth login
./mangahub manga search
