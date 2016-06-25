package cmd

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/bobisme/RestApiProject/conf"
	_ "github.com/mattn/go-sqlite3" // load sqlite3 suppor
	"github.com/tucnak/climax"
)

var dataDir string
var schemaFilename string

func init() {
	_, filename, _, _ := runtime.Caller(1)
	var err error
	dataDir, err = filepath.Abs(path.Join(path.Dir(filename), "../data"))
	if err != nil {
		logrus.Fatalln(err)
	}
	schemaFilename = path.Join(dataDir, "schema.sql")
}

func createDb(filename string, force bool) error {
	// if the database exists, and is not being forced, just quit
	_, err := os.Stat(filename)
	if err == nil && !force {
		return fmt.Errorf(
			"The database already exists. Try --force-recreate-database.")
	}

	f, err := os.Create(filename)
	f.Close()
	if err != nil {
		return err
	}

	// load schema
	schemaBytes, err := ioutil.ReadFile(schemaFilename)
	if err != nil {
		return fmt.Errorf("Could not read schema: %s", err.Error())
	}
	schema := string(schemaBytes)

	// connect to database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return fmt.Errorf("Could not open database: %s", err.Error())
	}
	defer db.Close()

	// execure schema
	logrus.Debugln("RUN SCHEMA: \n", schema)
	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("Could not execute schema: %s", err.Error())
	}
	return nil
}

func initDB(ctx climax.Context) int {
	handleDebugging(ctx)
	cfg := conf.LoadFile(getConfigPath(ctx))

	err := createDb(cfg.DBPath, ctx.Is("force-recreate-database"))
	if err != nil {
		logrus.Errorln(err)
		return 1
	}

	return 0
}

// InitDB command initialized the database and loads starting data
var InitDB = climax.Command{
	Name:  "init-db",
	Brief: "initialize the database",
	Help:  "Create the database and load the starting data.\n",
	Flags: []climax.Flag{
		debugFlag, configFlag,
		climax.Flag{
			Name:  "force-recreate-database",
			Usage: `--force-recreate-database`,
			Help:  "erase and recreate database if it exists",
		},
	},
	Handle: initDB,
}
