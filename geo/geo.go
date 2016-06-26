// Package geo just has some utility functions for trigonometry and geometry
package geo

import "math"

// DegToRad converts degrees to radians
func DegToRad(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// LatLonSinCos takes latitude and longitude in degrees and returns
// sin(lat), cos(lat), sin(lon), cos(lon)
func LatLonSinCos(lat, lon float64) (float64, float64, float64, float64) {
	latRad := DegToRad(lat)
	lonRad := DegToRad(lon)
	return math.Sin(latRad), math.Cos(latRad), math.Sin(lonRad), math.Cos(lonRad)
}
