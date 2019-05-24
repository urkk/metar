package phenomena

import "testing"

type testpair struct {
	input    string
	expected *Phenomena
	correct  bool
}

var tests = []testpair{

	// Vicinity     bool
	// Intensity    Intensity
	// Abbreviation string

	{"PLRADZ", &Phenomena{false, Moderate, "PLRADZ"}, true},
	{"+SHRASNGS", &Phenomena{false, Heavy, "SHRASNGS"}, true},
	{"VCBLDU", &Phenomena{true, Moderate, "BLDU"}, true},
	// not can VC
	{"VCSNDZ", &Phenomena{true, Moderate, "VCSNDZ"}, false},
	{"SHGRRA", &Phenomena{false, Moderate, "SHGRRA"}, true},
	// "Light" not be applicable
	{"-FC", &Phenomena{false, Light, "FC"}, false},
	{"-PLDZ", &Phenomena{false, Light, "PLDZ"}, true},
}

var recenttests = []testpair{

	// Vicinity     bool
	// Intensity    Intensity
	// Abbreviation string

	{"REFZDZ", &Phenomena{false, Moderate, "FZDZ"}, true},
	{"+REFZDZ", &Phenomena{false, Heavy, "FZDZ"}, false},
	{"RERASN", &Phenomena{false, Moderate, "RASN"}, true},
}

func TestParsePhenomena(t *testing.T) {
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
		} else if !pair.correct && ph != nil {
			t.Error("false positive at " + pair.input)
		}
	}
}

func TestParseRecentPhenomena(t *testing.T) {
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
		} else if !pair.correct && ph != nil {
			t.Error("false positive at " + pair.input)
		}
	}
}
