package api

import (
	"net/http"
	"strconv"

	"github.com/bobisme/RestApiProject/conf"
	"github.com/bobisme/RestApiProject/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const defaultLimit = 100
const maxLimit = 1000

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

func getLimitOffset(c *gin.Context) (uint, uint) {
	limit := defaultLimit
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if limit < 1 {
		limit = defaultLimit
	}

	offset := 0
	offStr := c.Query("offset")
	offset, err = strconv.Atoi(offStr)
	if err != nil {
		offset = 0
	}
	if offset < 0 {
		offset = 0
	}

	return uint(limit), uint(offset)
}

func getStateCitiesHandler(cfg *conf.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stateID, err := strconv.Atoi(c.Param("stateID"))
		if err != nil {
			jsonError(c, "Could not get id for state", err)
			return
		}
		var cities []models.City
		limit, offset := getLimitOffset(c)
		var count int
		q := db.Model(&models.City{}).Where("state_id = ?", stateID)
		q.Count(&count)
		q.Limit(limit).Offset(offset).Find(&cities)
		c.JSON(http.StatusOK, &MetaResponse{
			limit, offset, uint(count), cities,
		})
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
			v := models.Visit{UserID: user.ID, City: city}
			if err := db.Create(&v).Error; err != nil {
				jsonError(c, "error saving visit", err)
				return
			}
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

// MetaResponse wraps the response and provides more info
type MetaResponse struct {
	Limit  uint        `json:"limit"`
	Offset uint        `json:"offset"`
	Count  uint        `json:"count"`
	Data   interface{} `json:"data"`
}

func getVisitedCitiesHandler(cfg *conf.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c, db)
		if user == nil {
			return
		}
		var cities []models.City
		limit, offset := getLimitOffset(c)
		queryBase := `
			FROM cities
			WHERE cities.id IN (
				SELECT DISTINCT city_id
				FROM visits
				WHERE user_id = ?
			)
		`
		var count int
		q := db.Raw(`SELECT COUNT(*) `+queryBase, user.ID).Count(&count)
		if err := q.Error; err != nil {
			jsonError(c, "error counting cities", err)
			return
		}
		// not the most efficient method, but easiest to implement
		q = db.Raw(
			`SELECT cities.* `+queryBase+` LIMIT ? OFFSET ?`,
			user.ID, limit, offset).Scan(&cities)
		if err := q.Error; err != nil {
			jsonError(c, "error looking up cities", err)
			return
		}
		c.JSON(http.StatusOK, &MetaResponse{
			limit, offset, uint(count), cities,
		})
	}
}

func getVisitedStatesHandler(cfg *conf.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c, db)
		if user == nil {
			return
		}
		var states []models.State
		q := db.Raw(`
			SELECT *
			FROM states
			WHERE id IN (
				SELECT cities.state_id
				FROM cities
				LEFT JOIN visits ON cities.id = visits.city_id
				WHERE visits.user_id = ?
				GROUP BY cities.state_id
			)
		`, user.ID).Scan(&states)
		if err := q.Error; err != nil {
			jsonError(c, "error looking up states", err)
			return
		}
		c.JSON(http.StatusOK, &states)
	}
}

// GetRouter for the API Server router
func GetRouter(cfg *conf.Config, db *gorm.DB) *gin.Engine {
	// create a default router with logger and recovery
	r := gin.Default()
	SetRoutes(cfg, db, r)
	return r
}

// SetRoutes for the API Server router
func SetRoutes(cfg *conf.Config, db *gorm.DB, r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "HELLO")
	})
	r.GET("/state/:stateID/cities", getStateCitiesHandler(cfg, db))
	r.POST("/user/:userID/visits", getNewVisitHandler(cfg, db))
	r.DELETE("/user/:userID/visits/:visitID", getDeleteVisitHandler(cfg, db))
	r.GET("/user/:userID/visits/states", getVisitedStatesHandler(cfg, db))
	r.GET("/user/:userID/visits", getVisitedCitiesHandler(cfg, db))
}
