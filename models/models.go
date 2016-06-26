package models

import "github.com/jinzhu/gorm"

// State model*
type State struct {
	gorm.Model
	Name   string
	Abbrev string
}

// City model
type City struct {
	gorm.Model
	Name           string
	State          State
	StateID        uint
	Lat, Lon       float64
	LatSin, LatCos float64
	LonSin, LonCos float64
}

// User model
type User struct {
	gorm.Model
	FirstName string
	LastName  string

	Email        string
	PasswordHash string
}

// Visit model
type Visit struct {
	gorm.Model
	User User
	City City

	Lat, Lon       float64
	LatSin, LatCos float64
	LonSin, LonCos float64
	VisitMethod    string
}
