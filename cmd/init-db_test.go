package cmd

import (
	"database/sql"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("InitDB", func() {
	Describe("createDb", func() {
		AfterEach(func() {
			os.Remove("test-database.db")
		})

		It("will create the database file", func() {
			os.Remove("test-database.db")
			Ω(createDb("test-database.db", false)).Should(Succeed())
		})

		It("errors if database exists", func() {
			f, err := os.Create("test-database.db")
			f.Close()
			Ω(err).ShouldNot(HaveOccurred())
			err = createDb("test-database.db", false)
			Ω(err).Should(HaveOccurred())
		})

		It("will recreate the database if forced", func() {
			f, err := os.Create("test-database.db")
			f.Close()
			Ω(err).ShouldNot(HaveOccurred())
			err = createDb("test-database.db", true)
			Ω(err).ShouldNot(HaveOccurred())
		})

		Context("schema", func() {
			var db *sql.DB
			BeforeEach(func() {
				Ω(createDb("test-database.db", true)).Should(Succeed())
				Ω(loadInitalData("test-database.db")).Should(Succeed())
				db, _ = sql.Open("sqlite3", "test-database.db")
			})

			AfterEach(func() {
				db.Close()
			})

			DescribeTable(
				"creates tables",
				func(name string) {
					var tableName string
					err := db.QueryRow(
						`SELECT name FROM sqlite_master
						WHERE type='table' AND name=?`, name).Scan(&tableName)
					Ω(err).ShouldNot(HaveOccurred())
				},
				Entry("states", "states"),
				Entry("cities", "cities"),
				Entry("users", "users"),
				Entry("visits", "visits"),
			)

			DescribeTable(
				"loads initial data",
				func(tableName string, expectedCount int) {
					var count int
					err := db.QueryRow(
						`SELECT COUNT(*) AS count FROM ` + tableName).Scan(&count)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(count).Should(Equal(expectedCount))
				},
				// 51 "states" because DC is included
				// tough luck Puerto Rico and Guam
				Entry("states", "states", 51),
				Entry("cities", "cities", 505),
				Entry("users", "users", 9),
			)

			It("loads states correctly", func() {
				state := struct {
					Name, Abbrev      string
					Created, Modified time.Time
				}{}
				err := db.QueryRow(
					`SELECT name, abbrev, created_at, updated_at
					FROM states WHERE id=1`,
				).Scan(&state.Name, &state.Abbrev, &state.Created, &state.Modified)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(state.Name).Should(Equal("Alabama"))
				Ω(state.Abbrev).Should(Equal("AL"))
				Ω(state.Created).Should(Equal(marchFirst))
				Ω(state.Modified).Should(Equal(marchFirst))
			})

			It("loads cities correctly", func() {
				d := struct {
					name                                     string
					stateID                                  int
					lat, lon, latSin, latCos, lonSin, lonCos float64
					created, modified                        time.Time
				}{}
				err := db.QueryRow(
					`SELECT name, state_id, lat, lon,
					lat_sin, lat_cos, lon_sin, lon_cos,
					created_at, updated_at
					FROM cities WHERE id=1`,
				).Scan(
					&d.name, &d.stateID, &d.lat, &d.lon,
					&d.latSin, &d.latCos, &d.lonSin, &d.lonCos,
					&d.created, &d.modified,
				)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(d.name).Should(Equal("Akron"))
				Ω(d.stateID).Should(Equal(1))
				Ω(d.lat).Should(Equal(32.87802))
				Ω(d.lon).Should(Equal(-87.743989))
				Ω(d.latSin).Should(BeNumerically("~", 0.54285231218))
				Ω(d.latCos).Should(BeNumerically("~", 0.83982817716))
				Ω(d.lonSin).Should(BeNumerically("~", -0.99922491192))
				Ω(d.lonCos).Should(BeNumerically("~", 0.0393646464))
				Ω(d.created).Should(Equal(marchFirst))
				Ω(d.modified).Should(Equal(marchFirst))
			})

			It("loads users correctly", func() {
				d := struct {
					firstName, lastName        string
					created, modified, deleted *time.Time
				}{}
				err := db.QueryRow(
					`SELECT first_name, last_name,
					created_at, updated_at, deleted_at
					FROM users WHERE id=1`,
				).Scan(
					&d.firstName, &d.lastName,
					&d.created, &d.modified, &d.deleted,
				)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(d.firstName).Should(Equal("Henry"))
				Ω(d.lastName).Should(Equal("Harrison"))
				Ω(*d.created).Should(Equal(marchFirst))
				Ω(*d.modified).Should(Equal(marchFirst))
				Ω(d.deleted).Should(BeNil())
			})
		})

	})
})
