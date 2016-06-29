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

// parse and verify that the user exists
// sends a json error response and returns 0 if invalid
func getUser(c *gin.Context, db *gorm.DB) *models.User {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		jsonError(c, "could not parse user id", err)
		return nil
	}

	user := &models.User{}
	q := db.Where("id = ?", userID).First(user)
	if err := q.Error; err != nil {
		jsonError(c, "error looking up user", err)
		return nil
	} else if q.RecordNotFound() {
		jsonError(c, "user not found", nil)
		return nil
	}
	return user
}

func getNewVisitHandler(cfg *conf.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VisitRequest
		err := c.BindJSON(&req)
		if err != nil {
			jsonError(c, "could not understand your data", err)
			return
		}
		var state models.State
		var city models.City
		user := getUser(c, db)
		if user == nil {
			return
		}

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
			v := models.Visit{User: *user, City: city}
			db.Create(&v)
			c.JSON(http.StatusCreated, &v)
			return
		}
	}
}

func getDeleteVisitHandler(cfg *conf.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c, db)
		if user == nil {
			return
		}
		visitID, err := strconv.Atoi(c.Param("visitID"))
		if err != nil {
			jsonError(c, "could not parse visit id", err)
			return
		}
		var visit models.Visit
		q := db.Where("id = ?", visitID).First(&visit)
		if err := q.Error; err != nil {
			jsonError(c, "error looking up visit", err)
			return
		} else if q.RecordNotFound() {
			jsonError(c, "visit not found", nil)
			return
		}
		if err := db.Delete(&visit).Error; err != nil {
			jsonError(c, "error removing visit", err)
			return
		}
		c.Status(http.StatusNoContent)
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
	r.DELETE("/user/:userID/visits/:visitID", getDeleteVisitHandler(cfg, db))

	return r
}
