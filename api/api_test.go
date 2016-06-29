package api_test

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	. "github.com/bobisme/RestApiProject/api"
	"github.com/bobisme/RestApiProject/cmd"
	"github.com/bobisme/RestApiProject/conf"
	"github.com/bobisme/RestApiProject/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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
	_, err = db.Exec(
		`INSERT INTO states (name, abbrev, created_at, updated_at)
		VALUES (?, ?, ?, ?)`, "Essos", "ES", marchFirst, marchFirst)
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
		`INSERT INTO cities (name, state_id, lat, lon,
			lat_sin, lat_cos, lon_sin, lon_cos, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"Kings Landing", 1, 32.7765, -79.9311, 0.54136, 0.840789, -0.984598, 0.17483,
		marchFirst, marchFirst)
	check(err)
	_, err = db.Exec(
		`INSERT INTO cities (name, state_id, lat, lon,
			lat_sin, lat_cos, lon_sin, lon_cos, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"Qarth", 2, 26.8206, 30.8025, 0.45120, 0.892424, 0.51208, 0.858938,
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

func getRespBody(resp *http.Response) []byte {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	Ω(err).ShouldNot(HaveOccurred())
	return body
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
		Ω(err).ShouldNot(HaveOccurred())
		return getRespBody(resp)
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

	Context("new visit", func() {
		It("should be ok", func() {
			req := strings.NewReader(`{
				"city": "Winterfell",
				"state": "WS"
			}`)
			resp, err := http.Post(
				ts.URL+`/user/1/visits`, "application/json", req)
			Ω(err).ShouldNot(HaveOccurred())
			// body := getRespBody(resp)
			Ω(resp.StatusCode).Should(Equal(201))
			var visit models.Visit
			db.Model(&models.Visit{}).First(&visit)
			var visitCount int
			Ω(visit.CityID).Should(Equal(uint(1)))
			Ω(visit.UserID).Should(Equal(uint(1)))
			db.Model(&models.Visit{}).Count(&visitCount)
			Ω(visitCount).Should(Equal(1))

			req2 := strings.NewReader(`{
				"city": "Kings Landing",
				"state": "WS"
			}`)
			resp2, err := http.Post(
				ts.URL+`/user/1/visits`, "application/json", req2)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp2.StatusCode).Should(Equal(201))
			db.Model(&models.Visit{}).Count(&visitCount)
			Ω(visitCount).Should(Equal(2))

			req3 := strings.NewReader(`{
				"city": "Kings Landing",
				"state": "WS"
			}`)
			resp3, err := http.Post(
				ts.URL+`/user/1/visits`, "application/json", req3)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp3.StatusCode).Should(Equal(201))
			db.Model(&models.Visit{}).Count(&visitCount)
			Ω(visitCount).Should(Equal(3))
		})

		DescribeTable("fails on invalid user",
			func(url string) {
				reqData := `{ "city": "Winterfell", "state": "WS" }`
				req := strings.NewReader(reqData)
				resp, err := http.Post(ts.URL+url, "application/json", req)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(resp.StatusCode).Should(Equal(400))
			},
			Entry("0", `/user/0/visits`),
			Entry("non-existant", `/user/20/visits`),
			Entry("not a number", `/user/NO/visits`),
			Entry("blank", `/user//visits`),
		)

		DescribeTable("fails on invalid city",
			func(reqData string) {
				req := strings.NewReader(reqData)
				resp, err := http.Post(
					ts.URL+`/user/1/visits`, "application/json", req)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(resp.StatusCode).Should(Equal(400))
			},
			Entry("Not a city", `{ "city": "That one place", "state": "WS" }`),
			Entry("Wrong state", `{ "city": "Winterfell", "state": "ES" }`),
		)

		It("fails on invalid state", func() {
			reqData := `{ "city": "Winterfall", "state": "XS" }`
			req := strings.NewReader(reqData)
			resp, err := http.Post(ts.URL+`/user/1/visits`, "application/json", req)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp.StatusCode).Should(Equal(400))
		})
	})
})
