package visibility

import (
	"regexp"
	"strconv"
	"strings"

	cnv "github.com/urkk/metar/conversion"
)

// Unit of measurement.
type Unit string

const (
	// M - meters
	M = "M"
	// FT - feet
	FT = "FT"
	// SM - statute miles
	SM = "SM"
)

// Distance in units of measure
type Distance struct {
	// By default, meters. Or feet in US RVR. Both integer
	Value int
	// Used only for horizontal visibility in miles
	FractionValue float64
	Unit
}

// BaseVisibility - the basis of visibility: one measurement
type BaseVisibility struct {
	Distance
	AboveMax bool // more than reported value (P5000, P6SM...)
	BelowMin bool // less than reported value (M1/4SM, M0050...)
}

// Visibility - prevailing visibility
type Visibility struct {
	BaseVisibility
	// sector visibility
	LowerDistance  Distance
	LowerDirection string
}

// Meters - returns the distance in meters
func (d *Distance) Meters() (result int) {
	switch d.Unit {
	case M:
		result = d.Value
	case FT:
		result = cnv.FtToM(d.Value)
	case SM:
		result = cnv.SMileToM(float64(d.Value) + d.FractionValue)
	default:
		result = d.Value
	}
	return
}

// Feet - returns the distance in feet
func (d *Distance) Feet() (result int) {
	switch d.Unit {
	case M:
		result = cnv.MToFt(d.Value)
	case FT:
		result = d.Value
	case SM:
		result = cnv.SMileToFt(float64(d.Value) + d.FractionValue)
	default:
		result = cnv.MToFt(d.Value)
	}
	return
}

// Miles - returns the distance in miles
func (d *Distance) Miles() (result float64) {
	switch d.Unit {
	case M:
		result = cnv.MToSMile(d.Value)
	case FT:
		result = float64(cnv.FtToSMile(int(d.Value)))
	case SM:
		result = float64(d.Value) + d.FractionValue
	default:
		result = cnv.MToSMile(d.Value)
	}
	return
}

// ParseVisibility - identify and parses the representation oh horizontal visibility
func (v *Visibility) ParseVisibility(input []string) (tokensused int) {
	inputstring := strings.Join(input, " ")
	metric := regexp.MustCompile(`^(P|M)?(\d{4})(\s|$)((\d{4})(NE|SE|NW|SW|N|E|S|W))?`)
	// In US and CA sector visibility reported in the remarks section. (as VIS NW-SE 1/2; VIS NE 2 1/2 etc)
	imperial := regexp.MustCompile(`^(P|M)?(\d{1,2}|\d(\s)?)?((\d)/(\d))?SM`)

	switch {
	case metric.MatchString(inputstring):
		tokensused = 1
		v.Distance.Unit = M
		matches := metric.FindStringSubmatch(inputstring)
		v.BelowMin = matches[1] == "M"
		v.AboveMax = matches[1] == "P"
		v.Distance.Value, _ = strconv.Atoi(matches[2])
		if matches[4] != "" {
			v.LowerDistance.Value, _ = strconv.Atoi(matches[5])
			v.LowerDistance.Unit = M
			v.LowerDirection = matches[6]
			tokensused++
		}
	case imperial.MatchString(inputstring):
		tokensused = 1
		matches := imperial.FindStringSubmatch(inputstring)
		v.BelowMin = matches[1] == "M"
		v.AboveMax = matches[1] == "P"
		if matches[2] != "" {
			v.Distance.Value, _ = strconv.Atoi(strings.TrimSpace(matches[2]))
		}
		if matches[5] != "" && matches[6] != "" {
			numerator, _ := strconv.Atoi(matches[5])
			denominator, _ := strconv.Atoi(matches[6])
			if denominator != 0 {
				v.Distance.FractionValue += float64(numerator) / float64(denominator)
			}
		}
		v.Distance.Unit = SM
		if matches[3] == " " {
			tokensused++
		}
	default:
		return
	}
	return
}
