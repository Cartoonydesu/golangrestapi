package skill

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSkillById(c *gin.Context) {
	paramkey := c.Param("key")
	s, err := getSkillByKey(h.Db, paramkey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, s)
}
