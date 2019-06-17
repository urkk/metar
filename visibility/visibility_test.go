package visibility

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/urkk/metar/conversion"
)

func TestParseVisibility(t *testing.T) {

	type testpair struct {
		input    []string
		expected *Visibility
	}

	var onetokenpair = []testpair{
		{[]string{"2000"}, &Visibility{BaseVisibility{Distance{2000, 0.0, M}, false, false}, Distance{}, ""}},
		{[]string{"9999"}, &Visibility{BaseVisibility{Distance{9999, 0.0, M}, false, false}, Distance{}, ""}},
		{[]string{"5500"}, &Visibility{BaseVisibility{Distance{5500, 0.0, M}, false, false}, Distance{}, ""}},
		{[]string{"8SM"}, &Visibility{BaseVisibility{Distance{8, 0.0, SM}, false, false}, Distance{}, ""}},
		{[]string{"15SM"}, &Visibility{BaseVisibility{Distance{15, 0.0, SM}, false, false}, Distance{}, ""}},
		{[]string{"3/4SM"}, &Visibility{BaseVisibility{Distance{0, 0.75, SM}, false, false}, Distance{}, ""}},
		{[]string{"1/8SM"}, &Visibility{BaseVisibility{Distance{0, 0.125, SM}, false, false}, Distance{}, ""}},
		// like 2 and 1/1 with a missing space
		{[]string{"21/2SM"}, &Visibility{BaseVisibility{Distance{2, 0.5, SM}, false, false}, Distance{}, ""}},
		{[]string{"P6SM"}, &Visibility{BaseVisibility{Distance{6, 0.0, SM}, true, false}, Distance{}, ""}},
		{[]string{"P5000"}, &Visibility{BaseVisibility{Distance{5000, 0.0, M}, true, false}, Distance{}, ""}},
		{[]string{"M0050"}, &Visibility{BaseVisibility{Distance{50, 0.0, M}, false, true}, Distance{}, ""}},
		{[]string{"M1/4SM"}, &Visibility{BaseVisibility{Distance{0, 0.25, SM}, false, true}, Distance{}, ""}},
	}

	var twotokenpair = []testpair{
		{[]string{"3000", "1500NE"}, &Visibility{BaseVisibility{Distance{3000, 0.0, M}, false, false}, Distance{1500, .0, M}, "NE"}},
		{[]string{"5500", "1000S"}, &Visibility{BaseVisibility{Distance{5500, 0.0, M}, false, false}, Distance{1000, .0, M}, "S"}},
		{[]string{"7000", "5000W"}, &Visibility{BaseVisibility{Distance{7000, 0.0, M}, false, false}, Distance{5000, .0, M}, "W"}},
		{[]string{"5", "1/4SM"}, &Visibility{BaseVisibility{Distance{5, 0.25, SM}, false, false}, Distance{}, ""}},
		{[]string{"2", "1/2SM"}, &Visibility{BaseVisibility{Distance{2, 0.5, SM}, false, false}, Distance{}, ""}},
		{[]string{"1", "3/4SM"}, &Visibility{BaseVisibility{Distance{1, 0.75, SM}, false, false}, Distance{}, ""}},
	}

	var incorrecttokenpair = []testpair{
		{[]string{"20008MPS"}, &Visibility{BaseVisibility{Distance{}, false, false}, Distance{}, ""}},
		{[]string{"OVC020"}, &Visibility{BaseVisibility{Distance{}, false, false}, Distance{}, ""}},
	}

	Convey("Prevailing visibility parsing tests", t, func() {
		Convey("One token testing", func() {
			var t int
			for _, pair := range onetokenpair {
				vis := &Visibility{}
				t = vis.ParseVisibility(pair.input)
				So(vis, ShouldResemble, pair.expected)
				So(t, ShouldEqual, 1)
			}
		})

		Convey("Two token testing", func() {
			var t int
			for _, pair := range twotokenpair {
				vis := &Visibility{}
				t = vis.ParseVisibility(pair.input)
				So(vis, ShouldResemble, pair.expected)
				So(t, ShouldEqual, 2)
			}
		})

		Convey("Incorrect token testing", func() {
			var t int
			for _, pair := range incorrecttokenpair {
				vis := &Visibility{}
				t = vis.ParseVisibility(pair.input)
				So(vis, ShouldResemble, pair.expected)
				So(t, ShouldEqual, 0)
			}
		})

		Convey("Unit conversion testing", func() {
			for length := 100; length < 8000; {
				vis := &Visibility{}
				vis.Distance.Value = length
				So(vis.Distance.Meters(), ShouldEqual, length)
				So(vis.Distance.Feet(), ShouldEqual, MToFt(length))
				So(vis.Distance.Miles(), ShouldEqual, MToSMile(length))
				vis.Distance.Unit = M
				So(vis.Distance.Meters(), ShouldEqual, length)
				So(vis.Distance.Feet(), ShouldEqual, MToFt(length))
				So(vis.Distance.Miles(), ShouldEqual, MToSMile(length))
				vis.Distance.Unit = FT
				So(vis.Distance.Meters(), ShouldEqual, FtToM(length))
				So(vis.Distance.Feet(), ShouldEqual, length)
				So(vis.Distance.Miles(), ShouldEqual, FtToSMile(length))
				vis.Distance.Unit = SM
				So(vis.Distance.Meters(), ShouldEqual, SMileToM(float64(length)))
				So(vis.Distance.Feet(), ShouldEqual, SMileToFt(float64(length)))
				So(vis.Distance.Miles(), ShouldEqual, length)
				length += 500
			}
		})

	})

}
