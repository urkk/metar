package runways

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "github.com/urkk/metar/visibility"
)

func TestParseVisibility(t *testing.T) {
	type testpair struct {
		input    string
		expected VisualRange
	}

	var tests = []testpair{
		{"R25/M0075", VisualRange{Designator: RunwayDesignator{"25", false},
			Visibility:     BaseVisibility{Distance: Distance{Value: 75, FractionValue: 0.0, Unit: ""}, AboveMax: false, BelowMin: true},
			UpToVisibility: BaseVisibility{}, Trend: NotDefined}},
		{"R33L/P1500", VisualRange{Designator: RunwayDesignator{"33L", false},
			Visibility:     BaseVisibility{Distance: Distance{Value: 1500, FractionValue: 0.0, Unit: ""}, AboveMax: true, BelowMin: false},
			UpToVisibility: BaseVisibility{}, Trend: NotDefined}},
		{"R16R/1000U", VisualRange{Designator: RunwayDesignator{"16R", false},
			Visibility:     BaseVisibility{Distance: Distance{Value: 1000, FractionValue: 0.0, Unit: ""}, AboveMax: false, BelowMin: false},
			UpToVisibility: BaseVisibility{}, Trend: U}},
		//artificial situation
		{"R88/1000D", VisualRange{Designator: RunwayDesignator{"88", true},
			Visibility:     BaseVisibility{Distance: Distance{Value: 1000, FractionValue: 0.0, Unit: ""}, AboveMax: false, BelowMin: false},
			UpToVisibility: BaseVisibility{}, Trend: D}},
		{"R33C/0900N", VisualRange{Designator: RunwayDesignator{"33C", false},
			Visibility:     BaseVisibility{Distance: Distance{Value: 900, FractionValue: 0.0, Unit: ""}, AboveMax: false, BelowMin: false},
			UpToVisibility: BaseVisibility{}, Trend: N}},
		{"R06/P6000FT", VisualRange{Designator: RunwayDesignator{"06", false},
			Visibility:     BaseVisibility{Distance: Distance{Value: 6000, FractionValue: 0.0, Unit: "FT"}, AboveMax: true, BelowMin: false},
			UpToVisibility: BaseVisibility{}, Trend: NotDefined}},
		{"R32/2200V4000FT", VisualRange{Designator: RunwayDesignator{"32", false},
			Visibility:     BaseVisibility{Distance: Distance{Value: 2200, FractionValue: 0.0, Unit: "FT"}, AboveMax: false, BelowMin: false},
			UpToVisibility: BaseVisibility{Distance: Distance{Value: 4000, FractionValue: 0.0, Unit: "FT"}, AboveMax: false, BelowMin: false}, Trend: NotDefined}},
		{"R35/4500VP6000FT/D", VisualRange{Designator: RunwayDesignator{"35", false},
			Visibility:     BaseVisibility{Distance: Distance{Value: 4500, FractionValue: 0.0, Unit: "FT"}, AboveMax: false, BelowMin: false},
			UpToVisibility: BaseVisibility{Distance: Distance{Value: 6000, FractionValue: 0.0, Unit: "FT"}, AboveMax: true, BelowMin: false}, Trend: D}},

		{"OVC350", VisualRange{Designator: RunwayDesignator{"", false}, Visibility: BaseVisibility{}, UpToVisibility: BaseVisibility{}, Trend: NotDefined}},
	}

	Convey("Runway visual range parsing tests", t, func() {
		for _, pair := range tests {
			vis, _ := ParseVisualRange(pair.input)
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
