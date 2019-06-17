package runways

import (
	"fmt"
	"regexp"
	"strconv"

	vis "github.com/urkk/metar/visibility"
)

// VRTendency - average runway visual range tendency
type VRTendency string

const (
	//NotDefined - no changes are reported
	NotDefined = ""
	//U - upward
	U = "U"
	//N - no distinct
	N = "N"
	//D - downward
	D = "D"
)

// RunwayDesignator - two-digit runway number
type RunwayDesignator struct {
	Number     string
	AllRunways bool
}

// NewRD - construct new runway designator
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

// VisualRange - describes the horizontal distance you can expect to see down a runway
type VisualRange struct {
	Designator     RunwayDesignator
	Visibility     vis.BaseVisibility
	UpToVisibility vis.BaseVisibility
	Trend          VRTendency
}

// State - runway condition representation
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

// ParseVisualRange - identify and parses the representation of runway visual range
func ParseVisualRange(token string) (v VisualRange, result bool) {
	// TODO 0800V1000FT			R27/0150V0300U
	//pattern := `^(R\d{2}[LCR]?)/(M|P)?(\d{4})(U|D|N)?`
	regex := regexp.MustCompile(`^(R\d{2}[LCR]?)/(M|P)?(\d{4})(V(M|P)?(\d{4}))?(FT)?/?(U|D|N)?`)
	if matched := regex.MatchString(token); !matched {
		return v, false
	}
	matches := regex.FindStringSubmatch(token)
	v.Designator = NewRD(matches[1][1:])
	v.Visibility.AboveMax = matches[2] == "P" // plus
	v.Visibility.BelowMin = matches[2] == "M" // minus
	v.Visibility.Distance.Value, _ = strconv.Atoi(matches[3])
	if matches[7] == "FT" {
		v.Visibility.Distance.Unit = vis.FT
	}
	if matches[4] != "" {
		if matches[7] == "FT" {
			v.UpToVisibility.Unit = vis.FT
		}
		v.UpToVisibility.Value, _ = strconv.Atoi(matches[6])
		v.UpToVisibility.AboveMax = matches[5] == "P" // plus
		v.UpToVisibility.BelowMin = matches[5] == "M" // minus
	}
	v.Trend = VRTendency(matches[8])
	return v, true
}

// ParseState - identify and parses the representation of runway condition
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
