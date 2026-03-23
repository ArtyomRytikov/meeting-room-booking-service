package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"test-backend-1-ArtyomRytikov/internal/config"
	"test-backend-1-ArtyomRytikov/internal/handler"
	postgresrepo "test-backend-1-ArtyomRytikov/internal/repo/postgres"
	"test-backend-1-ArtyomRytikov/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	var pool *pgxpool.Pool
	var err error

	for i := 0; i < 15; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		pool, err = pgxpool.New(ctx, cfg.DatabaseURL())
		if err == nil {
			err = pool.Ping(ctx)
		}

		cancel()

		if err == nil {
			break
		}

		log.Printf("failed to connect to db (attempt %d/15): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	roomRepo := postgresrepo.NewRoomRepository(pool)
	roomService := service.NewRoomService(roomRepo)

	if err := roomService.Init(context.Background()); err != nil {
		log.Fatal(err)
	}

	r := handler.NewRouter(roomService)

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
