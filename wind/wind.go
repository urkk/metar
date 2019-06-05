package wind

import (
	"math"
	"regexp"
	"strconv"

	cnv "github.com/urkk/metar/conversion"
)

// Wind - wind on surface representation
type Wind struct {
	WindDirection int
	// The native unit of measurement is meter per second. To avoid conversion losses stored in float64
	speed        float64
	gustsSpeed   float64
	Variable     bool
	VariableFrom int
	VariableTo   int
	Above50MPS   bool
}

// SpeedKt - returns wind speed in knots. In Russian messages, the speed is specified in m/s, but it makes sense to receive it in knots for aviation purposes
func (w *Wind) SpeedKt() int {
	return int(math.Round(cnv.MpsToKts(w.speed)))
}

// SpeedMps - returns wind speed in meters per second.
func (w *Wind) SpeedMps() int {
	return int(math.Round(w.speed))
}

// GustsSpeedKt - returns gusts speed in knots.
func (w *Wind) GustsSpeedKt() int {
	return int(math.Round(cnv.MpsToKts(w.gustsSpeed)))
}

// GustsSpeedMps - returns gusts speed in meters per second.
func (w *Wind) GustsSpeedMps() int {
	return int(math.Round(w.gustsSpeed))
}

// ParseWind - identify and parses the representation of wind in the string
func (w *Wind) ParseWind(token string) (tokensused int) {

	rx := `^(\d{3}|VRB)(P)?(\d{2})(G\d\d)?(MPS|KT|KPH|KMH)\s?(\d{3}V\d{3})?`
	if matched, _ := regexp.MatchString(rx, token); !matched {
		return
	}
	tokensused = 1
	regex := regexp.MustCompile(rx)
	matches := regex.FindStringSubmatch(token)
	if matches[1] == "VRB" {
		w.Variable = true
	} else {
		w.WindDirection, _ = strconv.Atoi(matches[1])
	}
	w.Above50MPS = matches[2] != ""
	if matches[3] != "" {
		w.speed, _ = strconv.ParseFloat(matches[3], 64)
	}
	if matches[4] != "" {
		w.gustsSpeed, _ = strconv.ParseFloat(matches[4][1:], 64)
	}

	if matches[5] == "KT" {
		w.speed = cnv.KtsToMps(w.speed)
		w.gustsSpeed = cnv.KtsToMps(w.gustsSpeed)
	} else if matches[5] == "KPH" || matches[5] == "KMH" {
		w.speed = cnv.KphToMps(w.speed)
		w.gustsSpeed = cnv.KphToMps(w.gustsSpeed)
	}
	// Two tokens have been used
	if matches[6] != "" {
		tokensused++
		w.VariableFrom, _ = strconv.Atoi(matches[6][0:3])
		w.VariableTo, _ = strconv.Atoi(matches[6][4:])
	}
	return tokensused
}
