package main

import (
	"context"
	"log"
	"mangahub/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewMangaServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetManga(ctx, &proto.GetMangaRequest{Id: "1"})
	if err != nil {
		log.Fatalf("could not get manga: %v", err)
	}
	log.Printf("Manga Found: %s by %s", r.GetTitle(), r.GetAuthor())
}
