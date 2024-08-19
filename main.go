package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	connStr := "postgres://postgres:1234@127.0.0.1:5432/app?sslmode=disable"
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

	setRouter(router, h)

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

type PostSkill struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Logo        string   `json:"logo" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
}

type Successres struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

type Failres struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func setRouter(router *gin.Engine, h *handler) {
	router.GET("/ping", getPing)
	router.GET("/api/v1/skills", h.getAllSkills)
	router.GET("/api/v1/skills/:key", h.getSkillById)
	router.POST("/api/v1/skills", h.createSkill)
	router.PUT("/api/v1/skills/:key", h.updateSkill)
	router.PATCH("/api/v1/skills/:key/action/name", h.updateSkillName)
	router.PATCH("/api/v1/skills/:key/action/description", h.updateSkillDescription)
	router.PATCH("/api/v1/skills/:key/action/logo", h.updateSkillLogo)
	router.PATCH("/api/v1/skills/:key/action/tags", h.updateSkillTags)
	router.DELETE("/api/v1/skills/:key", h.deleteSkill)
}

func (h *handler) getAllSkills(context *gin.Context) {
	rows, err := h.db.Query("SELECT key, name, description, logo, tags FROM skill;")
	if err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Can not get all skill"})
		return
	}
	defer rows.Close()
	var skills []Skill
	for rows.Next() {
		var s Skill
		err := rows.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
		if err != nil {
			context.JSON(http.StatusBadRequest, Failres{"error", "Can not get all skill"})
			return
		}
		skills = append(skills, s)
	}
	context.JSON(http.StatusOK, Successres{"success", skills})
}

func (h *handler) getSkillById(context *gin.Context) {
	paramkey := context.Param("key")
	h.getSkillByKey(paramkey, context)
}

func (h *handler) createSkill(context *gin.Context) {
	var newSkill Skill
	err := context.BindJSON(&newSkill)
	if err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.db.Prepare("INSERT INTO skill (key, name, description, logo, tags) VALUES ($1, $2, $3, $4, $5) returning key;")
	if err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(newSkill.Key, newSkill.Name, newSkill.Description, newSkill.Logo, pq.Array(newSkill.Tags)); err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Skill already exists"})
		return
	}
	context.JSON(http.StatusOK, Successres{"success", newSkill})
}

func (h *handler) updateSkill(context *gin.Context) {
	var paramkey = context.Param("key")
	var s PostSkill
	if err := context.BindJSON(&s); err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.db.Prepare("UPDATE skill SET name = $1, description = $2, logo = $3, tags = $4 WHERE key = $5;")
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(s.Name, s.Description, s.Logo, pq.Array(s.Tags), paramkey); err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Skill already exists"})
		return
	}
	h.getSkillByKey(paramkey, context)
}

func (h *handler) getSkillByKey(key string, context *gin.Context) {
	skill := h.db.QueryRow(fmt.Sprintf("SELECT key, name, description, logo, tags FROM skill WHERE key = '%v';", key))
	var s Skill
	err := skill.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
	if err != nil {
		context.JSON(http.StatusNotFound, Failres{"error", "Skill not found"})
		return
	}
	context.JSON(http.StatusOK, Successres{"success", s})
}

func (h *handler) updateSkillName(context *gin.Context) {
	var paramkey = context.Param("key")
	var name struct {
		Name string `json:"name" binding:"required"`
	}
	if err := context.BindJSON(&name); err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.db.Prepare("UPDATE skill SET name = $1 WHERE key = $2;")
	if err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(name.Name, paramkey); err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Not be able to update name"})
		return
	}
	h.getSkillByKey(paramkey, context)
}

func (h *handler) updateSkillDescription(context *gin.Context) {
	var paramkey = context.Param("key")
	var description struct {
		Description string `json:"description" binding:"required"`
	}
	if err := context.BindJSON(&description); err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.db.Prepare("UPDATE skill SET description = $1 WHERE key = $2;")
	if err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Statement error"})
		return
	}
	defer stmt.Close()
	fmt.Print("what")
	if _, err := stmt.Exec(description.Description, paramkey); err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Not be able to update description"})
		return
	}
	h.getSkillByKey(paramkey, context)
}

func (h *handler) updateSkillLogo(context *gin.Context) {
	var paramkey = context.Param("key")
	var logo struct {
		Logo string `json:"logo" binding:"required"`
	}
	if err := context.BindJSON(&logo); err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.db.Prepare("UPDATE skill SET logo = $1 WHERE key = $2;")
	if err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(logo.Logo, paramkey); err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Not be able to update logo"})
		return
	}
	h.getSkillByKey(paramkey, context)
}

func (h *handler) updateSkillTags(context *gin.Context) {
	var paramkey = context.Param("key")
	var tags struct {
		Tags []string `json:"tags" binding:"required"`
	}
	if err := context.BindJSON(&tags); err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.db.Prepare("UPDATE skill SET tags = $1 WHERE key = $2;")
	if err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(pq.Array(tags.Tags), paramkey); err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Not be able to update tags"})
		//TODO error
		return
	}
	h.getSkillByKey(paramkey, context)

}

func (h *handler) deleteSkill(context *gin.Context) {
	paramkey := context.Param("key")
	skill := h.db.QueryRow(fmt.Sprintf("SELECT key, name, description, logo, tags FROM skill WHERE key = '%v';", paramkey))
	var s Skill
	err := skill.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
	if err != nil {
		context.JSON(http.StatusNotFound, Failres{"error", "Skill not found"})
		return
	}
	stmt, err := h.db.Prepare("DELETE FROM skill WHERE key = $1;")
	if err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(paramkey); err != nil {
		context.JSON(http.StatusBadRequest, Failres{"error", "Not be able to delete skill"})
		return
	}
	context.JSON(http.StatusOK, Successres{"success", "Skill deleted"})
}
