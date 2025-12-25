package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"

	"mangahub/pkg/database"
	"mangahub/proto"

	"google.golang.org/grpc"
)

type mangaServer struct {
	proto.UnimplementedMangaServiceServer
	DB *sql.DB
}

// Implement the GetManga RPC
func (s *mangaServer) GetManga(ctx context.Context, req *proto.GetMangaRequest) (*proto.MangaResponse, error) {
	fmt.Printf("gRPC Server received request for ID: %s\n", req.Id)

	var m proto.MangaResponse
	var genresRaw string // Temporary variable to hold the JSON text from DB

	// Use LOWER to make the search case-insensitive
	query := `SELECT id, title, author, genres, status, total_chapters, description 
              FROM manga 
              WHERE id LIKE ? OR title LIKE ? LIMIT 1`
	cleanId := strings.ReplaceAll(req.Id, " ", "%")
	searchParam := "%" + cleanId + "%"

	err := s.DB.QueryRow(query, searchParam, searchParam).Scan(
		&m.Id, &m.Title, &m.Author, &genresRaw, &m.Status, &m.TotalChapters, &m.Description,
	)

	if err != nil {
		fmt.Printf("❌ Database Search Error: %v\n", err)
		return nil, err
	}

	// Convert the text ["Shounen"] back into the gRPC slice []string
	// Simple way: trim the brackets and quotes, or use json.Unmarshal
	m.Genres = strings.Split(strings.Trim(genresRaw, "[]\" "), "\",\"")

	fmt.Printf("✅ Found: %s\n", m.Title)
	return &m, nil
}

func main() {
	// 1. Init Database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Start TCP Listener for gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 3. Register Server
	s := grpc.NewServer()
	proto.RegisterMangaServiceServer(s, &mangaServer{DB: db})

	log.Println("🚀 gRPC Internal Service running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
