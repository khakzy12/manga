@echo off
echo 🚀 Starting MangaHub Distributed System...

:: Ensure Data folder exists
if not exist "data" mkdir "data"

:: Step 1: Seed the database (runs and closes)
echo 💾 Seeding Database...
go run cmd/seed/main.go

:: Step 2: Start background services
echo 🛰️  Starting gRPC Server...
start "gRPC Server" go run cmd/grpc-server/main.go

echo 🛰️  Starting TCP Sync Server...
start "TCP Server" go run cmd/tcp-server/main.go

:: Step 3: Start the Gateway (This stays in the main window)
echo 🌐 Starting API Gateway on http://localhost:8080...
go run cmd/api-server/main.go

pause