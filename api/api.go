package api

import (
	"net/http"
	"strconv"

	"github.com/bobisme/RestApiProject/conf"
	"github.com/bobisme/RestApiProject/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// VisitRequest is the struct for posting visit data
type VisitRequest struct {
	City  string `json:"city"`
	State string `json:"state"`
}

func jsonError(c *gin.Context, message string, err error) {
	c.JSON(http.StatusBadRequest, []map[string]map[string]string{
		{"error": {"message": message, "detail": err.Error()}},
	})
}

func getStateCitiesHandler(cfg *conf.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stateID, err := strconv.Atoi(c.Param("stateID"))
		if err != nil {
			jsonError(c, "Could not get id for state", err)
			return
		}
		var cities []models.City
		db.Where("state_id = ?", stateID).Find(&cities)
		c.JSON(http.StatusOK, cities)
	}
}

func getNewVisitHandler(cfg *conf.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VisitRequest
		err := c.BindJSON(&req)
		if err != nil {
			jsonError(c, "could not understand your data", err)
			return
		}
		userID, err := strconv.Atoi(c.Param("userID"))
		if err != nil {
			jsonError(c, "could not parse user id", err)
			return
		}

		var user models.User
		q := db.Where("id = ?", userID).First(&user)
		if err := q.Error; err != nil {
			jsonError(c, "error looking up user", err)
			return
		} else if q.RecordNotFound() {
			jsonError(c, "user not found", nil)
			return
		}
		var state models.State
		var city models.City

		if req.City != "" && req.State != "" {
			q := db.Where("abbrev = ?", req.State).First(&state)
			if err := q.Error; err != nil {
				jsonError(c, "error looking up state", err)
				return
			} else if q.RecordNotFound() {
				jsonError(c, "state not found", nil)
				return
			}
			q = db.Where(
				"name = ? AND state_id = ?", req.City, state.ID).First(&city)
			if err := q.Error; err != nil {
				jsonError(c, "error looking up city", err)
				return
			} else if q.RecordNotFound() {
				jsonError(c, "city not found", nil)
				return
			}
			v := models.Visit{User: user, City: city}
			db.Create(&v)
			c.Status(http.StatusCreated)
			return
		}
	}
}

// GetRouter for the API Server router
func GetRouter(cfg *conf.Config, db *gorm.DB) *gin.Engine {
	// create a default router with logger and recovery
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "HELLO")
	})

	r.GET("/state/:stateID/cities", getStateCitiesHandler(cfg, db))
	r.POST("/user/:userID/visits", getNewVisitHandler(cfg, db))

	return r
}
