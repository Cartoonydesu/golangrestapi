package skill

import (
	"cartoon/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *Handler) DeleteSkill(c *gin.Context) {
	paramkey := c.Param("key")
	skill := h.Db.QueryRow(fmt.Sprintf("SELECT key, name, description, logo, tags FROM skill WHERE key = '%v';", paramkey))
	var s Skill
	err := skill.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, pq.Array(&s.Tags))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail{Status: "error", Message: "Skill not found"})
		return
	}
	stmt, err := h.Db.Prepare("DELETE FROM skill WHERE key = $1;")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(paramkey); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Not be able to delete skill"})
		return
	}
	c.JSON(http.StatusOK, response.Success{Status: "success", Data: "Skill deleted"})
}
