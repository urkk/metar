package conversion

import "testing"

type testpairint struct {
	input    int
	expected int
}

var testsHPaToMmHg = []testpairint{
	{1024, 768},
	{1018, 764},
	{1011, 758},
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
