package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// State model*
type State struct {
	gorm.Model
	Name   string
	Abbrev string
}

// City model
type City struct {
	gorm.Model
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
	gorm.Model
	FirstName string
	LastName  string

	Email        string
	PasswordHash []byte
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
	gorm.Model

	User   User `json:"-"`
	UserID uint
	City   City `json:"-"`
	CityID uint

	Lat, Lon       float64
	LatSin, LatCos float64
	LonSin, LonCos float64
	VisitMethod    string
}
