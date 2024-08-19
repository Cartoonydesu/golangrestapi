package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test/skill"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// slog.Info(os.Getenv("POSTGRES_URI"))
	connStr := os.Getenv("POSTGRES_URI")
	// connStr := "postgres://postgres:1234@127.0.0.1:5432/app?sslmode=disable"
	var err error
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	h := &skill.Handler{Db: db}
	router := gin.Default()
	skill.SetRouter(router, h)
	srv := http.Server{
		Addr:        ":" + os.Getenv("PORT"),
		Handler:     router,
		ReadTimeout: 3 * time.Second,
	}
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		slog.Info("shutting down")
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}
	slog.Info("bye")
}
