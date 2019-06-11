package phenomena

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParsePhenomena(t *testing.T) {
	arr := &Phenomena{}
	type testpair struct {
		input    string
		expected *Phenomenon
	}

	var tests = []testpair{
		// Vicinity     bool
		// Intensity    Intensity
		// Abbreviation string
		{"PLRADZ", &Phenomenon{false, Moderate, "PLRADZ"}},
		{"+SHRASNGS", &Phenomenon{false, Heavy, "SHRASNGS"}},
		{"VCBLDU", &Phenomenon{true, Moderate, "BLDU"}},
		// not can VC
		{"VCSNDZ", nil},
		{"SHGRRA", &Phenomenon{false, Moderate, "SHGRRA"}},
		// "Light" not be applicable
		{"-FC", nil},
		{"-PLDZ", &Phenomenon{false, Light, "PLDZ"}},
	}

	Convey("Phenomena parsing tests", t, func() {
		Convey("Phenomena must parsed correctly", func() {
			for _, pair := range tests {
				So(ParsePhenomena(pair.input), ShouldResemble, pair.expected)
			}
		})

		Convey("Correct phenomena must can be appended", func() {
			for _, pair := range tests {
				So(arr.AppendPhenomena(pair.input), ShouldResemble, ParsePhenomena(pair.input) != nil)
			}
		})
	})

}

func TestParseRecentPhenomena(t *testing.T) {
	arr := &Phenomena{}

	type testpair struct {
		input    string
		expected *Phenomenon
	}
	var recenttests = []testpair{
		{"REFZDZ", &Phenomenon{false, Moderate, "FZDZ"}},
		{"+REFZDZ", nil}, // + not applicable in recent weather
		{"RERASN", &Phenomenon{false, Moderate, "RASN"}},
	}

	Convey("Recent phenomena parsing tests", t, func() {
		Convey("Recent phenomena must parsed correctly", func() {
			for _, pair := range recenttests {
				So(ParseRecentPhenomena(pair.input), ShouldResemble, pair.expected)
			}
		})

		Convey("Correct recent phenomena must can be appended", func() {
			for _, pair := range recenttests {
				So(arr.AppendRecentPhenomena(pair.input), ShouldResemble, ParseRecentPhenomena(pair.input) != nil)
			}
		})

	})
}
