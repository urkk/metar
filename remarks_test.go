package metar

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseRemarks(t *testing.T) {
	type testpair struct {
		input    []string
		expected *Remark
	}

	var tests = []testpair{
		// WindOnRWY []WindOnRWY
		// QBB       int  // cloud base in meters
		// МТOBSC    bool // Mountains obscured
		// MASTOBSC  bool // Mast obscured
		// OBSTOBSC  bool // Obstacle obscured
		// QFE       int  // Q-code Field Elevation (mmHg/hPa)

		{[]string{"RMK", "R06/25002MPS", "QFE762"}, &Remark{WindOnRWY: []WindOnRWY{WindOnRWY{Runway: "06", Wind: getWind("25002MPS")}}, QBB: 0, МТOBSC: false, MASTOBSC: false, OBSTOBSC: false, QFE: 762}},
		{[]string{"QBB200", "MT", "OBSC", "QFE762"}, &Remark{WindOnRWY: nil, QBB: 200, МТOBSC: true, MASTOBSC: false, OBSTOBSC: false, QFE: 762}},
		{[]string{"QBB180"}, &Remark{WindOnRWY: nil, QBB: 180, МТOBSC: false, MASTOBSC: false, OBSTOBSC: false, QFE: 0}},
		{[]string{"MT", "OBSC", "MAST", "OBSC", "OBST", "OBSC", "QFE762/1004"}, &Remark{WindOnRWY: nil, QBB: 0, МТOBSC: true, MASTOBSC: true, OBSTOBSC: true, QFE: 762}},
		{[]string{"MT", "OBSC"}, &Remark{WindOnRWY: nil, QBB: 0, МТOBSC: true, MASTOBSC: false, OBSTOBSC: false, QFE: 0}},
		{[]string{"MAST", "OBSC"}, &Remark{WindOnRWY: nil, QBB: 0, МТOBSC: false, MASTOBSC: true, OBSTOBSC: false, QFE: 0}},
		{[]string{"OBST", "OBSC"}, &Remark{WindOnRWY: nil, QBB: 0, МТOBSC: false, MASTOBSC: false, OBSTOBSC: true, QFE: 0}},
		{[]string{"RMK", "R06/25002MPS", "120V180"}, &Remark{WindOnRWY: []WindOnRWY{WindOnRWY{Runway: "06", Wind: getWind("25002MPS 120V180")}}, QBB: 0, МТOBSC: false, MASTOBSC: false, OBSTOBSC: false, QFE: 0}},
	}

	Convey("Remarks parsing tests", t, func() {
		for _, pair := range tests {
			So(parseRemarks(pair.input), ShouldResemble, pair.expected)
		}

	})
}
