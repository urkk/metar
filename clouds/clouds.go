// Package clouds describe cloud amount and height for metar message decoding
package clouds

import (
	cnv "github.com/urkk/metar/conversion"
	"regexp"
	"strconv"
)

// Cloud - cloud representation
type Cloud struct {
	Type CloudType
	// the height is stored in hundreds of feet, as Flight Level
	height           int
	HeightNotDefined bool
	Cumulonimbus     bool
	ToweringCumulus  bool
	CBNotDefined     bool
}

// CloudType - Cloud amounts
type CloudType string

// predefined cloud amount code
const (
	FEW        = "FEW"
	SCT        = "SCT" //scattered
	BKN        = "BKN" //broken
	OVC        = "OVC" //overcast
	NSC        = "NSC" //nil significant cloud
	NCD        = "NCD" //nil cloud detected for automated METAR station
	SKC        = "SKC" //sky is clear
	CLR        = "CLR" //sky is clear for automated station
	NotDefined = "///"
)

// HeightM - returns height above surface of the lower base of cloudiness in meters
func (cl Cloud) HeightM() int {
	return cnv.FtToM(cl.height * 100)
}

// HeightFt - returns height above surface of the lower base of cloudiness in feet
func (cl Cloud) HeightFt() int {
	return cl.height * 100
}

// ParseCloud - identify and parses the representation of cloudiness in the string
func ParseCloud(token string) (cl Cloud, ok bool) {

	pattern := `^(FEW|SCT|BKN|OVC|NSC|SKC|NCD|CLR|///)(\d{3}|///)?(TCU|CB|///)?`
	if matched, _ := regexp.MatchString(pattern, token); !matched {
		return cl, false
	}
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(token)

	cl.Type = CloudType(matches[1])
	if cl.Type == NSC || cl.Type == NCD || cl.Type == CLR || cl.Type == SKC { // no clouds
		return cl, true
	}

	if matches[2] != "///" {
		cl.height, _ = strconv.Atoi(matches[2])
	} else {
		cl.HeightNotDefined = true
	}
	cl.CBNotDefined = matches[3] == "///"
	cl.Cumulonimbus = matches[3] == "CB"
	cl.ToweringCumulus = matches[3] == "TCU"
	return cl, true

}
