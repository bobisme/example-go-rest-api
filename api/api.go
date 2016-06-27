package api

import (
	"net/http"
	"strconv"

	"github.com/bobisme/RestApiProject/conf"
	"github.com/bobisme/RestApiProject/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func jsonError(c *gin.Context, message string, err error) {
	c.JSON(http.StatusBadRequest, []map[string]map[string]string{
		{"error": {"message": message, "detail": err.Error()}},
	})
}

// GetRouter for the API Server router
func GetRouter(cfg *conf.Config, db *gorm.DB) *gin.Engine {
	// create a default router with logger and recovery
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "HELLO")
	})

	r.GET("/state/:stateID/cities", func(c *gin.Context) {
		stateID, err := strconv.Atoi(c.Param("stateID"))
		if err != nil {
			jsonError(c, "Could not get id for state", err)
		}
		var cities []models.City
		db.Where("state_id = ?", stateID).Find(&cities)
		c.JSON(http.StatusOK, cities)
	})

	return r
}
