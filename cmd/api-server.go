package cmd

import (
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/bobisme/RestApiProject/api"
	"github.com/bobisme/RestApiProject/conf"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // load sqlite3 support
	"github.com/tucnak/climax"
)

func startAPIServer(ctx climax.Context) int {
	cfg := conf.LoadFile(getConfigPath(ctx))
	if cfg.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// load database
	db, err := gorm.Open("sqlite3", cfg.DBPath)
	if err != nil {
		logrus.Errorln(err)
		return 1
	}

	r := api.GetRouter(cfg, db)
	err = r.Run(":" + strconv.Itoa(cfg.Port))
	if err != nil {
		logrus.Errorln(err)
		return 1
	}

	return 0
}

// APIServer command starts the REST API server
var APIServer = climax.Command{
	Name:   "api-server",
	Brief:  "start the REST API server",
	Help:   "Start the REST API server.\n",
	Flags:  []climax.Flag{debugFlag, configFlag},
	Handle: startAPIServer,
}
