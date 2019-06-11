package wind

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/urkk/metar/conversion"
)

func TestParseWind(t *testing.T) {

	type testpair struct {
		input    string
		expected *Wind
		tokens   int
	}

	var windtests = []testpair{
		{"31005MPS", &Wind{310, 5, 0, false, 0, 0, false}, 1},
		{"31010KPH", &Wind{310, 2.7777777777777777, 0, false, 0, 0, false}, 1},
		{"VRB15MPS", &Wind{0, 15, 0, true, 0, 0, false}, 1},
		{"00000MPS", &Wind{0, 0, 0, false, 0, 0, false}, 1},
		{"240P49MPS", &Wind{240, 49, 0, false, 0, 0, true}, 1},
		{"04008G20MPS", &Wind{40, 8, 20, false, 0, 0, false}, 1},
		{"22003G08MPS 280V350", &Wind{220, 3, 8, false, 280, 350, false}, 2},
		{"14010KT", &Wind{140, 5.144456333854638, 0, false, 0, 0, false}, 1},
		{"BKN020", &Wind{0, 0, 0, false, 0, 0, false}, 0},
	}

	Convey("Wind parsing tests", t, func() {
		Convey("wind must parsed correctly", func() {
			for _, pair := range windtests {
				wnd := &Wind{}
				tokensused := wnd.ParseWind(pair.input)
				So(wnd, ShouldResemble, pair.expected)
				So(tokensused, ShouldEqual, pair.tokens)
			}
		})

		Convey("speed and gusts speed in meters per second must calculated correctly", func() {
			for _, pair := range windtests {
				wnd := &Wind{}
				wnd.ParseWind(pair.input)
				So(wnd.SpeedMps(), ShouldAlmostEqual, wnd.speed, .25)
				So(wnd.GustsSpeedMps(), ShouldAlmostEqual, wnd.gustsSpeed, .25)
			}
		})

		Convey("speed and gusts speed in knots must calculated correctly", func() {
			for _, pair := range windtests {
				wnd := &Wind{}
				wnd.ParseWind(pair.input)
				So(wnd.SpeedKt(), ShouldAlmostEqual, MpsToKts(wnd.speed), 1)
				So(wnd.GustsSpeedKt(), ShouldAlmostEqual, MpsToKts(wnd.gustsSpeed), 1)
			}
		})

	})
}
