package skill

import (
	"cartoon/response"
	"database/sql"
	"fmt"
	"net/http"

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

type CreateSkill struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Logo        string   `json:"logo" binding:"required"`
	Tags        []string `json:"tags" binding:"required"`
}

type Handler struct {
	Db *sql.DB
}

func GetPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h *Handler) GetAllSkills(c *gin.Context) {
	rows, err := h.Db.Query("SELECT key, name, description, logo, tags FROM skill;")
	if err != nil || rows.Err() != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Can not get all skill"})
		return
	}
	defer rows.Close()
	var skills []Skill
	for rows.Next() {
		var s Skill
		err := rows.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Can not get all skill"})
			return
		}
		skills = append(skills, s)
	}
	c.JSON(http.StatusOK, response.Success{Status: "success", Data: skills})
}

// func getSkillByKey(h *Handler, key string, c *gin.Context) {
// 	skill := h.Db.QueryRow(fmt.Sprintf("SELECT key, name, description, logo, tags FROM skill WHERE key = '%v';", key))
// 	var s Skill
// 	err := skill.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
// 	if err != nil {
// 		e := response.Fail{Status: "error", Message: "Skill not found"}
// 		c.JSON(http.StatusStatusInternalServerError, e)
// 		return
// 	}
// 	sx := response.Success{Status: "success", Data: s}
// 	c.JSON(http.StatusOK, sx)
// }

func getSkillByKey(db *sql.DB, key string) (Skill, error) {
	skill := db.QueryRow(fmt.Sprintf("SELECT key, name, description, logo, tags FROM skill WHERE key = '%v';", key))
	var s Skill
	err := skill.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
	if err != nil {
		return Skill{}, err
	}
	return s, nil
}
