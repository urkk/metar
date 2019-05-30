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

// TemperatureForecast - Forecast Max and Min temperature
type TemperatureForecast struct {
	Temp     int
	DateTime time.Time
	IsMax    bool
	IsMin    bool
}

// TAFMessage - Terminal Aerodrome Forecast struct
type TAFMessage struct {
	rawData   string    // The raw TAF
	COR       bool      // Correction of forecast due to a typo
	AMD       bool      // Amended forecast
	NIL       bool      // event of missing TAF
	Station   string    // 4-letter ICAO station identifier
	DateTime  time.Time // Time( in ISO8601 date/time format) this TAF was issued
	ValidFrom time.Time
	ValidTo   time.Time
	CNL       bool // The previously issued TAF for the period was cancelled

	Wind wind.Wind //	Surface wind
	// Ceiling And Visibility OK, indicating no cloud below 5,000 ft (1,500 m) or the highest minimum sector
	// altitude and no cumulonimbus or towering cumulus at any level, a visibility of 10 km (6 mi) or more and no significant weather change.
	CAVOK              bool
	Visibility         Visibility            // Horizontal visibility
	Phenomena          []ph.Phenomena        // Present Weather
	VerticalVisibility int                   // Vertical visibility (ft)
	Clouds             []clouds.Cloud        // Cloud amount and height
	Temperature        []TemperatureForecast // Temperature extremes
	// Prevision
	TREND []Trend

	NotDecodedTokens []string
}

// NewTAF - creates a new TAF forecast based on the original message
func NewTAF(inputtext string) *TAFMessage {
	t := &TAFMessage{
		rawData: inputtext,
	}
	headerRx := myRegexp{regexp.MustCompile(`^(?P<taf>TAF\s)?(?P<cor>COR\s)?(?P<amd>AMD\s)?(?P<station>\w{4})\s(?P<time>\d{6}Z)(?P<nil>\sNIL)?(\s(?P<from>\d{4})/(?P<to>\d{4}))?(?P<cnl>\sCNL)?`)}
	headermap := headerRx.FindStringSubmatchMap(t.rawData)
	t.Station = headermap["station"]
	t.DateTime, _ = time.Parse("200601021504Z", CurYearStr+CurMonthStr+headermap["time"])
	t.COR = headermap["cor"] != ""
	t.AMD = headermap["amd"] != ""
	t.NIL = headermap["nil"] != ""
	t.CNL = headermap["cnl"] != ""
	if t.Station == "" && t.DateTime.IsZero() {
		//not valid message?
		t.NotDecodedTokens = append(t.NotDecodedTokens, t.rawData)
		return t
	}
	if t.NIL { // End of TAF, if the forecast is lost
		return t
	}
	t.setTimeRange(headermap["from"], headermap["to"])
	if t.CNL { // End of TAF, if the forecast is cancelled
		return t
	}
	tokens := strings.Split(t.rawData, " ")

	position := 0
	for key, value := range headermap {
		if value != "" && key != "to" { // field "from" and "to" - it's one token (DDhh/DDhh), and they are mandatory.
			count++
		}
	}
	endposition := t.findTrendsInMessage(tokens, position)
	t.decodeTAF(tokens[position:endposition])
	return t
}

// RAW - returns the original message text
func (t *TAFMessage) RAW() string { return t.rawData }

func (t *TAFMessage) decodeTAF(tokens []string) {

	totalcount := len(tokens)

	for count := 0; count < totalcount; {
		// Wind - Visibility - Weather - Sky Condition

		// Surface wind
		if wnd, tokensused := wind.ParseWind(tokens[count]); tokensused > 0 {
			t.Wind = wnd
			count += tokensused
		}
		if tokens[count] == "CAVOK" {
			t.CAVOK = true
			count++
		} else {
			// Horizontal visibility
			if vis, tokensused := ParseVisibility(tokens[count]); tokensused > 0 {
				t.Visibility = vis
				count += tokensused
			}
			// Weather
			for count < len(tokens) && t.appendPhenomena(tokens[count]) {
				count++
			}
			// Vertical visibility
			if t.setVerticalVisibility(tokens[count]) {
				count++
			}
			// Cloudiness description
			for count < len(tokens) && t.appendCloud(tokens[count]) {
				count++
			}
		} // !CAVOK

		// Temperature
		for count < len(tokens) && t.addTempForecast(tokens[count]) {
			count++
		}
		// The token is not recognized or is located in the wrong position
		if count < totalcount {
			t.NotDecodedTokens = append(t.NotDecodedTokens, tokens[count])
			count++
		}
	} // End main section
}

