package cmd

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bobisme/RestApiProject/conf"
	"github.com/bobisme/RestApiProject/geo"
	_ "github.com/mattn/go-sqlite3" // load sqlite3 support
	"github.com/tucnak/climax"
)

var (
	dataDir        string
	schemaFilename string
	// this is just a cheat because all the initial date-times are the same
	// if they actually varied, I would probably use the time.Parse func
	marchFirst = time.Date(2015, time.Month(3), 1, 0, 0, 0, 0, time.UTC)
)

func init() {
	_, filename, _, _ := runtime.Caller(1)
	var err error
	dataDir, err = filepath.Abs(path.Join(path.Dir(filename), "../data"))
	if err != nil {
		log.Fatalln(err)
	}
	schemaFilename = path.Join(dataDir, "schema.sql")
}

func loadCSV(filename string) ([][]string, error) {
	f, err := os.Open(path.Join(dataDir, filename))
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("could open data file: %s", err)
	}
	reader := csv.NewReader(f)
	return reader.ReadAll()
}

// shortcut for fmt.Errorf...
func quickErr(fmtString string, err error) error {
	return fmt.Errorf("Could not read schema: %s", err)
}

// CreateInitialDatabase is a utility function to create a new database and
// load all the initial csv data
func CreateInitialDatabase(filename string) error {
	err := CreateDb(filename, false)
	if err != nil {
		return err
	}
	return loadInitalData(filename)
}

// CreateDb will create the database and load the initial schema.
// It will error if the database already exists, unless `force` is `true`
func CreateDb(filename string, force bool) error {
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
		return fmt.Errorf("Could not read schema: %s", err)
	}
	schema := string(schemaBytes)

	// connect to database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return fmt.Errorf("Could not open database: %s", err)
	}
	defer db.Close()

	// execure schema
	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("Could not execute schema: %s", err)
	}

	return nil
}

func loadInitalData(filename string) error {
	// connect to database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return fmt.Errorf("Could not open database: %s", err)
	}
	defer db.Close()

	stateData, err := loadCSV("State.csv")
	if err != nil {
		return fmt.Errorf("Could not load state data: %s", err)
	}
	for _, record := range stateData[1:] {
		_, err = db.Exec(
			`INSERT INTO states (name, abbrev, created_at, updated_at)
			VALUES (?, ?, ?, ?)`, record[0], record[1], marchFirst, marchFirst)
		if err != nil {
			return err
		}
	}

	cityData, err := loadCSV("City.csv")
	if err != nil {
		return fmt.Errorf("Could not load city data: %s", err)
	}
	for _, record := range cityData[1:] {
		lat, err := strconv.ParseFloat(record[3], 64)
		lon, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return fmt.Errorf("Could not convert lat/lon to float: %s", err)
		}

		latSin, latCos, lonSin, lonCos := geo.LatLonSinCos(lat, lon)
		_, err = db.Exec(
			`INSERT INTO cities (
				name, state_id, lat, lon, lat_sin, lat_cos, lon_sin, lon_cos,
				created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			record[0], record[1], lat, lon,
			latSin, latCos, lonSin, lonCos,
			marchFirst, marchFirst)
		if err != nil {
			return err
		}
	}

	userData, err := loadCSV("User.csv")
	if err != nil {
		return fmt.Errorf("Could not load user data: %s", err)
	}
	for _, record := range userData[1:] {
		_, err = db.Exec(
			`INSERT INTO users (
				first_name, last_name, created_at, updated_at)
			VALUES (?, ?, ?, ?)`,
			record[0], record[1], marchFirst, marchFirst)
		if err != nil {
			return err
		}
	}

	return nil
}

func initDB(ctx climax.Context) int {
	handleDebugging(ctx)
	cfg := conf.LoadFile(getConfigPath(ctx))

	err := CreateDb(cfg.DBPath, ctx.Is("force-recreate-database"))
	if err != nil {
		log.Errorln(err)
		return 1
	}

	err = loadInitalData(cfg.DBPath)
	if err != nil {
		log.Errorln(err)
		return 2
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
