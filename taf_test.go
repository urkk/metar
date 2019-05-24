package metar

import (
	"reflect"
	"testing"
	"time"

	"github.com/urkk/metar/clouds"
	"github.com/urkk/metar/phenomena"
)

type tafparsetest struct {
	input    string
	expected *TAFMessage
}

func getCloud(inp string) clouds.Cloud {
	cl, _ := clouds.ParseCloud(inp)
	return cl
}

var tafparsetests = []tafparsetest{

	{"TAF UUEE 170800Z 1709/1809 02003MPS 0700 FG BKN003 TX21/1712Z TNM07/1802Z",
		&TAFMessage{rawData: "TAF UUEE 170800Z 1709/1809 02003MPS 0700 FG BKN003 TX21/1712Z TNM07/1802Z",
			COR: false, AMD: false, NIL: false, Station: "UUEE",
			DateTime:         time.Date(curYear, curMonth, 17, 8, 0, 0, 0, time.UTC),
			ValidFrom:        time.Date(curYear, curMonth, 17, 9, 0, 0, 0, time.UTC),
			ValidTo:          time.Date(curYear, curMonth, 18, 9, 0, 0, 0, time.UTC),
			Wind:             getWind("02003MPS"),
			Visibility:       Visibility{Distance: 700, LowerDistance: 0, LowerDirection: ""},
			Phenomena:        []phenomena.Phenomena{phenomena.Phenomena{Vicinity: false, Abbreviation: "FG", Intensity: ""}},
			Clouds:           []clouds.Cloud{getCloud("BKN003")},
			Temperature:      []TemperatureForecast{TemperatureForecast{Temp: 21, DateTime: time.Date(curYear, curMonth, 17, 12, 0, 0, 0, time.UTC), IsMax: true, IsMin: false}, TemperatureForecast{Temp: -7, DateTime: time.Date(2019, 5, 18, 2, 0, 0, 0, time.UTC), IsMax: false, IsMin: true}},
			TREND:            nil,
			NotDecodedTokens: nil}},
	{"TAF UUEE 230500Z 2306/2324 02003MPS CAVOK END PROB40 TEMPO 1723/1810 18005MPS",
		&TAFMessage{rawData: "TAF UUEE 230500Z 2306/2324 02003MPS CAVOK END PROB40 TEMPO 1723/1810 18005MPS",
			COR: false, AMD: false, NIL: false, Station: "UUEE",
			DateTime:           time.Date(curYear, curMonth, 23, 5, 0, 0, 0, time.UTC),
			ValidFrom:          time.Date(curYear, curMonth, 23, 6, 0, 0, 0, time.UTC),
			ValidTo:            time.Date(curYear, curMonth, 24, 0, 0, 0, 0, time.UTC),
			CNL:                false,
			Wind:               getWind("02003MPS"),
			CAVOK:              true,
			Visibility:         Visibility{},
			Phenomena:          nil,
			VerticalVisibility: 0,
			Clouds:             nil,
			Temperature:        nil,
			TREND: []Trend{Trend{Type: TEMPO,
				Probability:                  40,
				Visibility:                   Visibility{},
				VerticalVisibility:           0,
				VerticalVisibilityNotDefined: false,
				Wind:                         getWind("18005MPS"),
				CAVOK:                        false,
				Phenomena:                    nil,
				Clouds:                       nil,
				FM:                           time.Date(curYear, curMonth, 17, 23, 0, 0, 0, time.UTC),
				TL:                           time.Date(curYear, curMonth, 18, 10, 0, 0, 0, time.UTC)}},
			NotDecodedTokens: []string{"END"}}},
	{"TAF UUEE 170800Z 1709/1809 02003MPS CAVOK PROB30 1723/1810 18005MPS",
		&TAFMessage{rawData: "TAF UUEE 170800Z 1709/1809 02003MPS CAVOK PROB30 1723/1810 18005MPS",
			COR: false, AMD: false, NIL: false, Station: "UUEE",
			DateTime:  time.Date(curYear, curMonth, 17, 8, 0, 0, 0, time.UTC),
			ValidFrom: time.Date(curYear, curMonth, 17, 9, 0, 0, 0, time.UTC),
			ValidTo:   time.Date(curYear, curMonth, 18, 9, 0, 0, 0, time.UTC),
			Wind:      getWind("02003MPS"),
			CAVOK:     true,
			TREND: []Trend{Trend{Type: TEMPO,
				Probability: 30,
				Wind:        getWind("18005MPS"),
				FM:          time.Date(curYear, 5, 17, 23, 0, 0, 0, time.UTC),
				TL:          time.Date(curYear, 5, 18, 10, 0, 0, 0, time.UTC)}},
			NotDecodedTokens: nil}},
	{"TAF UUEE 170800Z 1709/1809 02003MPS CAVOK PROB40 1724/1824 18005MPS",
		&TAFMessage{rawData: "TAF UUEE 170800Z 1709/1809 02003MPS CAVOK PROB40 1724/1824 18005MPS",
			COR: false, AMD: false, NIL: false, Station: "UUEE",
			DateTime:  time.Date(curYear, curMonth, 17, 8, 0, 0, 0, time.UTC),
			ValidFrom: time.Date(curYear, curMonth, 17, 9, 0, 0, 0, time.UTC),
			ValidTo:   time.Date(curYear, curMonth, 18, 9, 0, 0, 0, time.UTC),
			Wind:      getWind("02003MPS"),
			CAVOK:     true,
			TREND: []Trend{Trend{Type: TEMPO,
				Probability: 40,
				Wind:        getWind("18005MPS"),
				FM:          time.Date(curYear, curMonth, 18, 00, 0, 0, 0, time.UTC),
				TL:          time.Date(curYear, curMonth, 19, 00, 0, 0, 0, time.UTC)}},
			NotDecodedTokens: nil}},
	{"TAF UUEE 300800Z 0100/0124 02003MPS 3000 VV050 FM171230 18005MPS",
		&TAFMessage{rawData: "TAF UUEE 300800Z 0100/0124 02003MPS 3000 VV050 FM171230 18005MPS",
			COR: false, AMD: false, NIL: false, Station: "UUEE",
			DateTime:           time.Date(curYear, curMonth, 30, 8, 0, 0, 0, time.UTC),
			ValidFrom:          time.Date(curYear, curMonth+1, 1, 0, 0, 0, 0, time.UTC),
			ValidTo:            time.Date(curYear, curMonth+1, 2, 0, 0, 0, 0, time.UTC),
			Visibility:         Visibility{Distance: 3000, LowerDistance: 0, LowerDirection: ""},
			VerticalVisibility: 50,
			Wind:               getWind("02003MPS"),
			TREND: []Trend{Trend{Type: FM,
				Wind: getWind("18005MPS"),
				FM:   time.Date(curYear, curMonth, 17, 12, 30, 0, 0, time.UTC)}},
			NotDecodedTokens: nil}},
	{"TAF UERR 221100Z 2212/2312 27006G11MPS 9999 SCT030CB",
		&TAFMessage{rawData: "TAF UERR 221100Z 2212/2312 27006G11MPS 9999 SCT030CB",
			COR: false, AMD: false, NIL: false, Station: "UERR",
			DateTime:         time.Date(curYear, curMonth, 22, 11, 0, 0, 0, time.UTC),
			ValidFrom:        time.Date(curYear, curMonth, 22, 12, 0, 0, 0, time.UTC),
			ValidTo:          time.Date(curYear, curMonth, 23, 12, 0, 0, 0, time.UTC),
			Visibility:       Visibility{Distance: 9999, LowerDistance: 0, LowerDirection: ""},
			Wind:             getWind("27006G11MPS"),
			Clouds:           []clouds.Cloud{getCloud("SCT030CB")},
			NotDecodedTokens: nil}},
	{"TAF UUOL 221100Z 2212/2221 18003MPS CAVOK",
		&TAFMessage{rawData: "TAF UUOL 221100Z 2212/2221 18003MPS CAVOK",
			COR: false, AMD: false, NIL: false, Station: "UUOL",
			DateTime:         time.Date(curYear, curMonth, 22, 11, 0, 0, 0, time.UTC),
			ValidFrom:        time.Date(curYear, curMonth, 22, 12, 0, 0, 0, time.UTC),
			ValidTo:          time.Date(curYear, curMonth, 22, 21, 0, 0, 0, time.UTC),
			Visibility:       Visibility{Distance: 0, LowerDistance: 0, LowerDirection: ""},
			Wind:             getWind("18003MPS"),
			CAVOK:            true,
			NotDecodedTokens: nil}},
	{"TAF AMD UHMD 240602Z 2406/2412 CNL",
		&TAFMessage{rawData: "TAF AMD UHMD 240602Z 2406/2412 CNL",
			COR: false, AMD: true, NIL: false, Station: "UHMD", CNL: true,
			DateTime:  time.Date(curYear, curMonth, 24, 6, 2, 0, 0, time.UTC),
			ValidFrom: time.Date(curYear, curMonth, 24, 6, 0, 0, 0, time.UTC),
			ValidTo:   time.Date(curYear, curMonth, 24, 12, 0, 0, 0, time.UTC),
		}},
	{"TAF UHMD 221400Z NIL",
		&TAFMessage{rawData: "TAF UHMD 221400Z NIL",
			COR: false, AMD: false, NIL: true, Station: "UHMD", CNL: false,
			DateTime: time.Date(curYear, curMonth, 22, 14, 0, 0, 0, time.UTC),
		}},
}

func TestDecodeTAF(t *testing.T) {

	for _, pair := range tafparsetests {
		taf := NewTAF(pair.input)
		if !reflect.DeepEqual(taf, pair.expected) {
			t.Error(
				"For", pair.input,
				"expected", pair.expected,
				"got", taf,
			)
		}
	}
}