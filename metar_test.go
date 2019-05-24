package metar

import (
	"reflect"
	"testing"
	"time"

	"github.com/urkk/metar/clouds"
	"github.com/urkk/metar/phenomena"
	"github.com/urkk/metar/runways"
	"github.com/urkk/metar/wind"
)

var curYear = time.Now().Year()
var curMonth = time.Now().Month()
var curDay = time.Now().Day()

type visibilityparsetest struct {
	input    string
	expected Visibility
	multiple bool
}

var visibilityparsetests = []visibilityparsetest{
	// Distance       int
	// LowerDistance  int
	// LowerDirection string

	{"2000", Visibility{2000, 0, ""}, false},
	{"3000 1500NE", Visibility{3000, 1500, "NE"}, true},
	{"1500 1000S", Visibility{1500, 1000, "S"}, true},
	{"9999", Visibility{9999, 0, ""}, false},
	{"20008MPS", Visibility{0, 0, ""}, false},
}

func TestParseVisibility(t *testing.T) {
	for _, pair := range visibilityparsetests {
		v, ok, multiple := ParseVisibility(pair.input)
		if ok && v != pair.expected || multiple != pair.multiple {
			t.Error(
				"For", pair.input,
				"expected", pair.expected,
				"got", v,
			)
		}
	}
}

type remarksparsetest struct {
	input    []string
	expected Remark
}

func getWind(inp string) wind.Wind {
	w, _, _ := wind.ParseWind(inp)
	return w
}

var remarksparsetests = []remarksparsetest{
	// WindOnRWY []WindOnRWY
	// QBB       int  // cloud base in meters
	// МТOBSC    bool // Mountains obscured
	// MASTOBSC  bool // Mast obscured
	// OBSTOBSC  bool // Obstacle obscured
	// QFE       int  // Q-code Field Elevation (mmHg/hPa)

	{[]string{"RMK", "R06/25002MPS", "QFE762"}, Remark{WindOnRWY: []WindOnRWY{WindOnRWY{Runway: "06", Wind: getWind("25002MPS")}}, QBB: 0, МТOBSC: false, MASTOBSC: false, OBSTOBSC: false, QFE: 762}},
	{[]string{"QBB200", "MT", "OBSC", "QFE762"}, Remark{WindOnRWY: nil, QBB: 200, МТOBSC: true, MASTOBSC: false, OBSTOBSC: false, QFE: 762}},
	{[]string{"QBB180"}, Remark{WindOnRWY: nil, QBB: 180, МТOBSC: false, MASTOBSC: false, OBSTOBSC: false, QFE: 0}},
	{[]string{"MT", "OBSC", "MAST", "OBSC", "OBST", "OBSC", "QFE762/1004"}, Remark{WindOnRWY: nil, QBB: 0, МТOBSC: true, MASTOBSC: true, OBSTOBSC: true, QFE: 762}},
	{[]string{"MT", "OBSC"}, Remark{WindOnRWY: nil, QBB: 0, МТOBSC: true, MASTOBSC: false, OBSTOBSC: false, QFE: 0}},
	{[]string{"MAST", "OBSC"}, Remark{WindOnRWY: nil, QBB: 0, МТOBSC: false, MASTOBSC: true, OBSTOBSC: false, QFE: 0}},
	{[]string{"OBST", "OBSC"}, Remark{WindOnRWY: nil, QBB: 0, МТOBSC: false, MASTOBSC: false, OBSTOBSC: true, QFE: 0}},
	{[]string{"RMK", "R06/25002MPS", "120V180"}, Remark{WindOnRWY: []WindOnRWY{WindOnRWY{Runway: "06", Wind: getWind("25002MPS 120V180")}}, QBB: 0, МТOBSC: false, MASTOBSC: false, OBSTOBSC: false, QFE: 0}},
}

func TestParseRemarks(t *testing.T) {
	for _, pair := range remarksparsetests {
		r := *parseRemarks(pair.input)
		if !reflect.DeepEqual(r, pair.expected) {
			t.Error(
				"For", pair.input,
				"expected", pair.expected,
				"got", r,
			)
		}
	}
}

type metarparsetest struct {
	input    string
	expected *MetarMessage
}

