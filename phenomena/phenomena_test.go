package phenomena

import "testing"

type testpair struct {
	input    string
	expected *Phenomenon
	correct  bool
}

var tests = []testpair{

	// Vicinity     bool
	// Intensity    Intensity
	// Abbreviation string

	{"PLRADZ", &Phenomenon{false, Moderate, "PLRADZ"}, true},
	{"+SHRASNGS", &Phenomenon{false, Heavy, "SHRASNGS"}, true},
	{"VCBLDU", &Phenomenon{true, Moderate, "BLDU"}, true},
	// not can VC
	{"VCSNDZ", &Phenomenon{true, Moderate, "VCSNDZ"}, false},
	{"SHGRRA", &Phenomenon{false, Moderate, "SHGRRA"}, true},
	// "Light" not be applicable
	{"-FC", &Phenomenon{false, Light, "FC"}, false},
	{"-PLDZ", &Phenomenon{false, Light, "PLDZ"}, true},
}

var recenttests = []testpair{

	// Vicinity     bool
	// Intensity    Intensity
	// Abbreviation string

	{"REFZDZ", &Phenomenon{false, Moderate, "FZDZ"}, true},
	{"+REFZDZ", &Phenomenon{false, Heavy, "FZDZ"}, false},
	{"RERASN", &Phenomenon{false, Moderate, "RASN"}, true},
}

func TestParsePhenomena(t *testing.T) {
	arr := &Phenomena{}
	for _, pair := range tests {
		ph := ParsePhenomena(pair.input)
		if pair.correct {
			if *ph != *pair.expected {
				t.Error(
					"For", pair.input,
					"expected", pair.expected,
					"got", ph,
				)
			}
			if !arr.AppendPhenomena(pair.input) {
				t.Error("For", pair.input, "error append phenomenon")
			}
		} else if !pair.correct {
			if ph != nil || arr.AppendPhenomena(pair.input) {
				t.Error("false positive at " + pair.input)
			}

		}
	}
}

func TestParseRecentPhenomena(t *testing.T) {
	arr := &Phenomena{}
	for _, pair := range recenttests {
		ph := ParseRecentPhenomena(pair.input)
		if pair.correct {
			if *ph != *pair.expected {
				t.Error(
					"For", pair.input,
					"expected", pair.expected,
					"got", ph,
				)
			}
			if !arr.AppendRecentPhenomena(pair.input) {
				t.Error("For", pair.input, "error append phenomenon")
			}
		} else if !pair.correct {
			if ph != nil || arr.AppendRecentPhenomena(pair.input) {
				t.Error("false positive at " + pair.input)
			}
		}
	}
}
