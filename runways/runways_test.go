package runways

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseVisibility(t *testing.T) {
	type testpair struct {
		input    string
		expected VisualRange
	}

	var tests = []testpair{
		{"R25/M0075", VisualRange{RunwayDesignator{"25", false}, 75, false, true, NotDefined}},
		{"R33L/P1500", VisualRange{RunwayDesignator{"33L", false}, 1500, true, false, NotDefined}},
		{"R16R/1000U", VisualRange{RunwayDesignator{"16R", false}, 1000, false, false, U}},
		//artificial situation
		{"R88/1000D", VisualRange{RunwayDesignator{"88", true}, 1000, false, false, D}},
		{"R33C/0900N", VisualRange{RunwayDesignator{"33C", false}, 900, false, false, N}},
		{"OVC350", VisualRange{RunwayDesignator{"", false}, 0, false, false, NotDefined}},
	}

	Convey("Runway visual range parsing tests", t, func() {
		for _, pair := range tests {
			vis, _ := ParseVisibility(pair.input)
			So(vis, ShouldResemble, pair.expected)
		}
	})
}

func TestParseState(t *testing.T) {

	type testpair struct {
		input    string
		expected State
	}

	var tests = []testpair{

		{"R25/CLRD70", State{Designator: RunwayDesignator{"25", false},
			BrakingConditions: 70,
			CLRD:              true}},
		{"R24L/451293", State{Designator: RunwayDesignator{"24L", false},
			TypeOfCoverage:      4,
			DimensionOfCoverage: 5,
			HeightOfCoverage:    12,
			BrakingConditions:   93,
		}},
		{"R30/290250", State{Designator: RunwayDesignator{"30", false},
			TypeOfCoverage:      2,
			DimensionOfCoverage: 9,
			HeightOfCoverage:    2,
			BrakingConditions:   50,
		}},
		{"R21/0///65", State{Designator: RunwayDesignator{"21", false},
			TypeOfCoverage:            0,
			DimensionOfCoverageNotDef: true,
			HeightOfCoverageNotDef:    true,
			BrakingConditions:         65,
		}},
		{"R88///////", State{Designator: RunwayDesignator{"88", true},
			TypeOfCoverageNotDef:        true,
			DimensionOfCoverageNotDef:   true,
			HeightOfCoverageNotDef:      true,
			BrakingConditionsNotDefined: true,
		}},
		{"R74/4/1293", State{Designator: RunwayDesignator{"24R", false},
			TypeOfCoverage:            4,
			DimensionOfCoverageNotDef: true,
			HeightOfCoverage:          12,
			BrakingConditions:         93,
		}},
		{"R31/70D", State{Designator: RunwayDesignator{"31", false},
			BrakingConditions: 70,
			CLRD:              true,
		}},
		{"OVC350", State{Designator: RunwayDesignator{"", false}}},
	}

	Convey("Runway state parsing tests", t, func() {
		for _, pair := range tests {
			st, _ := ParseState(pair.input)
			So(st, ShouldResemble, pair.expected)
		}
	})
}
