package skill

import (
	"github.com/gin-gonic/gin"
)

// type handler struct {
// 	db *sql.DB
// }

func SetRouter(router *gin.Engine, h *Handler) {
	router.GET("/ping", GetPing)
	router.GET("/api/v1/skills", h.GetAllSkills)
	router.GET("/api/v1/skills/:key", h.GetSkillById)
	router.POST("/api/v1/skills", h.CreateSkill)
	router.PUT("/api/v1/skills/:key", h.UpdateSkill)
	router.PATCH("/api/v1/skills/:key/action/name", h.UpdateSkillName)
	router.PATCH("/api/v1/skills/:key/action/description", h.UpdateSkillDescription)
	router.PATCH("/api/v1/skills/:key/action/logo", h.UpdateSkillLogo)
	router.PATCH("/api/v1/skills/:key/action/tags", h.UpdateSkillTags)
	router.DELETE("/api/v1/skills/:key", h.DeleteSkill)
}
