package metar

import (
	"reflect"
	"testing"
	"time"

	"github.com/urkk/metar/clouds"
	"github.com/urkk/metar/phenomena"
)

type trendparsetest struct {
	input    []string
	expected Trend
}

var trendparsetests = []trendparsetest{

	{[]string{"TEMPO", "FM1230", "TL1330", "18005MPS", "CAVOK"},
		Trend{Type: TEMPO,
			Wind:  getWind("18005MPS"),
			CAVOK: true,
			FM:    time.Date(curYear, curMonth, curDay, 12, 30, 0, 0, time.UTC),
			TL:    time.Date(curYear, curMonth, curDay, 13, 30, 0, 0, time.UTC),
		},
	},
	{[]string{"TEMPO", "AT1230", "18005MPS", "BKN003"},
		Trend{Type: TEMPO,
			Wind:   getWind("18005MPS"),
			AT:     time.Date(curYear, curMonth, curDay, 12, 30, 0, 0, time.UTC),
			Clouds: []clouds.Cloud{getCloud("BKN003")},
		},
	},
	{[]string{"TEMPO", "2509/2515", "0500", "FG", "VV003"},
		Trend{Type: TEMPO,
			FM:                 time.Date(curYear, curMonth, 25, 9, 0, 0, 0, time.UTC),
			TL:                 time.Date(curYear, curMonth, 25, 15, 0, 0, 0, time.UTC),
			Visibility:         Visibility{Distance: 500, LowerDistance: 0, LowerDirection: ""},
			Phenomena:          []phenomena.Phenomena{phenomena.Phenomena{Vicinity: false, Abbreviation: "FG", Intensity: ""}},
			VerticalVisibility: 300,
		},
	},
	{[]string{"TEMPO", "2506/2512", "3100", "-SHRA", "BR", "BKN005", "OVC020CB"},
		Trend{Type: TEMPO,
			FM:         time.Date(curYear, curMonth, 25, 6, 0, 0, 0, time.UTC),
			TL:         time.Date(curYear, curMonth, 25, 12, 0, 0, 0, time.UTC),
			Visibility: Visibility{Distance: 3100, LowerDistance: 0, LowerDirection: ""},
			Phenomena:  []phenomena.Phenomena{phenomena.Phenomena{Vicinity: false, Abbreviation: "SHRA", Intensity: "-"}, phenomena.Phenomena{Vicinity: false, Abbreviation: "BR", Intensity: ""}},
			Clouds:     []clouds.Cloud{getCloud("BKN005"), getCloud("OVC020CB")},
		},
	},
	{[]string{"TEMPO", "2509/2515", "0500", "FG", "VV///"},
		Trend{Type: TEMPO,
			FM:                           time.Date(curYear, curMonth, 25, 9, 0, 0, 0, time.UTC),
			TL:                           time.Date(curYear, curMonth, 25, 15, 0, 0, 0, time.UTC),
			Visibility:                   Visibility{Distance: 500, LowerDistance: 0, LowerDirection: ""},
			Phenomena:                    []phenomena.Phenomena{phenomena.Phenomena{Vicinity: false, Abbreviation: "FG", Intensity: ""}},
			VerticalVisibility:           0,
			VerticalVisibilityNotDefined: true,
		},
	},
	{[]string{"FM241200", "18003MPS", "CAVOK"},
		Trend{Type: FM,
			FM:    time.Date(curYear, curMonth, 24, 12, 0, 0, 0, time.UTC),
			Wind:  getWind("18003MPS"),
			CAVOK: true,
		},
	},
	{[]string{"BECMG", "2405/2407", "09003G08MPS"},
		Trend{Type: BECMG,
			FM:   time.Date(curYear, curMonth, 24, 5, 0, 0, 0, time.UTC),
			TL:   time.Date(curYear, curMonth, 24, 7, 0, 0, 0, time.UTC),
			Wind: getWind("09003G08MPS"),
		},
	},
	{[]string{"BECMG", "FM2200", "TL2400", "BKN015//"},
		Trend{Type: BECMG,
			FM:     time.Date(curYear, curMonth, curDay, 22, 0, 0, 0, time.UTC),
			TL:     time.Date(curYear, curMonth, curDay+1, 00, 0, 0, 0, time.UTC),
			Clouds: []clouds.Cloud{getCloud("BKN015//")},
		},
	},
	// misspelled time
	{[]string{"BECMG", "FM220O", "TL23O0", "BKN015//"},
		Trend{Type: BECMG,
			Clouds: []clouds.Cloud{getCloud("BKN015//")},
		},
	},
	// misspelled time
	{[]string{"BECMG", "AT220O", "TL2300", "BKN015//"},
		Trend{Type: BECMG,
			TL:     time.Date(curYear, curMonth, curDay, 23, 0, 0, 0, time.UTC),
			Clouds: []clouds.Cloud{getCloud("BKN015//")},
		},
	},
	// clock error
	{[]string{"BECMG", "2526/2526", "0500", "FG"},
		Trend{Type: BECMG,
			Visibility: Visibility{Distance: 500, LowerDistance: 0, LowerDirection: ""},
			Phenomena:  []phenomena.Phenomena{phenomena.Phenomena{Vicinity: false, Abbreviation: "FG", Intensity: ""}},
		},
	},
}

func TestParseTrendData(t *testing.T) {

	for _, pair := range trendparsetests {
		tr := *parseTrendData(pair.input)
		if !reflect.DeepEqual(tr, pair.expected) {
			t.Error(
				"For", pair.input,
				"expected", pair.expected,
				"got", tr,
			)
		}
	}
}
