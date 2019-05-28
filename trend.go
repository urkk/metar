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
	// In case of in metar use values indicated time of changes. hhmm (BECMG FM1030 TL1130)
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
}

func parseTrendData(tokens []string) (trend *Trend) {
	trend = new(Trend)
	index := 0
	for index < len(tokens) {
		if tokens[index] == "PROB30" {
			trend.Probability = 30
			trend.Type = TEMPO
			index++
		} else if tokens[index] == "PROB40" {
			trend.Type = TEMPO
			trend.Probability = 40
			index++
		}
		if tokens[index] == TEMPO || tokens[index] == BECMG {
			trend.Type = TypeTrend(tokens[index])
			index++
		} else if strings.HasPrefix(tokens[index], "FM") {
			trend.Type = FM
			trend.FM, _ = time.Parse("200601021504", CurYearStr+CurMonthStr+tokens[index][2:])
			index++
		}
		// AT, FM, TL used in METAR trends
		if tokens[index][0:2] == "AT" {
			trend.AT, _ = time.Parse("200601021504", CurYearStr+CurMonthStr+CurDayStr+tokens[index][2:])
			index++
		}
		if tokens[index][0:2] == "FM" {
			trend.FM, _ = time.Parse("200601021504", CurYearStr+CurMonthStr+CurDayStr+tokens[index][2:])
			index++
		}
		if tokens[index][0:2] == "TL" {
			trend.TL, _ = time.Parse("200601021504", CurYearStr+CurMonthStr+CurDayStr+tokens[index][2:])
			index++
		}
		// date/time for TAF
		if trend.checkDateTime(tokens[index]) {
			index++
		}
		// Wind. Only the prevailing direction.
		if wnd, tokensused := wind.ParseWind(tokens[index]); tokensused > 0 {
			index += tokensused
			trend.Wind = wnd
		}
		if index < len(tokens) && tokens[index] == "CAVOK" {
			trend.CAVOK = true
			return trend
		} else {
			// Horizontal visibility. The distance and direction of the least visibility is not predicted
			if vis, tokensused := ParseVisibility(strings.Join(tokens[index:], " ")); tokensused > 0 {
				trend.Visibility = vis
				index += tokensused
			}
			// Weather or NSW - no significant weather
			for index < len(tokens) {
				if p := ph.ParsePhenomena(tokens[index]); p != nil {
					trend.Phenomena = append(trend.Phenomena, *p)
					index++
				} else {
					break
				}
			}
			// Vertical visibility
			if index < len(tokens) {
				regex := regexp.MustCompile(`^VV(\d{3}|///)`)
				matches := regex.FindStringSubmatch(tokens[index])
				if len(matches) != 0 && matches[0] != "" {
					if matches[1] != "///" {
						trend.VerticalVisibility, _ = strconv.Atoi(matches[1])
						trend.VerticalVisibility *= 100
					} else {
						trend.VerticalVisibilityNotDefined = true
					}
					index++
				}
			}
			// Clouds. No further information after the clouds in trend
			for index < len(tokens) {
				if cl, ok := clouds.ParseCloud(tokens[index]); ok {
					trend.Clouds = append(trend.Clouds, cl)
					index++
				} else {
					break
				}
			}
		}
		index++
	}
	return trend
}

func (trend *Trend) checkDateTime(input string) bool {

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
