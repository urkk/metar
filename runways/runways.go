package runways

import (
	"fmt"
	"regexp"
	"strconv"
)

type VRTendency string

const (
	NotDefined = ""
	U          = "U" //upward
	N          = "N" //no distinct
	D          = "D" //downward
)

type RunwayDesignator struct {
	Number     string
	AllRunways bool
}

func NewRD(number string) (rd RunwayDesignator) {
	pattern := `^(R)?(\d{2})[LCR]?`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(number)
	course, _ := strconv.Atoi(matches[2])
	rd.Number = number
	if course == 88 {
		rd.AllRunways = true
	} else if course >= 51 && course <= 86 {
		rd.Number = matches[1] + fmt.Sprintf("%02d", course-50) + "R"
	}
	// R99 - repetition of the last message received because there is no new information received. Not implemented.
	return rd
}

type VisualRange struct {
	Designator RunwayDesignator
	Distance   int
	AboveMax   bool
	BelowMin   bool
	Trend      VRTendency
}

type State struct {
	Designator                RunwayDesignator
	TypeOfCoverage            int
	TypeOfCoverageNotDef      bool
	DimensionOfCoverage       int
	DimensionOfCoverageNotDef bool
	HeightOfCoverage          int
	HeightOfCoverageNotDef    bool
	// Friction coefficient and braking action
	BrakingConditions           int
	BrakingConditionsNotDefined bool
	CLRD                        bool
	SNOCLO                      bool
}

func ParseVisibility(token string) (v VisualRange, result bool) {
	// TODO 0800V1000FT			R27/0150V0300U
	pattern := `^(R\d{2}[LCR]?)/(M|P)?(\d{4})(U|D|N)?`
	if matched, _ := regexp.MatchString(pattern, token); !matched {
		return v, false
	}
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(token)
	v.Designator = NewRD(matches[1][1:])
	v.AboveMax = matches[2] == "P" // plus
	v.BelowMin = matches[2] == "M" // minus
	v.Distance, _ = strconv.Atoi(matches[3])
	v.Trend = VRTendency(matches[4])
	return v, true
}

func ParseState(token string) (s State, result bool) {

	pattern := `^(R\d{2}[LCR]?)/((\d|\/)(\d|\/)(\d\d|\/\/)|CLRD)?(\d\d|\/\/)(D)?$`
	if matched, _ := regexp.MatchString(pattern, token); !matched {
		return s, false
	}
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(token)
	s.Designator = NewRD(matches[1][1:])

	if matches[6] == "//" {
		s.BrakingConditionsNotDefined = true
	} else {
		s.BrakingConditions, _ = strconv.Atoi(matches[6])
	}
	if matches[2] == "CLRD" || matches[7] == "D" {
		s.CLRD = true
		return s, true
	}
	if matches[3] == "/" {
		s.TypeOfCoverageNotDef = true
	} else {
		s.TypeOfCoverage, _ = strconv.Atoi(matches[3])
	}
	if matches[4] == "/" {
		s.DimensionOfCoverageNotDef = true
	} else {
		s.DimensionOfCoverage, _ = strconv.Atoi(matches[4])
	}
	if matches[5] == "//" {
		s.HeightOfCoverageNotDef = true
	} else {
		s.HeightOfCoverage, _ = strconv.Atoi(matches[5])
	}
	return s, true

}
