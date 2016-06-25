package cmd

import (
	"database/sql"
	"os"

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
			err := createDb("test-database.db", false)
			Ω(err).ShouldNot(HaveOccurred())
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
		})

	})
})
