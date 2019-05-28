// Package metar provides METAR (METeorological Aerodrome Report) message decoding
package metar

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/urkk/metar/clouds"
	cnv "github.com/urkk/metar/conversion"
	ph "github.com/urkk/metar/phenomena"
	rwy "github.com/urkk/metar/runways"
	"github.com/urkk/metar/wind"
)

// Year, month and day. By default read all messages in the current date. Can be redefined if necessary
var CurYearStr, CurMonthStr, CurDayStr string

func init() {
	now := time.Now()
	CurYearStr = strconv.Itoa(now.Year())
	CurMonthStr = fmt.Sprintf("%02d", now.Month())
	CurDayStr = fmt.Sprintf("%02d", now.Day())
}

// MetarMessage - Meteorological report presented as a data structure
type MetarMessage struct {
	rawData  string    // The raw METAR
	COR      bool      // Correction to observation
	Station  string    // 4-letter ICAO station identifier
	DateTime time.Time // Time (in ISO8601 date/time format) this METAR was observed
	Auto     bool      // METAR from automatic observing systems with no human intervention
	NIL      bool      // event of missing METAR
	Wind     wind.Wind //	Surface wind
	// Ceiling And Visibility OK, indicating no cloud below 5,000 ft (1,500 m) or the highest minimum sector
	// altitude and no cumulonimbus or towering cumulus at any level, a visibility of 10 km (6 mi) or more and no significant weather change.
	CAVOK                        bool
	Visibility                   Visibility        // Horizontal visibility. In meters
	RWYvisibility                []rwy.VisualRange // Runway visual range
	Phenomena                    []ph.Phenomena    // Present Weather
	PhenomenaNotDefined          bool              // Not detected by the automatic station - “//”
	VerticalVisibility           int               // Vertical visibility (ft)
	VerticalVisibilityNotDefined bool              // “///”
	Clouds                       []clouds.Cloud    // Cloud amount and height
	Temperature                  int               // Temperature in degrees Celsius
	Dewpoint                     int               // Dew point in degrees Celsius
	QNHhPa                       int               // Altimeter setting.  Atmospheric pressure adjusted to mean sea level
	// Supplementary informaton
	//Recent weather
	RecentPhenomena []ph.Phenomena
	// Information on the state of the runway(s)
	RWYState []rwy.State
	// Wind shear on runway(s)
	WindShear []rwy.RunwayDesignator
	// Prevision
	TREND []Trend
	//OR NO SIGnificant changes coming within the next two hours
	NOSIG bool
	// Remarks consisting of recent operationally significant weather as well as additive and automated maintenance data
	Remarks *Remark
	// An array of tokens that couldn't be decoded
	NotDecodedTokens []string
}

// RAW - returns the original message text
func (m *MetarMessage) RAW() string { return m.rawData }

// NewMETAR - creates a new METAR based on the original message
func NewMETAR(inputtext string) *MetarMessage {

	m := &MetarMessage{
		rawData: inputtext,
	}
	headerRx := myRegexp{regexp.MustCompile(`^(?P<type>(METAR|SPECI)\s)?(?P<cor>COR\s)?(?P<station>\w{4})\s(?P<time>\d{6}Z)(?P<auto>\sAUTO)?(?P<nil>\sNIL)?`)}
	headermap := headerRx.FindStringSubmatchMap(m.rawData)
	m.Station = headermap["station"]
	m.DateTime, _ = time.Parse("200601021504Z", CurYearStr+CurMonthStr+headermap["time"])
	m.COR = headermap["cor"] != ""
	m.Auto = headermap["auto"] != ""
	m.NIL = headermap["nil"] != ""
	if m.Station == "" && m.DateTime.IsZero() {
		//not valid message?
		m.NotDecodedTokens = append(m.NotDecodedTokens, m.rawData)
		return m
	}
	if m.NIL {
		return m
	}
	tokens := strings.Split(m.rawData, " ")

	count := 0
	totalcount := len(tokens)
	for _, value := range headermap {
		if value != "" {
			count++
		}
	}

	var trends [][]string
	var remarks []string
	for i := totalcount - 1; i > count; i-- {
		if tokens[i] == "RMK" {
			remarks = append(remarks, tokens[i:totalcount]...)
			totalcount = i
		}
		if tokens[i] == TEMPO || tokens[i] == BECMG {
			//for correct order of following on reverse parsing append []trends to current trend
			trends = append([][]string{tokens[i:totalcount]}, trends[0:]...)
			totalcount = i
		}
	}

	for _, trendstr := range trends {
		if trend := parseTrendData(trendstr); trend != nil {
			m.TREND = append(m.TREND, *trend)
		}
	}

	if len(remarks) > 0 {
		m.Remarks = parseRemarks(remarks)
	}
	m.decodeMetar(tokens[count:totalcount])
	return m
}

// Visibility - prevailing visibility
type Visibility struct {
	Distance       int
	LowerDistance  int
	LowerDirection string
}

type myRegexp struct {
	*regexp.Regexp
}

func (r *myRegexp) FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)
	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}
	for i, name := range r.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}
		captures[name] = match[i]
	}
	return captures
}

