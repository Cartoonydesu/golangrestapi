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
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Logo        string   `json:"logo"`
	Tags        []string `json:"tags"`
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
		context.JSON(http.StatusBadRequest, gin.H{"message": err})
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
	context.JSON(http.StatusOK, Successres{"success", skills})
}

func (h *handler) getSkillById(context *gin.Context) {
	paramkey := context.Param("key")
	// slog.Info(paramkey)
	row := h.db.QueryRow(fmt.Sprintf("SELECT key, name, description, logo, tags FROM skill WHERE key = '%v';", paramkey))
	var s Skill
	err := row.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
	if err != nil {
		// TODO error
		context.JSON(http.StatusBadRequest, err)
		return
	}
	context.JSON(http.StatusOK, Successres{"success", s})
}

func (h *handler) getSkillByKey(key string) Skill {
	row := h.db.QueryRow(fmt.Sprintf("SELECT key, name, description, logo, tags FROM skill WHERE key = '%v';", key))
	var s Skill
	err := row.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
	if err != nil {
		// TODO error
		// context.JSON(http.StatusBadRequest, err)
		return Skill{}
	}
	// context.JSON(http.StatusOK, Successres{"success", s})
	// context.JSON(http.StatusOK, )
	return s
}

func (h *handler) createSkill(context *gin.Context) {
	var newSkill Skill
	err := context.BindJSON(&newSkill)
	if err != nil {
		//TODO error
		context.JSON(http.StatusBadRequest, err)
		return
	}
	stmt, err := h.db.Prepare("INSERT INTO skill (key, name, description, logo, tags) VALUES ($1, $2, $3, $4, $5) returning key;")
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		//TODO error
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(newSkill.Key, newSkill.Name, newSkill.Description, newSkill.Logo, pq.Array(newSkill.Tags)); err != nil {
		context.JSON(http.StatusBadRequest, err)
		//TODO error
		return
	}
	context.JSON(http.StatusOK, Successres{"success", newSkill})
}

func (h *handler) updateSkill(context *gin.Context) {
	var paramkey = context.Param("key")
	var s PostSkill
	if err := context.BindJSON(&s); err != nil {
		context.JSON(http.StatusBadRequest, err)
		//TODO error
		return
	}
	stmt, err := h.db.Prepare("UPDATE skill SET name = $2, description = $3, logo = $4, tags = $5 WHERE key = $1;")
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		//TODO error
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(paramkey, s.Name, s.Description, s.Logo, pq.Array(s.Tags)); err != nil {
		context.JSON(http.StatusBadRequest, err)
		//TODO error
		return
	}

	row := h.db.QueryRow(fmt.Sprintf("SELECT key, name, description, logo, tags FROM skill WHERE key = '%v';", paramkey))
	var updatedSkill Skill
	err = row.Scan(&updatedSkill.Key, &updatedSkill.Name, &updatedSkill.Description, &updatedSkill.Logo, pq.Array(&updatedSkill.Tags))
	if err != nil {
		// TODO error
		context.JSON(http.StatusBadRequest, err)
		return
	}
	context.JSON(http.StatusOK, Successres{"success", updatedSkill})
}

// type name struct {
// 	Name string `json:"name"`
// }

func (h *handler) updateSkillName(context *gin.Context) {
	var paramkey = context.Param("key")
	var name struct{ Name string }
	if err := context.BindJSON(&name); err != nil {
		context.JSON(http.StatusBadRequest, err)
		//TODO error
		return
	}
	stmt, err := h.db.Prepare("UPDATE skill SET name = $2 WHERE key = $1 RETURNING key, name, description, logo, tags;")
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		//TODO error
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(paramkey, name.Name); err != nil {
		context.JSON(http.StatusBadRequest, err)
		//TODO error
		return
	}

	updatedskill := h.db.QueryRow(fmt.Sprintf("SELECT key, name, description, logo, tags FROM skill WHERE key = '%v';", paramkey))
	var updatedSkill Skill
	err = updatedskill.Scan(&updatedSkill.Key, &updatedSkill.Name, &updatedSkill.Description, &updatedSkill.Logo, pq.Array(&updatedSkill.Tags))
	if err != nil {
		// TODO error
		context.JSON(http.StatusBadRequest, err)
		return
	}
	context.JSON(http.StatusOK, Successres{"success", updatedSkill})
	// context.JSON(http.StatusOK, row)
	// row := h.db.QueryRow(fmt.Sprintf("UPDATE skill SET name = %v WHERE key = %v RETURNING key, name, description, logo, tags;", name.Name, paramkey))
	// err = row.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
	// log.Info

}

func (h *handler) updateSkillDescription(context *gin.Context) {

}

func (h *handler) updateSkillLogo(context *gin.Context) {

}

func (h *handler) updateSkillTags(context *gin.Context) {

}

func (h *handler) deleteSkill(context *gin.Context) {
	paramkey := context.Param("key")
	stmt, err := h.db.Prepare("DELETE FROM skill WHERE key = $1;")
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		//TODO error
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(paramkey); err != nil {
		context.JSON(http.StatusBadRequest, Failres{"success", "Skill deleted"})
		return
	}
	context.JSON(http.StatusOK, Successres{"success", "Skill deleted"})
}
