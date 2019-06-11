package clouds

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/urkk/metar/conversion"
)

func TestParseCloud(t *testing.T) {
	arr := &Clouds{}

	type testpair struct {
		input    string
		expected Cloud
	}

	var tests = []testpair{

		// Type CloudType
		// height           int
		// HeightNotDefined bool
		// Cumulonimbus     bool
		// ToweringCumulus  bool
		// CBNotDefined     bool

		{"FEW005", Cloud{FEW, 5, false, false, false, false}},
		{"FEW010CB", Cloud{FEW, 10, false, true, false, false}},
		{"SCT018", Cloud{SCT, 18, false, false, false, false}},
		{"BKN025///", Cloud{BKN, 25, false, false, false, true}},
		{"OVC///", Cloud{OVC, 0, true, false, false, false}},
		{"///015", Cloud{NotDefined, 15, false, false, false, false}},
		{"//////", Cloud{NotDefined, 0, true, false, false, false}},
		{"//////CB", Cloud{NotDefined, 0, true, true, false, false}},
		{"BKN020TCU", Cloud{BKN, 20, false, false, true, false}},
		{"NSC", Cloud{NSC, 0, false, false, false, false}},
		{"RESHRA", Cloud{"", 0, false, false, false, false}},
	}
	Convey("Cloud layer parsing tests", t, func() {
		Convey("cloud must parsed correctly", func() {
			for _, pair := range tests {
				cloud, _ := ParseCloud(pair.input)
				So(cloud, ShouldResemble, pair.expected)

			}
		})

		Convey("height in feet must calculated correctly", func() {
			for _, pair := range tests {
				cloud, _ := ParseCloud(pair.input)
				So(cloud.HeightFt(), ShouldEqual, pair.expected.height*100)
			}
		})

		Convey("height in meters must calculated correctly", func() {
			for _, pair := range tests {
				cloud, _ := ParseCloud(pair.input)
				So(cloud.HeightM(), ShouldEqual, FtToM(pair.expected.height*100))
			}
		})

		Convey("correct cloud must can be appended", func() {
			for _, pair := range tests {
				_, ok := ParseCloud(pair.input)
				So(arr.AppendCloud(pair.input), ShouldEqual, ok)
			}
		})

	})
}
