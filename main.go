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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func main() {
	// fmt.Println("Hello world.")
	// logs := logger.New()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	connStr := "postgres://postgres:1234@127.0.0.1:5432/postgres?sslmode=disable"
	var err error
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	h := &handler{db: db}
	router := gin.Default()
	router.GET("/ping", getPing)
	router.GET("/api/v1/skill", h.GetAllSkills)

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

type handler struct {
	db *sql.DB
}

func getPing(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}

//	func newRouter() {
//		router := gin.Default()
//	}
type Skill struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Logo        string   `json:"logo"`
	Tags        []string `json:"tags"`
}

func (h *handler) GetAllSkills(context *gin.Context) {
	rows, err := h.db.Query("SELECT key, name, description, logo, tags FROM skill;")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "get all skill"})
		return
	}
	defer rows.Close()
	var skills []Skill
	for rows.Next() {
		var s Skill
		err := rows.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		skills = append(skills, s)
	}
	context.JSON(http.StatusOK, skills)
}
