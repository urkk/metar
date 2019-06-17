package conversion

import (
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLengthConversionTests(t *testing.T) {

	Convey("Length conversion tests", t, func() {
		Convey("converts statute miles to meters", func() {
			for input := .0; input < 6.25; {
				So(SMileToM(input), ShouldAlmostEqual, input*1609.344, 1)
				input += .25
			}
		})

		Convey("converts feet to meters (rounded to 10 meters)", func() {
			for input := 300; input < 3000; {
				So(FtToM(input), ShouldAlmostEqual, int(math.Round(float64(input)*0.3048/10)*10), 10)
				input += 300
			}
		})

		Convey("converts meters to feet", func() {
			for input := 0; input < 8000; {
				So(MToFt(input), ShouldAlmostEqual, float64(input)*3.28084, 10)
				input += 100
			}
		})

		Convey("converts metres to statute miles", func() {
			for input := 0; input < 8000; {
				So(MToSMile(input), ShouldAlmostEqual, float64(input)/1609.344, 0.1)
				input += 500
			}
		})

		Convey("converts feet to statute miles", func() {
			for input := 0; input < 8000; {
				So(FtToSMile(input), ShouldEqual, float64(input)/5280)
				input += 500
			}
		})

		Convey("converts statute miles to feet", func() {
			for input := .0; input < 6; {
				So(SMileToFt(input), ShouldEqual, input*5280)
				input += .25
			}
		})

	})

}

func TestSpeedConversionTests(t *testing.T) {
	Convey("Speed conversion tests", t, func() {

		Convey("converts kilometres per hour to meters per second", func() {
			for input := 0; input < 200; {
				So(KphToMps(input), ShouldAlmostEqual, float64(input)/3.6, .00001)
				input += 5
			}
		})

		Convey("converts knots to meters per second", func() {
			for input := .0; input < 50; input++ {
				So(KtsToMps(input), ShouldAlmostEqual, input/1.94384, .00001)
			}
		})

		Convey("converts meters per second to knots", func() {
			for input := .0; input < 50; input++ {
				So(MpsToKts(input), ShouldAlmostEqual, input*1.94384, .00001)
			}
		})

		Convey("converts kilometres per hour to knots", func() {
			for input := 0; input < 120; {
				So(KphToKts(input), ShouldAlmostEqual, float64(input)/1.852, 0.1)
				input += 2
			}
		})
	})
}

func TestPressureConversionTests(t *testing.T) {
	Convey("Pressure conversion tests", t, func() {
		Convey("converts hectopascal to mm of mercury", func() {
			for input := 985; input < 1027; input++ {
				So(HPaToMmHg(input), ShouldAlmostEqual, float64(input)*0.75006375541921, 1)
			}
		})

		Convey("converts mm of mercury to hectopascal", func() {
			for input := 740; input < 780; input++ {
				So(MmHgToHPa(input), ShouldAlmostEqual, float64(input)*1.333223684, 1)
			}
		})

		Convey("converts inch of mercury to hectopascal", func() {
			for input := 29.0; input < 30.3; {
				So(InHgTohPa(input), ShouldAlmostEqual, input*33.86389, 1)
				input += .1
			}
		})
	})
}

func TestDirectionToCardinalDirection(t *testing.T) {

	type testpairDirection struct {
		input    int
		expected string
	}

	var testsDirection = []testpairDirection{
		{360, "N"},
		{30, "NE"},
		{275, "W"},
		{0, "N"},
	}

	Convey("converts direction in degrees to points of the compass", t, func() {
		for _, pair := range testsDirection {
			So(DirectionToCardinalDirection(pair.input), ShouldEqual, pair.expected)
		}
	})
}

func TestCalcRelativeHumidity(t *testing.T) {
	type testRelativeHumidity struct {
		temp, dewpoint, rh int
	}
	var testsRelativeHumidity = []testRelativeHumidity{
		{20, 7, 43},
		{31, 7, 22},
		{25, 25, 100},
		{25, 19, 69},
		{-7, -13, 62},
		{0, -9, 51},
	}
	Convey("calculates the relative humidity of the dew point and temperature", t, func() {
		for _, data := range testsRelativeHumidity {
			So(CalcRelativeHumidity(data.temp, data.dewpoint), ShouldEqual, data.rh)
		}
	})
}
