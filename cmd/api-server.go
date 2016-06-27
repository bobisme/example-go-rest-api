package cmd

import "github.com/tucnak/climax"

func startAPIServer(ctx climax.Context) int {
	// cfg := conf.LoadFile(getConfigPath(ctx))
	// if cfg.ReleaseMode {
	// 	gin.SetMode(gin.ReleaseMode)
	// }
	//
	// // load database
	// db, err := gorm.Open("sqlite", cfg.DBPath)
	// if err != nil {
	// 	logrus.Panicln(err)
	// }
	//
	// r := api.GetRouter(cfg, db)
	// err = r.Run(":" + strconv.Itoa(cfg.Port))
	// if err != nil {
	// 	logrus.Panicln(err)
	// }
	//
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
