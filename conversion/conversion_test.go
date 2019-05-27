package conversion

import "testing"

type testpairint struct {
	input    int
	expected int
}

var testsHPaToMmHg = []testpairint{
	{1020, 765}, //765.0650305276
	{1015, 761}, //761.3147117505
	{1011, 758}, //758.3144567288
}

func TestHPaToMmHg(t *testing.T) {
	for _, pair := range testsHPaToMmHg {
		v := HPaToMmHg(pair.input)
		if v != pair.expected {
			t.Error(
				"For", pair.input,
				"expected", pair.expected,
				"got", v,
			)
		}
	}
}

var testsMmHgToHPa = []testpairint{
	{768, 1024}, //1023.91296
	{764, 1019}, //1018.58008
	{758, 1011}, //1010.58076
}

func TestMmHgToHPa(t *testing.T) {
	for _, pair := range testsMmHgToHPa {
		v := MmHgToHPa(pair.input)
		if v != pair.expected {
			t.Error(
				"For", pair.input,
				"expected", pair.expected,
				"got", v,
			)
		}
	}
}

type testpairDirection struct {
	input    int
	expected string
}

var testsDirection = []testpairDirection{
	{360, "N"},
	{30, "NE"},
	{275, "W"},
	{0, "N"},
}

func TestDirectionToCardinalDirection(t *testing.T) {
	for _, pair := range testsDirection {
		v := DirectionToCardinalDirection(pair.input)
		if v != pair.expected {
			t.Error(
				"For", pair.input,
				"expected", pair.expected,
				"got", v,
			)
		}
	}
}

type testRelativeHumidity struct {
	temp, dewpoint, rh int
}

var testsRelativeHumidity = []testRelativeHumidity{
	{20, 7, 43},
	{31, 7, 22},
	{25, 25, 100},
	{25, 19, 69},
	{-7, -13, 62},
	{0, -9, 51},
}

func TestCalcRelativeHumidity(t *testing.T) {
	for _, data := range testsRelativeHumidity {
		rh := CalcRelativeHumidity(data.temp, data.dewpoint)
		if rh != data.rh {
			t.Error(
				"For t", data.temp, " dew point ", data.dewpoint,
				"expected", data.rh,
				"got", rh,
			)
		}
	}
}
