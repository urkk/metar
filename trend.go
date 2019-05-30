package metar

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/urkk/metar/clouds"
	ph "github.com/urkk/metar/phenomena"
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
	FM                           time.Time // FroM (time)
	TL                           time.Time // unTiL (time)
	AT                           time.Time // AT time
	Visibility                   Visibility
	VerticalVisibility           int
	VerticalVisibilityNotDefined bool
	Wind                         wind.Wind
	CAVOK                        bool
	Phenomena                    []ph.Phenomena
	Clouds                       []clouds.Cloud
	setData                      func(input string) bool
}

func parseTrendData(tokens []string) (trend *Trend) {
	trend = new(Trend)
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
		if wnd, tokensused := wind.ParseWind(tokens[count]); tokensused > 0 {
			count += tokensused
			trend.Wind = wnd
		}
		if count < len(tokens) && tokens[count] == "CAVOK" {
			trend.CAVOK = true
			// no data after CAVOK
		} else {
			// Horizontal visibility. The distance and direction of the least visibility is not predicted
			if vis, tokensused := ParseVisibility(strings.Join(tokens[count:], " ")); tokensused > 0 {
				trend.Visibility = vis
				count += tokensused
			}
			// Weather or NSW - no significant weather
			for count < len(tokens) && trend.appendPhenomena(tokens[count]) {
				count++
			}
			// Vertical visibility
			if count < len(tokens) && trend.setVerticalVisibility(tokens[count]) {
				count++
			}
			// Clouds. No further information after the clouds in trend
			for count < len(tokens) && trend.appendCloud(tokens[count]) {
				count++
			}
		}
	}
	return trend
}

// Checks the string for correspondence to the forecast date/time template. Sets the date/time in case of success.
func (trend *Trend) setDateTime(input string) bool {

	regex := regexp.MustCompile(`^(\d{4})/(\d{4})`)
	matches := regex.FindStringSubmatch(input)
	if len(matches) != 0 && matches[0] != "" {
		starttime := matches[1]
		endtime := matches[2]
		// hours maybe 24
		if starttime[2:] == "24" {
			starttime = starttime[:2] + "23"
			trend.FM, _ = time.Parse("2006010215", CurYearStr+CurMonthStr+starttime)
			trend.FM = trend.FM.Add(time.Hour)
		} else {
			trend.FM, _ = time.Parse("2006010215", CurYearStr+CurMonthStr+starttime)
		}
		if endtime[2:] == "24" {
			endtime = endtime[:2] + "23"
			trend.TL, _ = time.Parse("2006010215", CurYearStr+CurMonthStr+endtime)
			trend.TL = trend.TL.Add(time.Hour)
		} else {
			trend.TL, _ = time.Parse("2006010215", CurYearStr+CurMonthStr+endtime)
		}
		return true
	}
	return false
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

func (trend *Trend) setPeriodOfChanges(input string) bool {

	switch input[0:2] {
	case "AT":
		timeofaction, err := time.Parse("200601021504", CurYearStr+CurMonthStr+CurDayStr+input[2:])
		if err != nil {
			return false
		}
		trend.AT = timeofaction
		return true
	case "FM":
		timeofaction, err := time.Parse("200601021504", CurYearStr+CurMonthStr+CurDayStr+input[2:])
		if err != nil {
			return false
		}
		trend.FM = timeofaction
		return true
	case "TL":
		if input[2:] == "2400" {
			trend.TL, _ = time.Parse("200601021504", CurYearStr+CurMonthStr+CurDayStr+"2300")
			trend.TL = trend.TL.Add(time.Hour)
		} else {
			timeofaction, err := time.Parse("200601021504", CurYearStr+CurMonthStr+CurDayStr+input[2:])
			if err != nil {
				return false
			}
			trend.TL = timeofaction
			return true
		}
	default:
		return false
	}
	return false
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

	regex := regexp.MustCompile(`^VV(\d{3}|///)`)
	matches := regex.FindStringSubmatch(input)
	if len(matches) != 0 && matches[0] != "" {
		if matches[1] != "///" {
			trend.VerticalVisibility, _ = strconv.Atoi(matches[1])
			trend.VerticalVisibility *= 100
		} else {
			trend.VerticalVisibilityNotDefined = true
		}
		return true
	}
	return false
}

func (trend *Trend) appendCloud(input string) bool {

	if cl, ok := clouds.ParseCloud(input); ok {
		trend.Clouds = append(trend.Clouds, cl)
		return true
	}
	return false
}

func (trend *Trend) appendPhenomena(input string) bool {

	if p := ph.ParsePhenomena(input); p != nil {
		trend.Phenomena = append(trend.Phenomena, *p)
		return true
	}
	return false
}
