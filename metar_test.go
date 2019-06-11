package metar

import (
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/urkk/metar/clouds"
	. "github.com/urkk/metar/phenomena"
	"github.com/urkk/metar/runways"
	"github.com/urkk/metar/wind"
)

var curYear = time.Now().Year()
var curMonth = time.Now().Month()
var curDay = time.Now().Day()

func TestParseVisibility(t *testing.T) {

	type testpair struct {
		input    []string
		expected *Visibility
	}

	var onetokenpair = []testpair{
		{[]string{"2000"}, &Visibility{2000, 0, ""}},
		{[]string{"9999"}, &Visibility{9999, 0, ""}},
		{[]string{"5500"}, &Visibility{5500, 0, ""}},
	}

	var twotokenpair = []testpair{
		{[]string{"3000", "1500NE"}, &Visibility{3000, 1500, "NE"}},
		{[]string{"1500", "1000S"}, &Visibility{1500, 1000, "S"}},
		{[]string{"7000", "5000W"}, &Visibility{7000, 5000, "W"}},
	}

	var incorrecttokenpair = []testpair{
		{[]string{"20008MPS"}, &Visibility{0, 0, ""}},
		{[]string{"OVC020"}, &Visibility{0, 0, ""}},
	}

	Convey("Prevailing visibility parsing tests", t, func() {
		vis := &Visibility{}
		Convey("One token testing", func() {
			var t int
			for _, pair := range onetokenpair {
				t = vis.ParseVisibility(pair.input)
				So(vis, ShouldResemble, pair.expected)
				So(t, ShouldEqual, 1)
			}
		})

		Convey("Two token testing", func() {
			var t int
			for _, pair := range twotokenpair {
				t = vis.ParseVisibility(pair.input)
				So(vis, ShouldResemble, pair.expected)
				So(t, ShouldEqual, 2)
			}
		})

		Convey("Incorrect token testing", func() {
			var t int
			for _, pair := range incorrecttokenpair {
				t = vis.ParseVisibility(pair.input)
				So(vis, ShouldResemble, pair.expected)
				So(t, ShouldEqual, 0)
			}
		})

	})

}

func getWind(inp string) wind.Wind {
	w := &wind.Wind{}
	w.ParseWind(inp)
	return *w
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
			Phenomena:                    []Phenomenon{Phenomenon{Vicinity: true, Intensity: "", Abbreviation: "FG"}},
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
			Phenomena:   []Phenomenon{Phenomenon{Vicinity: true, Intensity: "", Abbreviation: "FG"}},
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
				Phenomena:          []Phenomenon{Phenomenon{Vicinity: false, Abbreviation: "FG", Intensity: ""}},
				VerticalVisibility: 8000,
			},
			},
			Remarks:          &Remark{QFE: 748},
			NotDecodedTokens: []string{"END"},
		}},
	{"METAR URSS 270600Z 22003MPS 180V250 5000 4000NW R24/P2000 R30/6000 // VV100 M17/M11 A3006 RESN WS ALL RWY",
		&MetarMessage{rawData: "METAR URSS 270600Z 22003MPS 180V250 5000 4000NW R24/P2000 R30/6000 // VV100 M17/M11 A3006 RESN WS ALL RWY",
			Station:    "URSS",
			DateTime:   time.Date(curYear, curMonth, 27, 6, 0, 0, 0, time.UTC),
			Wind:       getWind("22003MPS 180V250"),
			Visibility: Visibility{Distance: 5000, LowerDistance: 4000, LowerDirection: "NW"},
			RWYvisibility: []runways.VisualRange{runways.VisualRange{Designator: runways.RunwayDesignator{Number: "24", AllRunways: false},
				Distance: 2000,
				AboveMax: true,
				Trend:    ""},
				runways.VisualRange{Designator: runways.RunwayDesignator{Number: "30", AllRunways: false},
					Distance: 6000,
					AboveMax: false,
					Trend:    ""},
			},
			PhenomenaNotDefined: true,
			VerticalVisibility:  10000,
			Temperature:         -17,
			Dewpoint:            -11,
			QNHhPa:              1018,
			WindShear:           []runways.RunwayDesignator{runways.RunwayDesignator{Number: "", AllRunways: true}},
			RecentPhenomena:     []Phenomenon{Phenomenon{Vicinity: false, Abbreviation: "SN", Intensity: ""}},
		}},
	{"TAF UUWW 121350Z 1215/1315 VRB01MPS 9999 SCT040 TX22/1215Z TN13/1302Z TEMPO 1221/1315 3100 -SHRA FEW007 BKN011CB",
		&MetarMessage{rawData: "TAF UUWW 121350Z 1215/1315 VRB01MPS 9999 SCT040 TX22/1215Z TN13/1302Z TEMPO 1221/1315 3100 -SHRA FEW007 BKN011CB"}},
	{"COR UHMM 240700Z 11006MPS 9999 -SHRA SCT011 OVC018CB 05/04 Q1002 RMK R01/18004MPS",
		&MetarMessage{rawData: "COR UHMM 240700Z 11006MPS 9999 -SHRA SCT011 OVC018CB 05/04 Q1002 RMK R01/18004MPS",
			COR:         true,
			Station:     "UHMM",
			DateTime:    time.Date(curYear, curMonth, 24, 7, 0, 0, 0, time.UTC),
			Wind:        getWind("11006MPS"),
			Visibility:  Visibility{Distance: 9999, LowerDistance: 0, LowerDirection: ""},
			Clouds:      clouds.Clouds{getCloud("SCT011"), getCloud("OVC018CB")},
			Temperature: 5,
			Dewpoint:    4,
			QNHhPa:      1002,
			Phenomena:   []Phenomenon{Phenomenon{Vicinity: false, Abbreviation: "SHRA", Intensity: "-"}},
			Remarks:     &Remark{WindOnRWY: []WindOnRWY{WindOnRWY{Runway: "01", Wind: getWind("18004MPS")}}},
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
		msg, err := NewMETAR(pair.input)
		if err == nil && !reflect.DeepEqual(msg, pair.expected) || msg.RAW() != pair.input {
			t.Error(
				"For", pair.input,
				"expected", pair.expected,
				"got", msg,
			)
		}
	}
}
