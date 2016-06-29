package conf_test

import (
	. "."
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("conf", func() {
	Describe("Default()", func() {
		c := Default()
		It("should set db_path to database.sqlite3", func() {
			Ω(c.DBPath).Should(Equal("database.sqlite3"))
		})
	})
})
