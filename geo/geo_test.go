package geo_test

import (
	"math"

	. "github.com/bobisme/RestApiProject/geo"
	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Geo", func() {
	Describe("DegToRad", func() {
		DescribeTable(
			"works",
			func(deg, expected float64) {
				Ω(DegToRad(deg)).Should(BeNumerically("~", expected))
			},
			Entry("positive number", 35.2271, 0.6148288809),
			Entry("negative number", -80.8431, -1.41097827),
			Entry("0°", 0.0, 0.0),
			Entry("180°", 180.0, math.Pi),
			Entry("360°", 360.0, 2*math.Pi),
			Entry("-180°", -180.0, -math.Pi),
			Entry("-360°", -360.0, -2*math.Pi),
		)
	})
	Describe("LatLonSinCos", func() {
		It("works", func() {
			latSin, latCos, lonSin, lonCos := LatLonSinCos(35.2271, -80.8431)
			Ω(latSin).Should(BeNumerically("~", 0.57681874832))
			Ω(latCos).Should(BeNumerically("~", 0.81687216354))
			Ω(lonSin).Should(BeNumerically("~", -0.9872562543))
			Ω(lonCos).Should(BeNumerically("~", 0.15913858219))
		})
	})
})
