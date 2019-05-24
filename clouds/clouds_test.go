package clouds

import "testing"

type testpair struct {
	input             string
	expected          Cloud
	expectedHeightInM int
}

var tests = []testpair{

	// Type CloudType
	// height           int
	// HeightNotDefined bool
	// Cumulonimbus     bool
	// ToweringCumulus  bool
	// CBNotDefined     bool

	{"FEW005", Cloud{FEW, 5, false, false, false, false}, 150},
	{"FEW010CB", Cloud{FEW, 10, false, true, false, false}, 300},
	{"SCT018", Cloud{SCT, 18, false, false, false, false}, 550},
	{"BKN025///", Cloud{BKN, 25, false, false, false, true}, 760},
	{"OVC///", Cloud{OVC, 0, true, false, false, false}, 0},
	{"///015", Cloud{NotDefined, 15, false, false, false, false}, 460},
	{"//////", Cloud{NotDefined, 0, true, false, false, false}, 0},
	{"//////CB", Cloud{NotDefined, 0, true, true, false, false}, 0},
	{"BKN020TCU", Cloud{BKN, 20, false, false, true, false}, 610},
	{"NSC", Cloud{NSC, 0, false, false, false, false}, 0},
	{"RESHRA", Cloud{NotDefined, 0, false, false, false, false}, 0},
}

func TestParseCloud(t *testing.T) {
	for _, pair := range tests {
		v, ok := ParseCloud(pair.input)
		if ok && v != pair.expected {
			t.Error(
				"For", pair.input,
				"expected", pair.expected,
				"got", v,
			)
		}
	}
}

func TestHeightFt(t *testing.T) {
	for _, pair := range tests {
		v, ok := ParseCloud(pair.input)
		if ok && v.HeightFt() != pair.expected.height*100 {
			t.Error(
				"For", pair.input,
				"expected", pair.expected.height*100,
				"got", v.HeightFt(),
			)
		}
	}
}

func TestHeightM(t *testing.T) {
	for _, pair := range tests {
		v, ok := ParseCloud(pair.input)
		if ok && v.HeightM() != pair.expectedHeightInM {
			t.Error(
				"For", pair.input,
				"expected", pair.expectedHeightInM,
				"got", v.HeightM(),
			)
		}
	}
}
