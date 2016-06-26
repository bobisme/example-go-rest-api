package models_test

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/bobisme/RestApiProject/cmd"
	. "github.com/bobisme/RestApiProject/models"

	"os"

	"database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // load sqlite3 support
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
		VALUES (?, ?, ?, ?)`, "North Carolina", "NC", marchFirst, marchFirst)
	check(err)
	_, err = db.Exec(
		`INSERT INTO states (name, abbrev, created_at, updated_at)
		VALUES (?, ?, ?, ?)`, "South Carolina", "SC", marchFirst, marchFirst)
	check(err)
	_, err = db.Exec(
		`UPDATE states SET deleted_at = ? WHERE abbrev = ?`, marchFirst, "SC")
	check(err)
	// city
	_, err = db.Exec(
		`INSERT INTO cities (name, state_id, lat, lon,
			lat_sin, lat_cos, lon_sin, lon_cos, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"Charlotte", 1, 35.2271, -80.8431,
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

var _ = Describe("Models", func() {
	var db *gorm.DB

	BeforeEach(func() {
		var err error
		Ω(cmd.CreateDb("test-models.sqlite3", true)).Should(Succeed())
		loadTestData("test-models.sqlite3")
		db, err = gorm.Open("sqlite3", "test-models.sqlite3")
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		db.Close()
		os.Remove("test-models.sqlite3")
	})

	Context("State", func() {
		var nc State
		var sc State

		BeforeEach(func() {
			db.Where("abbrev = ?", "NC").First(&nc)
			db.Where("abbrev = ?", "SC").First(&sc)
		})

		It("loads ID correctly", func() {
			Ω(nc.Model.ID).Should(Equal(uint(1)))
		})

		It("loads name correctly", func() {
			Ω(nc.Name).Should(Equal("North Carolina"))
		})

		It("loads abbrev correctly", func() {
			Ω(nc.Abbrev).Should(Equal("NC"))
		})

		It("loads created correctly", func() {
			Ω(nc.CreatedAt).Should(Equal(marchFirst))
		})

		It("loads updated correctly", func() {
			Ω(nc.UpdatedAt).Should(Equal(marchFirst))
		})

		It("loads deleted correctly", func() {
			Ω(nc.DeletedAt).Should(BeNil())
		})

		It("SC was soft deleted", func() {
			Ω(sc.ID).Should(Equal(uint(0)))
			var scCount int
			db.Model(&State{}).Where("abbrev = ?", "SC").Count(&scCount)
			Ω(scCount).Should(Equal(0))
			// this does a *more* raw query, so it shows SC in the results
			db.Table("states").Where("abbrev = ?", "SC").Count(&scCount)
			Ω(scCount).Should(Equal(1))
		})
	})

	Context("City", func() {
		var charlotte City

		BeforeEach(func() {
			db.Where("name = ?", "Charlotte").First(&charlotte)
		})

		It("loads ID correctly", func() {
			Ω(charlotte.ID).Should(Equal(uint(1)))
		})

		It("loads name correctly", func() {
			Ω(charlotte.Name).Should(Equal("Charlotte"))
		})

		It("loads state_id correctly", func() {
			Ω(charlotte.StateID).Should(Equal(uint(1)))
		})

		It("loads state correctly", func() {
			db.Model(&charlotte).Related(&charlotte.State)
			Ω(charlotte.State.ID).Should(Equal(uint(1)))
			Ω(charlotte.State.Name).Should(Equal("North Carolina"))
		})

		It("loads lat, long fields", func() {
			Ω(charlotte.Lat).Should(BeNumerically("~", 35.2271))
			Ω(charlotte.Lon).Should(BeNumerically("~", -80.8431))
			Ω(charlotte.LatSin).Should(BeNumerically("~", 0.57681874832))
			Ω(charlotte.LatCos).Should(BeNumerically("~", 0.81687216354))
			Ω(charlotte.LonSin).Should(BeNumerically("~", -0.9872562543))
			Ω(charlotte.LonCos).Should(BeNumerically("~", 0.15913858219))
		})

		It("loads created correctly", func() {
			Ω(charlotte.Model.CreatedAt).Should(Equal(marchFirst))
		})

		It("loads updated correctly", func() {
			Ω(charlotte.Model.UpdatedAt).Should(Equal(marchFirst))
		})

		It("loads deleted correctly", func() {
			Ω(charlotte.Model.DeletedAt).Should(BeNil())
		})
	})

	Context("User", func() {
		var snow User

		BeforeEach(func() {
			db.Where("first_name = ? AND last_name = ?", "John", "Snow").First(&snow)
		})

		It("loads ID correctly", func() {
			Ω(snow.ID).Should(Equal(uint(1)))
		})

		It("loads first name correctly", func() {
			Ω(snow.FirstName).Should(Equal("John"))
		})

		It("loads last name correctly", func() {
			Ω(snow.LastName).Should(Equal("Snow"))
		})

		It("loads email correctly", func() {
			Ω(snow.Email).Should(Equal("john@northernbastards.net"))
		})

		It("loads password_hash correctly", func() {
			Ω(snow.PasswordHash).Should(Equal([]byte("HASHNOTTESTEDHERE")))
		})

		It("loads created correctly", func() {
			Ω(snow.Model.CreatedAt).Should(Equal(marchFirst))
		})

		It("loads updated correctly", func() {
			Ω(snow.Model.UpdatedAt).Should(Equal(marchFirst))
		})

		It("loads deleted correctly", func() {
			Ω(snow.Model.DeletedAt).Should(BeNil())
		})
	})

	Context("User passwords", func() {
		var daeny User

		BeforeEach(func() {
			daeny = User{
				FirstName: "Daenerys",
				LastName:  "Targaryen",
				Email:     "mhysa@khaleesi.org",
			}
			db.Create(&daeny)
		})

		AfterEach(func() {
			db.Delete(&daeny)
		})

		It("doesn't take blank passwords", func() {
			Ω(SetPassword(db, &daeny, "")).ShouldNot(Succeed())
		})

		It("sets hashes", func() {
			Ω(SetPassword(db, &daeny, "drogon84")).Should(Succeed())
			var u User
			db.Where("email = ?", "mhysa@khaleesi.org").First(&u)
			Ω(bcrypt.CompareHashAndPassword(
				u.PasswordHash, []byte("drogon84"))).Should(Succeed())
		})

		It("checks passwords", func() {
			Ω(SetPassword(db, &daeny, "drogon84")).Should(Succeed())
			var u User
			db.Where("email = ?", "mhysa@khaleesi.org").First(&u)
			Ω(CheckPassword(db, &u, "drogon84")).Should(Succeed())
			Ω(CheckPassword(db, &u, "drogon85")).ShouldNot(Succeed())
		})
	})
})
