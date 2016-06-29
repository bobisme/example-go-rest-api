package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Model is a better base model because it aknowledges JSON
type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

// State model*
type State struct {
	Model
	Name   string `json:"name"`
	Abbrev string `json:"abbrev"`
}

// City model
type City struct {
	Model
	Name    string  `json:"name"`
	State   State   `json:"-"`
	StateID uint    `json:"stateId"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	LatSin  float64 `json:"-"`
	LatCos  float64 `json:"-"`
	LonSin  float64 `json:"-"`
	LonCos  float64 `json:"-"`
}

// User model
type User struct {
	Model
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Email        string  `json:"email"`
	PasswordHash []byte  `json:"-"`
	Visits       []Visit `json:"visits"`
}

// SetPassword for the user
func SetPassword(db *gorm.DB, user *User, password string) error {
	if password == "" {
		return fmt.Errorf("Password must not be empty.")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}
	user.PasswordHash = hash
	if err = db.Model(user).Update("password_hash", hash).Error; err != nil {
		return err
	}
	return nil
}

// CheckPassword for the user. Returns nil on success, error otherwise.
func CheckPassword(db *gorm.DB, user *User, password string) error {
	return bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
}

// Visit model
type Visit struct {
	Model

	User   User `json:"-"`
	UserID uint `json:"userId"`
	City   City `json:"-"`
	CityID uint `json:"cityId"`

	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	LatSin, LatCos float64 `json:"-"`
	LonSin, LonCos float64 `json:"-"`
	VisitMethod    string  `json:"visitMethod"`
}
