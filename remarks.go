package metar

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/urkk/metar/wind"
)

type Remark struct {
	WindOnRWY []WindOnRWY
	QBB       int  // Cloud base in meters
	МТOBSC    bool // Mountains obscured
	MASTOBSC  bool // Mast obscured
	OBSTOBSC  bool // Obstacle obscured
	QFE       int  // Q-code Field Elevation (mmHg)
}

func parseRemarks(tokens []string) *Remark {

	RMK := new(Remark)
	var count = 0
	for count < len(tokens) {
		// Wind value on runway. Not documented, used in URSS
		regex := regexp.MustCompile(`(R\d{2}[LCR]?)/((\d{3})?(VRB)?(P)?(\d{2})?(G\d\d)?(MPS|KT))`)
		matches := regex.FindStringSubmatch(tokens[count])
		if len(matches) != 0 && matches[0] != "" {
			wnd := new(WindOnRWY)
			wnd.Runway = matches[1][1:]
			wind, _, multiple := wind.ParseWind(matches[2] + tokens[count+1])
			if multiple {
				count++
			}
			wnd.Wind = wind
			RMK.WindOnRWY = append(RMK.WindOnRWY, *wnd)
			count++
		}
		if count >= len(tokens) {
			break
		}
		if strings.HasPrefix(tokens[count], "QBB") {
			RMK.QBB, _ = strconv.Atoi(tokens[count][3:])
			count++
		}
		if count >= len(tokens) {
			break
		}
		if tokens[count] == "MT" && tokens[count+1] == "OBSC" {
			RMK.МТOBSC = true
			count += 2
		}
		if count >= len(tokens) {
			break
		}
		if tokens[count] == "MAST" && tokens[count+1] == "OBSC" {
			RMK.MASTOBSC = true
			count += 2
		}
		if count >= len(tokens) {
			break
		}
		if tokens[count] == "OBST" && tokens[count+1] == "OBSC" {
			RMK.OBSTOBSC = true
			count += 2
		}
		if count >= len(tokens) {
			break
		}
		// may be QFE767/1022 (mmHg/hPa)
		if strings.HasPrefix(tokens[count], "QFE") {
			RMK.QFE, _ = strconv.Atoi(tokens[count][3:6])
			count++
		}
		count++
	}
	return RMK
}
