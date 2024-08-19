package skill

import (
	"cartoon/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *Handler) CreateSkill(c *gin.Context) {
	var newSkill Skill
	err := c.BindJSON(&newSkill)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.Db.Prepare("INSERT INTO skill (key, name, description, logo, tags) VALUES ($1, $2, $3, $4, $5) returning key;")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(newSkill.Key, newSkill.Name, newSkill.Description, newSkill.Logo, pq.Array(newSkill.Tags)); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Skill already exists"})
		return
	}
	c.JSON(http.StatusOK, response.Success{Status: "success", Data: newSkill})
}
