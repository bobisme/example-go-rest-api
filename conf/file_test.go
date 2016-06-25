package conf_test

import (
	"os"
	"path/filepath"

	. "."
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	missingFile, _ = filepath.Abs("test-missing-config.toml")
	blankFile, _   = filepath.Abs("test-blank-config.toml")
	normalFile, _  = filepath.Abs("test-non-blank-config.toml")
)

func createTestFiles() {
	f, err := os.Create(blankFile)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	f2, err := os.Create(normalFile)
	defer f2.Close()
	f2.WriteString(`db_path = "test.db"`)
	if err != nil {
		panic(err)
	}
}

var _ = Describe("File", func() {
	Describe("LoadFile", func() {
		BeforeEach(func() {
			os.Remove(missingFile)
			createTestFiles()
		})

		AfterEach(func() {
			os.Remove(blankFile)
			os.Remove(normalFile)
		})

		It("should panic on invalid path", func() {
			fatalFn := func() { LoadFile(missingFile) }
			Ω(fatalFn).Should(Panic())
		})

		It("should return a default config if file blank", func() {
			c := LoadFile(blankFile)
			d := Default()
			Ω(c).Should(Equal(d))
		})

		It("should return the proper config settings", func() {
			c := LoadFile(normalFile)
			d := Default()
			Ω(c).ShouldNot(Equal(d))
			Ω(c.DBPath).Should(Equal("test.db"))
		})
	})
})
