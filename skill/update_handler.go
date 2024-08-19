package skill

import (
	"cartoon/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *Handler) UpdateSkill(c *gin.Context) {
	var paramkey = c.Param("key")
	var s CreateSkill
	if err := c.BindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.Db.Prepare("UPDATE skill SET name = $1, description = $2, logo = $3, tags = $4 WHERE key = $5;")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(s.Name, s.Description, s.Logo, pq.Array(s.Tags), paramkey); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Skill already exists"})
		return
	}
	sk, err := getSkillByKey(h.Db, paramkey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, sk)
}
