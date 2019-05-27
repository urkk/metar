package wind

import "testing"

type parsetest struct {
	input    string
	expected Wind
}

var parsetests = []parsetest{
	// Speed         int
	// WindDirection int
	// GustsSpeed    int
	// Variable      bool
	// VariableFrom  int
	// VariableTo    int
	// Above50MPS    bool
	{"31005MPS", Wind{310, 5, 0, false, 0, 0, false}},
	{"31010KPH", Wind{310, 2.7777777777777777, 0, false, 0, 0, false}},
	{"VRB15MPS", Wind{0, 15, 0, true, 0, 0, false}},
	{"00000MPS", Wind{0, 0, 0, false, 0, 0, false}},
	{"240P49MPS", Wind{240, 49, 0, false, 0, 0, true}},
	{"04008G20MPS", Wind{40, 8, 20, false, 0, 0, false}},
	{"22003G08MPS 280V350", Wind{220, 3, 8, false, 280, 350, false}},
	{"14010KT", Wind{140, 5.144456333854638, 0, false, 0, 0, false}},
	{"BKN020", Wind{0, 0, 0, false, 0, 0, false}},
}

type functest struct {
	input            Wind
	expectedKt       int
	expectedMps      int
	expectedGustsKt  int
	expectedGustsMps int
}

var functests = []functest{
	{Wind{0, 10, 0, false, 0, 0, false}, 19, 10, 0, 0},
	{Wind{0, 0, 0, false, 0, 0, false}, 0, 0, 0, 0},
	{Wind{0, 15, 0, false, 0, 0, false}, 29, 15, 0, 0},
	{Wind{0, 7.33, 0, false, 0, 0, false}, 14, 7, 0, 0},
	{Wind{0, 10, 10, false, 0, 0, false}, 19, 10, 19, 10},
	{Wind{0, 15, 15, false, 0, 0, false}, 29, 15, 29, 15},
}

func TestParseWind(t *testing.T) {
	for _, pair := range parsetests {
		v, ok, _ := ParseWind(pair.input)
		if ok && v != pair.expected {
			t.Error(
				"For", pair.input,
				"expected", pair.expected,
				"got", v,
			)
		}
	}
}

func TestSpeedKt(t *testing.T) {
	for _, pair := range functests {
		vkt := pair.input.SpeedKt()
		if vkt != pair.expectedKt {
			t.Error(
				"For", pair.input,
				"expected", pair.expectedKt,
				"got", vkt,
			)
		}
	}
}

func TestSpeedMps(t *testing.T) {
	for _, pair := range functests {
		vmps := pair.input.SpeedMps()
		if vmps != pair.expectedMps {
			t.Error(
				"For", pair.input,
				"expected", pair.expectedMps,
				"got", vmps,
			)
		}
	}
}

func TestGustsSpeedMps(t *testing.T) {
	for _, pair := range functests {
		vmps := pair.input.GustsSpeedMps()
		if vmps != pair.expectedGustsMps {
			t.Error(
				"For", pair.input,
				"expected", pair.expectedGustsMps,
				"got", vmps,
			)
		}
	}
}

func TestGustsSpeedKt(t *testing.T) {
	for _, pair := range functests {
		vmps := pair.input.GustsSpeedKt()
		if vmps != pair.expectedGustsKt {
			t.Error(
				"For", pair.input,
				"expected", pair.expectedGustsKt,
				"got", vmps,
			)
		}
	}
}
