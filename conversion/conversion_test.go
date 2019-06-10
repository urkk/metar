package conversion

import (
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHPaToMmHg(t *testing.T) {
	Convey("converts hectopascal to mm of mercury", t, func() {
		for input := 985; input < 1027; input++ {
			So(HPaToMmHg(input), ShouldAlmostEqual, float64(input)*0.75006375541921, 1)
		}
	})
}

func TestMmHgToHPa(t *testing.T) {
	Convey("converts mm of mercury to hectopascal", t, func() {
		for input := 740; input < 780; input++ {
			So(MmHgToHPa(input), ShouldAlmostEqual, float64(input)*1.333223684, 1)
		}
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

func TestKphToMps(t *testing.T) {
	Convey("converts kilometres per hour to meters per second", t, func() {
		for input := .0; input < 200; input++ {
			So(KphToMps(input), ShouldAlmostEqual, input/3.6, .00001)
		}
	})
}

func TestKtsToMps(t *testing.T) {
	Convey("converts knots to meters per second", t, func() {
		for input := .0; input < 50; input++ {
			So(KtsToMps(input), ShouldAlmostEqual, input/1.94384, .00001)
		}
	})
}

func TestMpsToKts(t *testing.T) {
	Convey("converts meters per second to knots", t, func() {
		for input := .0; input < 50; input++ {
			So(MpsToKts(input), ShouldAlmostEqual, input*1.94384, .00001)
		}
	})
}

func TestSMileToM(t *testing.T) {
	Convey("converts statute miles to meters", t, func() {
		for input := .0; input < 6.25; {
			So(SMileToM(input), ShouldAlmostEqual, input*1609.344, 1)
			input += .25
		}
	})
}

func TestFtToM(t *testing.T) {
	Convey("converts feet to meters (rounded to 10 meters)", t, func() {
		for input := 300; input < 3000; {
			So(FtToM(input), ShouldAlmostEqual, int(math.Round(float64(input)*0.3048/10)*10), 10)
			input += 300
		}
	})
}

func TestInHgTohPa(t *testing.T) {
	Convey("converts inch of mercury to hectopascal", t, func() {
		for input := 29.0; input < 30.3; {
			So(InHgTohPa(input), ShouldAlmostEqual, input*33.86389, 1)
			input += .1
		}
	})
}
