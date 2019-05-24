// Package clouds describe cloud amount and height for metar message decoding
package clouds

import (
	"regexp"
	"strconv"

	cnv "github.com/urkk/metar/conversion"
)

type Cloud struct {
	Type CloudType
	// the height is stored in hundreds of feet, as Flight Level
	height           int
	HeightNotDefined bool
	Cumulonimbus     bool
	ToweringCumulus  bool
	CBNotDefined     bool
}

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

// returns height above surface of the lower base of cloudiness in meters
func (cl Cloud) HeightM() int {
	return cnv.FtToM(cl.height * 100)
}

// returns height above surface of the lower base of cloudiness in feet
func (cl Cloud) HeightFt() int {
	return cl.height * 100
}

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
