package skill

import (
	"cartoon/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *Handler) UpdateSkillName(c *gin.Context) {
	var key = c.Param("key")
	var name struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.BindJSON(&name); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.Db.Prepare("UPDATE skill SET name = $1 WHERE key = $2;")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(name.Name, key); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Not be able to update name"})
		return
	}
	sk, err := getSkillByKey(h.Db, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail{Status: "error", Message: "Skill not found" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Success{Status: "success", Data: sk})
}

func (h *Handler) UpdateSkillDescription(c *gin.Context) {
	var key = c.Param("key")
	var description struct {
		Description string `json:"description" binding:"required"`
	}
	if err := c.BindJSON(&description); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.Db.Prepare("UPDATE skill SET description = $1 WHERE key = $2;")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(description.Description, key); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Not be able to update description"})
		return
	}
	sk, err := getSkillByKey(h.Db, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail{Status: "error", Message: "Skill not found" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Success{Status: "success", Data: sk})
}

func (h *Handler) UpdateSkillLogo(c *gin.Context) {
	var key = c.Param("key")
	var logo struct {
		Logo string `json:"logo" binding:"required"`
	}
	if err := c.BindJSON(&logo); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.Db.Prepare("UPDATE skill SET logo = $1 WHERE key = $2;")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(logo.Logo, key); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Not be able to update logo"})
		return
	}
	sk, err := getSkillByKey(h.Db, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail{Status: "error", Message: "Skill not found" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Success{Status: "success", Data: sk})
}

func (h *Handler) UpdateSkillTags(c *gin.Context) {
	var key = c.Param("key")
	var tags struct {
		Tags []string `json:"tags" binding:"required"`
	}
	if err := c.BindJSON(&tags); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Cannot extract data from JSON"})
		return
	}
	stmt, err := h.Db.Prepare("UPDATE skill SET tags = $1 WHERE key = $2;")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Statement error"})
		return
	}
	defer stmt.Close()
	if _, err := stmt.Exec(pq.Array(tags.Tags), key); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail{Status: "error", Message: "Not be able to update tags"})
		return
	}
	sk, err := getSkillByKey(h.Db, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail{Status: "error", Message: "Skill not found" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.Success{Status: "success", Data: sk})
}
