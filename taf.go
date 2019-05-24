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

type TemperatureForecast struct {
	Temp     int
	DateTime time.Time
	IsMax    bool
	IsMin    bool
}

// Terminal Aerodrome Forecast
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
	VerticalVisibility int                   // Vertical visibility
	Clouds             []clouds.Cloud        // Cloud amount and height
	Temperature        []TemperatureForecast // Temperature extremes
	// Prevision
	TREND []Trend

	NotDecodedTokens []string
}

func NewTAF(inputtext string) *TAFMessage {
	t := &TAFMessage{
		rawData: inputtext,
	}
	t.decodeTAF()
	return t
}
func (t *TAFMessage) RAW() string { return t.rawData }

func (t *TAFMessage) decodeTAF() {

	headerRx := myRegexp{regexp.MustCompile(`^(?P<taf>TAF\s)?(?P<cor>COR\s)?(?P<amd>AMD\s)?(?P<station>\w{4})\s(?P<time>\d{6}Z)(?P<nil>\sNIL)?(\s(?P<from>\d{4})/(?P<to>\d{4}))?(?P<cnl>\sCNL)?`)}
	headermap := headerRx.FindStringSubmatchMap(t.RAW())
	t.Station = headermap["station"]
	t.DateTime, _ = time.Parse("200601021504Z", CurYearStr+CurMonthStr+headermap["time"])
	t.COR = headermap["cor"] != ""
	t.AMD = headermap["amd"] != ""
	t.NIL = headermap["nil"] != ""
	t.CNL = headermap["cnl"] != ""

	if t.NIL { // End of TAF, if the forecast is lost
		return
	}
	t.ValidFrom, _ = time.Parse("2006010215", CurYearStr+CurMonthStr+headermap["from"])
	// hours maybe 24
	if headermap["to"][2:] == "24" {
		t.ValidTo, _ = time.Parse("2006010215", CurYearStr+CurMonthStr+headermap["to"][:2]+"23")
		t.ValidTo = t.ValidTo.Add(time.Hour)
	} else {
		t.ValidTo, _ = time.Parse("2006010215", CurYearStr+CurMonthStr+headermap["to"])
	}

	//	forecast for next month
	if t.ValidFrom.Day() < t.DateTime.Day() {
		t.ValidFrom = t.ValidFrom.AddDate(0, 1, 0)
	}
	if t.ValidTo.Day() < t.DateTime.Day() {
		t.ValidTo = t.ValidTo.AddDate(0, 1, 0)
	}
	if t.CNL { // End of TAF, if the forecast is cancelled
		return
	}
	tokens := strings.Split(t.RAW(), " ")

	count := 0
	for _, value := range headermap {
		if value != "" {
			count++
		}
	}
	// field "from" and "to" - it's one token (DDhh/DDhh), and they are mandatory. step back.
	if count > 3 { // station id, issued, valid from, valid to
		count--
	}

	var trends [][]string

	totalcount := len(tokens)
	for i := len(tokens) - 1; i > count; i-- {
		if tokens[i] == TEMPO || tokens[i] == BECMG || strings.HasPrefix(tokens[i], "PROB") || strings.HasPrefix(tokens[i], "FM") {
			if strings.HasPrefix(tokens[i-1], "PROB") {
				i--
			}
			trends = append([][]string{tokens[i:totalcount]}, trends[0:]...)
			totalcount = i
		}
	}

	for _, trend := range trends {
		trend := parseTrendData(trend)
		if trend != nil {
			t.TREND = append(t.TREND, *trend)
		}
	}

	for count < totalcount {
		// Wind - Visibility - Weather - Sky Condition

		// Surface wind
		if wnd, ok, _ := wind.ParseWind(tokens[count]); ok {
			t.Wind = wnd
			count++
		}
		if tokens[count] == "CAVOK" {
			t.CAVOK = true
			count++
		}
		var regex *regexp.Regexp
		var matches []string
		if !t.CAVOK {
			// Horizontal visibility
			if vis, ok, _ := ParseVisibility(tokens[count]); ok {
				t.Visibility = vis
				count++
			}
			// Weather
			for i := count; i < len(tokens); i++ {
				p := ph.ParsePhenomena(tokens[count])
				if p != nil {
					t.Phenomena = append(t.Phenomena, *p)
					count++
				} else {
					break // the end of the weather group
				}
			}
			// Vertical visibility
			regex = regexp.MustCompile(`VV(\d{3})`)
			matches = regex.FindStringSubmatch(tokens[count])
			if len(matches) != 0 && matches[0] != "" {
				count++
				if matches[1] != "" {
					t.VerticalVisibility, _ = strconv.Atoi(matches[1])
				}
			}
			// Cloudiness description
			for i := count; i < len(tokens); i++ {
				if cl, ok := clouds.ParseCloud(tokens[count]); ok {
					t.Clouds = append(t.Clouds, cl)
					count++
				} else {
					break
				}
			}
		} // !CAVOK
		if count >= len(tokens) {
			break
		}
		// Temperature
		regex = regexp.MustCompile(`^T(X|N)(M)?(\d\d)\/(\d{4}Z)`)
		matches = regex.FindStringSubmatch(tokens[count])
		for ; len(matches) > 0; matches = regex.FindStringSubmatch(tokens[count]) {
			tempf := new(TemperatureForecast)
			tempf.Temp, _ = strconv.Atoi(matches[3])
			if matches[2] == "M" {
				tempf.Temp = -tempf.Temp
			}
			if matches[1] == "N" {
				tempf.IsMin = true
			} else {
				tempf.IsMax = true
			}
			tempf.DateTime, _ = time.Parse("2006010215Z", time.Now().Format("200601")+matches[4])
			t.Temperature = append(t.Temperature, *tempf)
			count++
			if count >= len(tokens) {
				break
			}
		}
		// The token is not recognized or is located in the wrong position
		if count < totalcount {
			t.NotDecodedTokens = append(t.NotDecodedTokens, tokens[count])
			count++
		}
	} // End main section
}
