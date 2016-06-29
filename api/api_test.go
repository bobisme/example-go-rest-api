package api_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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
	gin.SetMode(gin.ReleaseMode)

	var (
		r   *gin.Engine
		ts  *httptest.Server
		cfg *conf.Config
		db  *gorm.DB
	)

	get := func(url string) []byte {
		resp, err := http.Get(ts.URL + url)
		Ω(err).ShouldNot(HaveOccurred())
		body := getRespBody(resp)
		if resp.StatusCode != 200 {
			fmt.Println("ERROR >>>>>", string(body))
		}
		Ω(resp.StatusCode).Should(Equal(200))
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
		// don't use the default router. it's too noisy
		r = gin.New()
		SetRoutes(cfg, db, r)
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
			var out struct {
				Limit, Offset, Count int
				Data                 []models.City
			}
			getJSON("/state/1/cities", &out)
			Ω(len(out.Data)).Should(Equal(2))
			Ω(out.Limit).Should(Equal(100))
			Ω(out.Offset).Should(Equal(0))
			Ω(out.Count).Should(Equal(2))
			Ω(out.Data[0].Name).Should(Equal("Winterfell"))
			Ω(out.Data[0].Lat).Should(BeNumerically("~", 35.2271))
			Ω(out.Data[0].Lon).Should(BeNumerically("~", -80.8431))
			Ω(out.Data[0].StateID).Should(Equal(uint(1)))
		})

		It("limit works", func() {
			var out struct {
				Limit, Offset, Count int
				Data                 []models.City
			}
			getJSON("/state/1/cities?limit=1", &out)
			Ω(len(out.Data)).Should(Equal(1))
			Ω(out.Limit).Should(Equal(1))
			Ω(out.Offset).Should(Equal(0))
			Ω(out.Count).Should(Equal(2))
			Ω(out.Data[0].ID).Should(Equal(uint(1)))
		})

		It("offset works", func() {
			var out struct {
				Limit, Offset, Count int
				Data                 []models.City
			}
			getJSON("/state/1/cities?offset=1", &out)
			Ω(len(out.Data)).Should(Equal(1))
			Ω(out.Limit).Should(Equal(100))
			Ω(out.Offset).Should(Equal(1))
			Ω(out.Count).Should(Equal(2))
			Ω(out.Data[0].ID).Should(Equal(uint(2)))
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
			Ω(visit.CityID).Should(Equal(uint(1)))
			Ω(visit.UserID).Should(Equal(uint(1)))
			var visitCount int
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

	Context("delete visit", func() {
		var ids []uint
		BeforeEach(func() {
			ids = []uint{}
			visits := []string{
				`{ "city": "Winterfell", "state": "WS" }`,
				`{ "city": "Qarth", "state": "ES" }`,
				`{ "city": "Kings Landing", "state": "WS" }`,
			}
			for _, data := range visits {
				req := strings.NewReader(data)
				resp, _ := http.Post(
					ts.URL+`/user/1/visits`, "application/json", req)
				body := getRespBody(resp)
				var visit models.Visit
				json.Unmarshal(body, &visit)
				ids = append(ids, visit.ID)
			}
		})

		It("should be ok", func() {
			var visitCount int
			db.Model(&models.Visit{}).Count(&visitCount)
			Ω(visitCount).Should(Equal(3))

			id := strconv.Itoa(int(ids[1]))
			req, err := http.NewRequest("DELETE", ts.URL+`/user/1/visits/`+id, nil)
			Ω(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(resp.StatusCode).Should(Equal(204))
			db.Model(&models.Visit{}).Count(&visitCount)
			Ω(visitCount).Should(Equal(2))
		})

		DescribeTable("fails on invalid user",
			func(url string) {
				req, _ := http.NewRequest("DELETE", ts.URL+url, nil)
				resp, _ := http.DefaultClient.Do(req)
				Ω(resp.StatusCode).Should(Equal(400))
			},
			Entry("0", `/user/0/visits/2`),
			Entry("non-existant", `/user/20/visits/2`),
			Entry("not a number", `/user/NO/visits/2`),
			Entry("blank", `/user//visits/2`),
		)

		DescribeTable("fails on invalid visit",
			func(url string) {
				req, _ := http.NewRequest("DELETE", ts.URL+url, nil)
				resp, _ := http.DefaultClient.Do(req)
				Ω(resp.StatusCode).Should(Equal(400))
			},
			Entry("0", `/user/1/visits/0`),
			Entry("non-existant", `/user/1/visits/20`),
			Entry("not a number", `/user/1/visits/NO`),
		)
	})
	Context("cities visited", func() {
		BeforeEach(func() {
			visits := []string{
				`{ "city": "Winterfell", "state": "WS" }`,
				`{ "city": "Qarth", "state": "ES" }`,
				`{ "city": "Kings Landing", "state": "WS" }`,
				`{ "city": "Winterfell", "state": "WS" }`,
				`{ "city": "Winterfell", "state": "WS" }`,
				`{ "city": "Kings Landing", "state": "WS" }`,
				`{ "city": "Qarth", "state": "ES" }`,
			}
			for _, data := range visits {
				req := strings.NewReader(data)
				resp, _ := http.Post(
					ts.URL+`/user/1/visits`, "application/json", req)
				body := getRespBody(resp)
				var visit models.Visit
				json.Unmarshal(body, &visit)
			}
		})

		It("should be ok", func() {
			var out struct {
				Limit, Offset, Count int
				Data                 []models.City
			}
			body := get("/user/1/visits")
			json.Unmarshal(body, &out)
			Ω(len(out.Data)).Should(Equal(3))
			Ω(out.Limit).Should(Equal(100))
			Ω(out.Offset).Should(Equal(0))
			Ω(out.Count).Should(Equal(3))
			Ω(out.Data[0].Name).Should(Equal("Winterfell"))
			Ω(out.Data[0].Lat).Should(BeNumerically("~", 35.2271))
			Ω(out.Data[0].Lon).Should(BeNumerically("~", -80.8431))
			Ω(out.Data[0].StateID).Should(Equal(uint(1)))
		})

		It("accepts limit", func() {
			var out struct {
				Limit, Offset, Count int
				Data                 []models.City
			}
			body := get("/user/1/visits?limit=2")
			json.Unmarshal(body, &out)
			Ω(len(out.Data)).Should(Equal(2))
			Ω(out.Limit).Should(Equal(2))
			Ω(out.Offset).Should(Equal(0))
			Ω(out.Count).Should(Equal(3))
			Ω(out.Data[0].ID).Should(Equal(uint(1)))
		})

		It("accepts offset", func() {
			var out struct {
				Limit, Offset, Count int
				Data                 []models.City
			}
			body := get("/user/1/visits?offset=1")
			json.Unmarshal(body, &out)
			Ω(len(out.Data)).Should(Equal(2))
			Ω(out.Limit).Should(Equal(100))
			Ω(out.Offset).Should(Equal(1))
			Ω(out.Count).Should(Equal(3))
			Ω(out.Data[0].ID).Should(Equal(uint(2)))
		})
	})
	Context("states visited", func() {
		BeforeEach(func() {
			visits := []string{
				`{ "city": "Winterfell", "state": "WS" }`,
				`{ "city": "Qarth", "state": "ES" }`,
				`{ "city": "Kings Landing", "state": "WS" }`,
				`{ "city": "Winterfell", "state": "WS" }`,
				`{ "city": "Winterfell", "state": "WS" }`,
				`{ "city": "Kings Landing", "state": "WS" }`,
				`{ "city": "Qarth", "state": "ES" }`,
			}
			for _, data := range visits {
				req := strings.NewReader(data)
				resp, _ := http.Post(
					ts.URL+`/user/1/visits`, "application/json", req)
				body := getRespBody(resp)
				var visit models.Visit
				json.Unmarshal(body, &visit)
			}
		})
		It("should be ok", func() {
			var out []models.State
			body := get("/user/1/visits/states")
			json.Unmarshal(body, &out)
			Ω(len(out)).Should(Equal(2))
			Ω(out[0].Name).Should(Equal("Westeros"))
			Ω(out[0].Abbrev).Should(Equal("WS"))
			Ω(out[0].ID).Should(Equal(uint(1)))
		})
	})
})
