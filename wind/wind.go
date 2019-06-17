package wind

import (
	"math"
	"regexp"
	"strconv"

	cnv "github.com/urkk/metar/conversion"
)

type speedUnit string

const (
	// MPS - meters per second
	MPS = "MPS"
	// KT - knots
	KT = "KT"
	// KPH - kilometers per hour
	KPH = "KPH"
	// KMH - kilometers per hour
	KMH = "KMH"
)

// Wind - wind on surface representation
type Wind struct {
	WindDirection       int
	speed               int
	gustsSpeed          int
	Variable            bool
	VariableFrom        int
	VariableTo          int
	Above50MPS          bool
	SpeedNotDefined     bool
	DirectionNotDefined bool
	unit                speedUnit
}

// SpeedKt - returns wind speed in knots. In Russian messages, the speed is specified in m/s, but it makes sense to receive it in knots for aviation purposes
func (w *Wind) SpeedKt() (speed int) {
	return getSpeedValue(w.speed, w.unit, KT)
}

// SpeedMps - returns wind speed in meters per second.
func (w *Wind) SpeedMps() (speed int) {
	return getSpeedValue(w.speed, w.unit, MPS)
}

// GustsSpeedKt - returns gusts speed in knots.
func (w *Wind) GustsSpeedKt() (speed int) {
	return getSpeedValue(w.gustsSpeed, w.unit, KT)
}

// GustsSpeedMps - returns gusts speed in meters per second.
func (w *Wind) GustsSpeedMps() (speed int) {
	return getSpeedValue(w.gustsSpeed, w.unit, MPS)
}

func getSpeedValue(speed int, unit speedUnit, needUnit speedUnit) (result int) {
	if unit == needUnit {
		result = speed
	} else {
		switch {
		case unit == MPS && needUnit == KT:
			result = int(math.Round(cnv.MpsToKts(float64(speed))))
		case unit == KT && needUnit == MPS:
			result = int(math.Round(cnv.KtsToMps(float64(speed))))
		case (unit == KMH || unit == KPH) && needUnit == MPS:
			result = int(math.Round(cnv.KphToMps(speed)))
		case (unit == KMH || unit == KPH) && needUnit == KT:
			result = int(math.Round(cnv.KphToKts(speed)))
		}
	}
	return
}

// ParseWind - identify and parses the representation of wind in the string
func (w *Wind) ParseWind(token string) (tokensused int) {

	rx := `^(\d{3}|VRB|///)(P)?(\d{2}|//)(G\d\d)?(MPS|KT|KPH|KMH)\s?(\d{3}V\d{3})?`
	if matched, _ := regexp.MatchString(rx, token); !matched {
		return
	}
	tokensused = 1
	regex := regexp.MustCompile(rx)
	matches := regex.FindStringSubmatch(token)
	w.Variable = matches[1] == "VRB"
	w.DirectionNotDefined = matches[1] == "///"
	if !w.Variable && !w.DirectionNotDefined {
		w.WindDirection, _ = strconv.Atoi(matches[1])
	}
	w.Above50MPS = matches[2] != ""
	w.SpeedNotDefined = matches[3] == "//"
	if matches[3] != "" && !w.SpeedNotDefined {
		w.speed, _ = strconv.Atoi(matches[3])
	}
	if matches[4] != "" {
		w.gustsSpeed, _ = strconv.Atoi(matches[4][1:])
	}
	w.unit = speedUnit(matches[5])
	// Two tokens have been used
	if matches[6] != "" {
		tokensused++
		w.VariableFrom, _ = strconv.Atoi(matches[6][0:3])
		w.VariableTo, _ = strconv.Atoi(matches[6][4:])
	}
	return tokensused
}
