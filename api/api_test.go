package api_test

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	. "github.com/bobisme/RestApiProject/api"
	"github.com/bobisme/RestApiProject/cmd"
	"github.com/bobisme/RestApiProject/conf"
	"github.com/bobisme/RestApiProject/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var marchFirst = time.Date(2015, time.Month(3), 1, 0, 0, 0, 0, time.UTC)

func loadTestData(filename string) {
	check := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	// connect to database
	db, err := sql.Open("sqlite3", filename)
	check(err)
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO states (name, abbrev, created_at, updated_at)
		VALUES (?, ?, ?, ?)`, "Westeros", "WS", marchFirst, marchFirst)
	check(err)
	// city
	_, err = db.Exec(
		`INSERT INTO cities (name, state_id, lat, lon,
			lat_sin, lat_cos, lon_sin, lon_cos, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"Winterfell", 1, 35.2271, -80.8431,
		0.57681874832, 0.81687216354, -0.9872562543, 0.15913858219,
		marchFirst, marchFirst)
	check(err)
	_, err = db.Exec(
		`INSERT INTO users (
			first_name, last_name, email, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		"John", "Snow", "john@northernbastards.net", "HASHNOTTESTEDHERE",
		marchFirst, marchFirst)
	check(err)
}

var _ = Describe("Api", func() {
	var (
		r   *gin.Engine
		ts  *httptest.Server
		cfg *conf.Config
		db  *gorm.DB
	)

	get := func(url string) []byte {
		resp, err := http.Get(ts.URL + url)
		defer resp.Body.Close()
		Ω(err).ShouldNot(HaveOccurred())
		body, err := ioutil.ReadAll(resp.Body)
		Ω(err).ShouldNot(HaveOccurred())
		return body
	}

	getJSON := func(url string, out interface{}) {
		json.Unmarshal(get(url), out)
	}

	BeforeEach(func() {
		var err error
		cfg = conf.Default()
		cfg.DBPath = "test-rest-api.db"
		cmd.CreateDb("test-rest-api.db", true)
		loadTestData("test-rest-api.db")
		db, err = gorm.Open("sqlite3", "test-rest-api.db")
		Ω(err).ShouldNot(HaveOccurred())
		r = GetRouter(cfg, db)
		ts = httptest.NewServer(r)
	})

	AfterEach(func() {
		ts.Close()
		ts = nil
		r = nil
		cfg = nil
		db.Close()
		os.Remove("test-rest-api.db")
	})

	It("handles the root fine", func() {
		body := get("/")
		Ω(body).Should(Equal([]byte("HELLO")))
	})

	Context("cities in state", func() {
		It("should work", func() {
			var out []models.City
			getJSON("/state/1/cities", &out)
			Ω(len(out)).ShouldNot(Equal(0))
			Ω(out[0].Name).Should(Equal("Winterfell"))
			Ω(out[0].Lat).Should(BeNumerically("~", 35.2271))
			Ω(out[0].Lon).Should(BeNumerically("~", -80.8431))
			Ω(out[0].StateID).Should(Equal(uint(1)))
		})
		It("should err if not an id", func() {
			resp, err := http.Get(ts.URL + "/state/NO/cities")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp.StatusCode).Should(Equal(400))
		})
	})
})
