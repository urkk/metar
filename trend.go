package metar

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/urkk/metar/clouds"
	ph "github.com/urkk/metar/phenomena"
	v "github.com/urkk/metar/visibility"
	"github.com/urkk/metar/wind"
)

// TypeTrend - type of trend: temporary or permanently expected changes
type TypeTrend string

const (
	// BECMG - Weather development (BECoMinG)
	BECMG = "BECMG"
	// TEMPO - TEMPOrary existing weather phenomena
	TEMPO = "TEMPO"
	// FM - FroM (in TAF reports)
	FM = "FM"
)

// Trend - forecast of changes for a specified period
type Trend struct {
	Type        TypeTrend
	Probability int // used only in TAFs. Maybe only 30 or 40. The PROBdd group is not used in conjunction with BECMG and FM
	// In case of in metar use values indicated time of changes. hh:mm (BECMG FM1030 TL1130)
	// In TAFs used from - until fields as date/time. ddhh/ddhh (TEMPO 2208/2218)
	FM time.Time // FroM (time)
	TL time.Time // unTiL (time)
	AT time.Time // AT time
	v.Visibility
	VerticalVisibility           int
	VerticalVisibilityNotDefined bool
	wind.Wind
	CAVOK bool
	ph.Phenomena
	clouds.Clouds
}

func parseTrendData(tokens []string) (trend *Trend) {
	trend = &Trend{}
	for count := 0; count < len(tokens); count++ {
		// PROB30 (40)
		if trend.setProbability(tokens[count]) {
			count++
		}
		// TEMPO, BECMG or FM
		if trend.setTypeOfTrend(tokens[count]) {
			count++
		}
		// AT, FM, TL used in METAR trends
		for trend.setPeriodOfChanges(tokens[count]) {
			count++
		}
		// date/time for TAF
		if trend.setDateTime(tokens[count]) {
			count++
		}
		// Wind. Only the prevailing direction.
		count += trend.ParseWind(tokens[count])

		if count < len(tokens) && tokens[count] == "CAVOK" {
			trend.CAVOK = true
			// no data after CAVOK
			return
		}
		count = decodeWeatherCondition(trend, count, tokens)
	}
	return trend
}

// Checks the string for correspondence to the forecast date/time template. Sets the date/time in case of success.
func (trend *Trend) setDateTime(input string) (ok bool) {

	regex := regexp.MustCompile(`^(\d{4})/(\d{4})`)
	matches := regex.FindStringSubmatch(input)
	if len(matches) != 0 && matches[0] != "" {
		ok = true
		t, err := parseTime(matches[1])
		if err == nil {
			trend.FM = t
		} else {
			log.Println(err)
			ok = false
		}
		t, err = parseTime(matches[2])
		if err == nil {
			trend.TL = t
		} else {
			log.Println(err)
			ok = false
		}
	}
	return
}

// parses the transmitted string, taking into account that the number of hours can be 24
func parseTime(input string) (t time.Time, err error) {
	var inputString string
	if input[2:] == "24" {
		inputString = input[:2] + "23"
		t, err = time.Parse("2006010215", CurYearStr+CurMonthStr+inputString)
		t = t.Add(time.Hour)
	} else {
		t, err = time.Parse("2006010215", CurYearStr+CurMonthStr+input)
	}
	return t, err
}

func (trend *Trend) setProbability(input string) bool {
	// Other probability values are not allowed
	if input == "PROB30" {
		trend.Probability = 30
		trend.Type = TEMPO
		return true
	} else if input == "PROB40" {
		trend.Type = TEMPO
		trend.Probability = 40
		return true
	}
	return false
}

func (trend *Trend) setPeriodOfChanges(input string) (ok bool) {
	switch input[0:2] {
	case "AT":
		timeofaction, err := time.Parse("200601021504", CurYearStr+CurMonthStr+CurDayStr+input[2:])
		if err == nil {
			trend.AT = timeofaction
			ok = true
		} else {
			log.Println(err)
		}
	case "FM":
		timeofaction, err := time.Parse("200601021504", CurYearStr+CurMonthStr+CurDayStr+input[2:])
		if err == nil {
			trend.FM = timeofaction
			ok = true
		} else {
			log.Println(err)
		}
	case "TL":
		var t time.Time
		var err error
		if input[2:] == "2400" {
			t, err = time.Parse("200601021504", CurYearStr+CurMonthStr+CurDayStr+"2300")
			t = t.Add(time.Hour)
		} else {
			t, err = time.Parse("200601021504", CurYearStr+CurMonthStr+CurDayStr+input[2:])
		}
		if err == nil {
			trend.TL = t
			ok = true
		} else {
			log.Println(err)
		}
	default:
		return
	}
	return
}

func (trend *Trend) setTypeOfTrend(input string) bool {

	if input == TEMPO || input == BECMG {
		trend.Type = TypeTrend(input)
		return true
	} else if strings.HasPrefix(input, "FM") {
		trend.Type = FM
		trend.FM, _ = time.Parse("200601021504", CurYearStr+CurMonthStr+input[2:])
		return true
	}
	return false
}

func (trend *Trend) setVerticalVisibility(input string) bool {

	if vv, nd, ok := parseVerticalVisibility(input); ok {
		trend.VerticalVisibility = vv
		trend.VerticalVisibilityNotDefined = nd
		return true
	}
	return false
}