func (t *TAFMessage) addTempForecast(input string) bool {

	regex := regexp.MustCompile(`^T(X|N)(M)?(\d\d)\/(\d{4}Z)`)
	matches := regex.FindStringSubmatch(input)
	if len(matches) > 0 {
		tempf := new(TemperatureForecast)
		tempf.Temp, _ = strconv.Atoi(matches[3])
		if matches[2] == "M" {
			tempf.Temp = -tempf.Temp
		}
		tempf.IsMin = matches[1] == "N"
		tempf.IsMax = matches[1] == "X"
		tempf.DateTime, _ = time.Parse("2006010215Z", CurYearStr+CurMonthStr+matches[4])
		// if date in next month
		if tempf.DateTime.Day() < t.DateTime.Day() {
			tempf.DateTime = tempf.DateTime.AddDate(0, 1, 0)
		}
		t.Temperature = append(t.Temperature, *tempf)
		return true
	}
	return false
}

func (t *TAFMessage) setTimeRange(fromStr, toStr string) {

	t.ValidFrom, _ = time.Parse("2006010215", CurYearStr+CurMonthStr+fromStr)
	// hours maybe 24
	if toStr[2:] == "24" {
		t.ValidTo, _ = time.Parse("2006010215", CurYearStr+CurMonthStr+toStr[:2]+"23")
		t.ValidTo = t.ValidTo.Add(time.Hour)
	} else {
		t.ValidTo, _ = time.Parse("2006010215", CurYearStr+CurMonthStr+toStr)
	}

	//	forecast for next month
	if t.ValidFrom.Day() < t.DateTime.Day() {
		t.ValidFrom = t.ValidFrom.AddDate(0, 1, 0)
	}
	if t.ValidTo.Day() < t.DateTime.Day() {
		t.ValidTo = t.ValidTo.AddDate(0, 1, 0)
	}
}

func (t *TAFMessage) appendPhenomena(input string) bool {

	if p := ph.ParsePhenomena(input); p != nil {
		t.Phenomena = append(t.Phenomena, *p)
		return true
	}
	return false
}

func (t *TAFMessage) appendCloud(input string) bool {

	if cl, ok := clouds.ParseCloud(input); ok {
		t.Clouds = append(t.Clouds, cl)
		return true
	}
	return false
}

func (t *TAFMessage) setVerticalVisibility(input string) bool {

	regex := regexp.MustCompile(`VV(\d{3})`)
	matches := regex.FindStringSubmatch(input)
	if len(matches) != 0 && matches[1] != "" {
		t.VerticalVisibility, _ = strconv.Atoi(matches[1])
		t.VerticalVisibility *= 100
		return true
	}
	return false
}

func (t *TAFMessage) findTrendsInMessage(tokens []string, startposition int) (endposition int) {

	var trends [][]string
	endposition = len(tokens)
	for i := len(tokens) - 1; i > startposition; i-- {
		if tokens[i] == TEMPO || tokens[i] == BECMG || strings.HasPrefix(tokens[i], "PROB") || strings.HasPrefix(tokens[i], "FM") {
			if strings.HasPrefix(tokens[i-1], "PROB") {
				i--
			}
			trends = append([][]string{tokens[i:endposition]}, trends[0:]...)
			endposition = i
		}
	}
	for _, trendstr := range trends {
		if trend := parseTrendData(trendstr); trend != nil {
			t.TREND = append(t.TREND, *trend)
		}
	}
	return

}