func (m *MetarMessage) decodeMetar(tokens []string) {

	var regex *regexp.Regexp
	var matches []string
	totalcount := len(tokens)
	for count := 0; count < totalcount; {
		// Surface wind
		if wnd, tokensused := wind.ParseWind(strings.Join(tokens[count:], " ")); tokensused > 0 {
			m.Wind = wnd
			count += tokensused

		}
		if tokens[count] == "CAVOK" {
			m.CAVOK = true
			count++
		} else {
			// Horizontal visibility
			if vis, tokensused := ParseVisibility(strings.Join(tokens[count:], " ")); tokensused > 0 {
				m.Visibility = vis
				count += tokensused
			}
			// Runway visual range
			for count < totalcount {
				if RWYvis, ok := rwy.ParseVisibility(tokens[count]); ok {
					m.RWYvisibility = append(m.RWYvisibility, RWYvis)
					count++
				} else {
					break
				}
			}
			// Present Weather
			if tokens[count] == "//" {
				m.PhenomenaNotDefined = true
				count++
			}
			for count < totalcount {
				if p := ph.ParsePhenomena(tokens[count]); p != nil {
					m.Phenomena = append(m.Phenomena, *p)
					count++
				} else {
					break // the end of the weather group
				}
			}
			// Vertical visibility
			regex = regexp.MustCompile(`VV(\d{3}|///)`)
			matches = regex.FindStringSubmatch(tokens[count])
			if len(matches) != 0 && matches[0] != "" {
				count++
				if matches[1] != "///" {
					m.VerticalVisibility, _ = strconv.Atoi(matches[1])
					m.VerticalVisibility *= 100
				} else {
					m.VerticalVisibilityNotDefined = true
				}
			}
			// Cloudiness description
			for count < totalcount {
				if cl, ok := clouds.ParseCloud(tokens[count]); ok {
					m.Clouds = append(m.Clouds, cl)
					count++
				} else {
					break
				}
			}
		} //!CAVOK
		// Temperature and dew point
		if m.checkTemp(tokens[count]) {
			count++
		}
		// Altimeter setting
		if m.checkAltimetr(tokens[count]) {
			count++
		}
		//	All the following elements are optional
		// Recent weather
		for count < totalcount {
			if p := ph.ParseRecentPhenomena(tokens[count]); p != nil {
				m.RecentPhenomena = append(m.RecentPhenomena, *p)
				count++
			} else {
				break // the end of the weather group
			}
		}
		// Wind shear
		regex = regexp.MustCompile(`^WS\s((R\d{2}[LCR]?)|(ALL\sRWY))`)
		matches = regex.FindStringSubmatch(strings.Join(tokens[count:], " "))
		for ; len(matches) > 0; matches = regex.FindStringSubmatch(strings.Join(tokens[count:], " ")) {
			if matches[3] != "" { // WS ALL RWY
				rd := new(rwy.RunwayDesignator)
				rd.AllRunways = true
				m.WindShear = append(m.WindShear, *rd)
				count += 3
			}
			if matches[2] != "" { // WS R03
				m.WindShear = append(m.WindShear, rwy.NewRD(matches[1]))
				count += 2
			}
		}
		// TODO Sea surface condition
		// W19/S4  W15/Н7  W15/Н17 W15/Н175

		// State of the runway(s)
		if count < totalcount && tokens[count] == "R/SNOCLO" {
			rwc := new(rwy.State)
			rwc.SNOCLO = true
			m.RWYState = append(m.RWYState, *rwc)
			count++
		}
		for count < totalcount {
			if rwc, ok := rwy.ParseState(tokens[count]); ok {
				m.RWYState = append(m.RWYState, rwc)
				count++
			} else {
				break
			}
		}
		if count < totalcount && tokens[count] == "NOSIG" {
			m.NOSIG = true
			count++
		}
		// The token is not recognized or is located in the wrong position
		if count < totalcount {
			m.NotDecodedTokens = append(m.NotDecodedTokens, tokens[count])
			count++
		}
	} // End main section
}

// ParseVisibility - identify and parses the representation oh horizontal visibility
func ParseVisibility(token string) (v Visibility, tokensused int) {
	// The literal P (M) is not listed in the documentation, but is used in messages
	pattern := `^(P|M)?(\d{4})(\s|$)((\d{4})(NE|SE|NW|SW|N|E|S|W))?`
	if matched, _ := regexp.MatchString(pattern, token); !matched {
		return v, tokensused
	}
	tokensused = 1
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(token)
	v.Distance, _ = strconv.Atoi(matches[2])
	if matches[4] != "" {
		v.LowerDistance, _ = strconv.Atoi(matches[5])
		v.LowerDirection = matches[6]
		tokensused++
	}
	return v, tokensused
}

// Checks whether the string is a temperature and dew point values and writes this values
func (m *MetarMessage) checkTemp(input string) bool {
	regex := regexp.MustCompile(`^(M)?(\d{2})/(M)?(\d{2})$`)
	matches := regex.FindStringSubmatch(input)
	if len(matches) != 0 {
		m.Temperature, _ = strconv.Atoi(matches[2])
		m.Dewpoint, _ = strconv.Atoi(matches[4])
		if matches[1] == "M" {
			m.Temperature = -m.Temperature
		}
		if matches[3] == "M" {
			m.Dewpoint = -m.Dewpoint
		}
		return true
	}
	return false
}

func (m *MetarMessage) checkAltimetr(input string) bool {
	regex := regexp.MustCompile(`([Q|A])(\d{4})`)
	matches := regex.FindStringSubmatch(input)
	if len(matches) != 0 {
		if matches[1] == "A" {
			inHg, _ := strconv.ParseFloat(matches[2][:2]+"."+matches[2][2:4], 64)
			m.QNHhPa = int(cnv.InHgTohPa(inHg))
		} else {
			m.QNHhPa, _ = strconv.Atoi(matches[2])
		}
		return true
	}
	return false
}