var metarparsetests = []metarparsetest{

	{"METAR URSS 270600Z 22003MPS 9999 VCFG VV/// 17/11 Q1018 R02/010070 NOSIG",
		&MetarMessage{rawData: "METAR URSS 270600Z 22003MPS 9999 VCFG VV/// 17/11 Q1018 R02/010070 NOSIG",
			Station:                      "URSS",
			DateTime:                     time.Date(curYear, curMonth, 27, 6, 0, 0, 0, time.UTC),
			Wind:                         getWind("22003MPS"),
			Visibility:                   Visibility{Distance: 9999, LowerDistance: 0, LowerDirection: ""},
			VerticalVisibilityNotDefined: true,
			Phenomena:                    []phenomena.Phenomena{phenomena.Phenomena{Vicinity: true, Intensity: "", Abbreviation: "FG"}},
			Temperature:                  17,
			Dewpoint:                     11,
			QNHhPa:                       1018,
			RWYState:                     []runways.State{runways.State{Designator: runways.RunwayDesignator{Number: "02", AllRunways: false}, TypeOfCoverage: 0, TypeOfCoverageNotDef: false, DimensionOfCoverage: 1, DimensionOfCoverageNotDef: false, HeightOfCoverage: 0, HeightOfCoverageNotDef: false, BrakingConditions: 70, BrakingConditionsNotDefined: false, CLRD: false, SNOCLO: false}},
			NOSIG:                        true}},
	{"METAR URSS 270600Z 22003MPS 9999 VCFG NSC 17/11 Q1018 R02/010070 NOSIG",
		&MetarMessage{rawData: "METAR URSS 270600Z 22003MPS 9999 VCFG NSC 17/11 Q1018 R02/010070 NOSIG",
			Station:     "URSS",
			DateTime:    time.Date(curYear, curMonth, 27, 6, 0, 0, 0, time.UTC),
			Wind:        getWind("22003MPS"),
			Visibility:  Visibility{Distance: 9999, LowerDistance: 0, LowerDirection: ""},
			Phenomena:   []phenomena.Phenomena{phenomena.Phenomena{Vicinity: true, Intensity: "", Abbreviation: "FG"}},
			Clouds:      []clouds.Cloud{getCloud("NSC")},
			Temperature: 17,
			Dewpoint:    11,
			QNHhPa:      1018,
			RWYState:    []runways.State{runways.State{Designator: runways.RunwayDesignator{Number: "02", AllRunways: false}, TypeOfCoverage: 0, TypeOfCoverageNotDef: false, DimensionOfCoverage: 1, DimensionOfCoverageNotDef: false, HeightOfCoverage: 0, HeightOfCoverageNotDef: false, BrakingConditions: 70, BrakingConditionsNotDefined: false, CLRD: false, SNOCLO: false}},
			NOSIG:       true}},
	{"METAR URSS 270600Z 22003MPS CAVOK M17/M11 Q1018 WS R06R R02/010070 NOSIG",
		&MetarMessage{rawData: "METAR URSS 270600Z 22003MPS CAVOK M17/M11 Q1018 WS R06R R02/010070 NOSIG",
			Station:     "URSS",
			DateTime:    time.Date(curYear, curMonth, 27, 6, 0, 0, 0, time.UTC),
			Wind:        getWind("22003MPS"),
			Visibility:  Visibility{Distance: 0, LowerDistance: 0, LowerDirection: ""},
			CAVOK:       true,
			Temperature: -17,
			Dewpoint:    -11,
			QNHhPa:      1018,
			RWYState:    []runways.State{runways.State{Designator: runways.RunwayDesignator{Number: "02", AllRunways: false}, TypeOfCoverage: 0, TypeOfCoverageNotDef: false, DimensionOfCoverage: 1, DimensionOfCoverageNotDef: false, HeightOfCoverage: 0, HeightOfCoverageNotDef: false, BrakingConditions: 70, BrakingConditionsNotDefined: false, CLRD: false, SNOCLO: false}},
			WindShear:   []runways.RunwayDesignator{runways.RunwayDesignator{Number: "R06R", AllRunways: false}},
			NOSIG:       true}},
	{"METAR URSS 270600Z 22003MPS CAVOK END M17/M11 Q1018 WS R06R R/SNOCLO TEMPO 32010G17MPS 8000 FG VV080 RMK QFE748",
		&MetarMessage{rawData: "METAR URSS 270600Z 22003MPS CAVOK END M17/M11 Q1018 WS R06R R/SNOCLO TEMPO 32010G17MPS 8000 FG VV080 RMK QFE748",
			Station:     "URSS",
			DateTime:    time.Date(curYear, curMonth, 27, 6, 0, 0, 0, time.UTC),
			Wind:        getWind("22003MPS"),
			Visibility:  Visibility{Distance: 0, LowerDistance: 0, LowerDirection: ""},
			CAVOK:       true,
			Temperature: -17,
			Dewpoint:    -11,
			QNHhPa:      1018,
			RWYState:    []runways.State{runways.State{Designator: runways.RunwayDesignator{Number: "", AllRunways: false}, TypeOfCoverage: 0, TypeOfCoverageNotDef: false, DimensionOfCoverage: 0, DimensionOfCoverageNotDef: false, HeightOfCoverage: 0, HeightOfCoverageNotDef: false, BrakingConditions: 0, BrakingConditionsNotDefined: false, CLRD: false, SNOCLO: true}},
			WindShear:   []runways.RunwayDesignator{runways.RunwayDesignator{Number: "R06R", AllRunways: false}},
			TREND: []Trend{Trend{
				Type:               TEMPO,
				Wind:               getWind("32010G17MPS"),
				Visibility:         Visibility{Distance: 8000, LowerDistance: 0, LowerDirection: ""},
				Phenomena:          []phenomena.Phenomena{phenomena.Phenomena{Vicinity: false, Abbreviation: "FG", Intensity: ""}},
				VerticalVisibility: 8000,
			},
			},
			Remarks:          &Remark{QFE: 748},
			NotDecodedTokens: []string{"END"},
		}},
	{"METAR URSS 270600Z 22003MPS 180V250 5000 4000NW R24/P2000 // VV100 M17/M11 A3006 RESN WS ALL RWY",
		&MetarMessage{rawData: "METAR URSS 270600Z 22003MPS 180V250 5000 4000NW R24/P2000 // VV100 M17/M11 A3006 RESN WS ALL RWY",
			Station:    "URSS",
			DateTime:   time.Date(curYear, curMonth, 27, 6, 0, 0, 0, time.UTC),
			Wind:       getWind("22003MPS 180V250"),
			Visibility: Visibility{Distance: 5000, LowerDistance: 4000, LowerDirection: "NW"},
			RWYvisibility: []runways.VisualRange{runways.VisualRange{Designator: runways.RunwayDesignator{Number: "24", AllRunways: false},
				Distance: 2000,
				AboveMax: true,
				Trend:    ""}},
			PhenomenaNotDefined: true,
			VerticalVisibility:  10000,
			Temperature:         -17,
			Dewpoint:            -11,
			QNHhPa:              1018,
			WindShear:           []runways.RunwayDesignator{runways.RunwayDesignator{Number: "", AllRunways: true}},
			RecentPhenomena:     []phenomena.Phenomena{phenomena.Phenomena{Vicinity: false, Abbreviation: "SN", Intensity: ""}},
		}},
	{"TAF UUWW 121350Z 1215/1315 VRB01MPS 9999 SCT040 TX22/1215Z TN13/1302Z TEMPO 1221/1315 3100 -SHRA FEW007 BKN011CB",
		&MetarMessage{rawData: "TAF UUWW 121350Z 1215/1315 VRB01MPS 9999 SCT040 TX22/1215Z TN13/1302Z TEMPO 1221/1315 3100 -SHRA FEW007 BKN011CB",
			NotDecodedTokens: []string{"TAF UUWW 121350Z 1215/1315 VRB01MPS 9999 SCT040 TX22/1215Z TN13/1302Z TEMPO 1221/1315 3100 -SHRA FEW007 BKN011CB"},
		}},
	{"COR UHMM 240700Z 11006MPS 9999 -SHRA SCT011 OVC018CB 05/04 Q1002",
		&MetarMessage{rawData: "COR UHMM 240700Z 11006MPS 9999 -SHRA SCT011 OVC018CB 05/04 Q1002",
			COR:         true,
			Station:     "UHMM",
			DateTime:    time.Date(curYear, curMonth, 24, 7, 0, 0, 0, time.UTC),
			Wind:        getWind("11006MPS"),
			Visibility:  Visibility{Distance: 9999, LowerDistance: 0, LowerDirection: ""},
			Clouds:      []clouds.Cloud{getCloud("SCT011"), getCloud("OVC018CB")},
			Temperature: 5,
			Dewpoint:    4,
			QNHhPa:      1002,
			Phenomena:   []phenomena.Phenomena{phenomena.Phenomena{Vicinity: false, Abbreviation: "SHRA", Intensity: "-"}},
		}},
	{"METAR UOOO 052030Z NIL",
		&MetarMessage{rawData: "METAR UOOO 052030Z NIL",
			Station:  "UOOO",
			DateTime: time.Date(curYear, curMonth, 5, 20, 30, 0, 0, time.UTC),
			NIL:      true,
		}},
}

func TestDecode(t *testing.T) {
	for _, pair := range metarparsetests {
		msg := NewMETAR(pair.input)
		if !reflect.DeepEqual(msg, pair.expected) {
			t.Error(
				"For", pair.input,
				"expected", pair.expected,
				"got", msg,
			)
		}
	}
}
