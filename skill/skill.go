package skill

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"test/response"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

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

type Handler struct {
	Db *sql.DB
}

func GetPing(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h *Handler) GetAllSkills(context *gin.Context) {
	rows, err := h.Db.Query("SELECT key, name, description, logo, tags FROM skill;")
	if err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Can not get all skill"})
		return
	}
	defer rows.Close()
	var skills []Skill
	for rows.Next() {
		var s Skill
		err := rows.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
		if err != nil {
			context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Can not get all skill"})
			return
		}
		skills = append(skills, s)
	}
	context.JSON(http.StatusOK, response.Success{Status: "success", Data: skills})
}

func (h *Handler) GetSkillById(context *gin.Context) {
	paramkey := context.Param("key")
	getSkillByKey(h, paramkey, context)
}

func (h *Handler) CreateSkill(context *gin.Context) {
	var newSkill Skill
	err := context.BindJSON(&newSkill)
	if err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.Db.Prepare("INSERT INTO skill (key, name, description, logo, tags) VALUES ($1, $2, $3, $4, $5) returning key;")
	if err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(newSkill.Key, newSkill.Name, newSkill.Description, newSkill.Logo, pq.Array(newSkill.Tags)); err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Skill already exists"})
		return
	}
	context.JSON(http.StatusOK, response.Success{Status: "success", Data: newSkill})
}

func (h *Handler) UpdateSkill(context *gin.Context) {
	var paramkey = context.Param("key")
	var s PostSkill
	if err := context.BindJSON(&s); err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.Db.Prepare("UPDATE skill SET name = $1, description = $2, logo = $3, tags = $4 WHERE key = $5;")
	if err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(s.Name, s.Description, s.Logo, pq.Array(s.Tags), paramkey); err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Skill already exists"})
		return
	}
	getSkillByKey(h, paramkey, context)
}

func getSkillByKey(h *Handler, key string, context *gin.Context) {
	skill := h.Db.QueryRow(fmt.Sprintf("SELECT key, name, description, logo, tags FROM skill WHERE key = '%v';", key))
	var s Skill
	err := skill.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
	if err != nil {
		context.JSON(http.StatusNotFound, response.Fail{Status: "error", Message: "Skill not found"})
		return
	}
	context.JSON(http.StatusOK, response.Success{Status: "success", Data: s})
}

func (h *Handler) UpdateSkillName(context *gin.Context) {
	var paramkey = context.Param("key")
	var name struct {
		Name string `json:"name" binding:"required"`
	}
	if err := context.BindJSON(&name); err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.Db.Prepare("UPDATE skill SET name = $1 WHERE key = $2;")
	if err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(name.Name, paramkey); err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Not be able to update name"})
		return
	}
	getSkillByKey(h, paramkey, context)
}

func (h *Handler) UpdateSkillDescription(context *gin.Context) {
	var paramkey = context.Param("key")
	var description struct {
		Description string `json:"description" binding:"required"`
	}
	if err := context.BindJSON(&description); err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.Db.Prepare("UPDATE skill SET description = $1 WHERE key = $2;")
	if err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Statement error"})
		return
	}
	defer stmt.Close()
	fmt.Print("what")
	if _, err := stmt.Exec(description.Description, paramkey); err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Not be able to update description"})
		return
	}
	getSkillByKey(h, paramkey, context)
}

func (h *Handler) UpdateSkillLogo(context *gin.Context) {
	var paramkey = context.Param("key")
	var logo struct {
		Logo string `json:"logo" binding:"required"`
	}
	if err := context.BindJSON(&logo); err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.Db.Prepare("UPDATE skill SET logo = $1 WHERE key = $2;")
	if err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(logo.Logo, paramkey); err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Not be able to update logo"})
		return
	}
	getSkillByKey(h, paramkey, context)
}

func (h *Handler) UpdateSkillTags(context *gin.Context) {
	var paramkey = context.Param("key")
	var tags struct {
		Tags []string `json:"tags" binding:"required"`
	}
	if err := context.BindJSON(&tags); err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.Db.Prepare("UPDATE skill SET tags = $1 WHERE key = $2;")
	if err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(pq.Array(tags.Tags), paramkey); err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Not be able to update tags"})
		return
	}
	getSkillByKey(h, paramkey, context)
}

func (h *Handler) DeleteSkill(context *gin.Context) {
	paramkey := context.Param("key")
	var mutex = &sync.Mutex{}
	mutex.Lock()
	skill := h.Db.QueryRow(fmt.Sprintf("SELECT key, name, description, logo, tags FROM skill WHERE key = '%v';", paramkey))
	var s Skill
	err := skill.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
	if err != nil {
		context.JSON(http.StatusNotFound, response.Fail{Status: "error", Message: "Skill not found"})
		return
	}
	stmt, err := h.Db.Prepare("DELETE FROM skill WHERE key = $1;")
	if err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(paramkey); err != nil {
		context.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Not be able to delete skill"})
		return
	}
	mutex.Unlock()
	context.JSON(http.StatusOK, response.Success{Status: "success", Data: "Skill deleted"})
}
