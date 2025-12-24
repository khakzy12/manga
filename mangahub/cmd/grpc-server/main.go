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
	var genresRaw sql.NullString // Temporary variable to hold the JSON text from DB
	var chapterRaw sql.NullInt64

	query := `SELECT id, title, author, genres, status, total_chapters, chapter, description 
              FROM manga 
              WHERE id LIKE ? OR title LIKE ? LIMIT 1`
	cleanId := strings.ReplaceAll(req.Id, " ", "%")
	searchParam := "%" + cleanId + "%"

	var statusRaw sql.NullString
	var totalChRaw sql.NullInt64
	var descRaw sql.NullString
	err := s.DB.QueryRow(query, searchParam, searchParam).Scan(
		&m.Id, &m.Title, &m.Author, &genresRaw, &statusRaw, &totalChRaw, &chapterRaw, &descRaw,
	)
	if err != nil {
		fmt.Printf("❌ Database Search Error: %v\n", err)
		return nil, err
	}

	// Handle nullable fields safely
	if genresRaw.Valid {
		m.Genres = strings.Split(strings.Trim(genresRaw.String, "[]\" "), "\",\"")
	} else {
		m.Genres = []string{}
	}
	if chapterRaw.Valid {
		m.Chapter = int32(chapterRaw.Int64)
	} else {
		m.Chapter = 0
	}
	if statusRaw.Valid {
		m.Status = statusRaw.String
	} else {
		m.Status = ""
	}
	if totalChRaw.Valid {
		m.TotalChapters = int32(totalChRaw.Int64)
	} else {
		m.TotalChapters = 0
	}
	if descRaw.Valid {
		m.Description = descRaw.String
	} else {
		m.Description = ""
	}

	fmt.Printf("✅ Found: %s\n", m.Title)
	return &m, nil
}

// Implement SearchManga (returns multiple results)
func (s *mangaServer) SearchManga(ctx context.Context, req *proto.SearchRequest) (*proto.SearchResponse, error) {
	fmt.Printf("gRPC Search received query: %s\n", req.Query)
	query := `SELECT id, title, author, genres, status, total_chapters, chapter, description 
	          FROM manga 
	          WHERE id LIKE ? OR title LIKE ? LIMIT 50`
	searchParam := "%" + strings.ReplaceAll(req.Query, " ", "%") + "%"

	rows, err := s.DB.Query(query, searchParam, searchParam)
	if err != nil {
		fmt.Printf("❌ Search Query Error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var results []*proto.MangaResponse
	for rows.Next() {
		var m proto.MangaResponse
		var genresRaw sql.NullString
		var chapterRaw sql.NullInt64
		var statusRaw sql.NullString
		var totalChRaw sql.NullInt64
		var descRaw sql.NullString
		if err := rows.Scan(&m.Id, &m.Title, &m.Author, &genresRaw, &statusRaw, &totalChRaw, &chapterRaw, &descRaw); err != nil {
			fmt.Printf("⚠️ Scan row error: %v\n", err)
			continue
		}
		if genresRaw.Valid {
			m.Genres = strings.Split(strings.Trim(genresRaw.String, "[]\" "), "\",\"")
		} else {
			m.Genres = []string{}
		}
		if chapterRaw.Valid {
			m.Chapter = int32(chapterRaw.Int64)
		} else {
			m.Chapter = 0
		}
		if statusRaw.Valid {
			m.Status = statusRaw.String
		} else {
			m.Status = ""
		}
		if totalChRaw.Valid {
			m.TotalChapters = int32(totalChRaw.Int64)
		} else {
			m.TotalChapters = 0
		}
		if descRaw.Valid {
			m.Description = descRaw.String
		} else {
			m.Description = ""
		}
		results = append(results, &m)
	}

	return &proto.SearchResponse{Results: results}, nil
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
