package metar

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/urkk/metar/wind"
)

// Remark - Additional information not included in the main message
type Remark struct {
	WindOnRWY []WindOnRWY
	QBB       int  // Cloud base in meters
	МТOBSC    bool // Mountains obscured
	MASTOBSC  bool // Mast obscured
	OBSTOBSC  bool // Obstacle obscured
	QFE       int  // Q-code Field Elevation (mmHg)
}

// WindOnRWY - surface wind observations on the runways
type WindOnRWY struct {
	Runway string
	Wind   wind.Wind
}

func parseRemarks(tokens []string) *Remark {

	RMK := new(Remark)
	var count = 0
	for count < len(tokens) {
		// Wind value on runway. Not documented, but used in URSS and UHMA
		regex := regexp.MustCompile(`^(R\d{2}[LCR]?)/((\d{3})?(VRB)?(P)?(\d{2})?(G\d\d)?(MPS|KT))`)
		matches := regex.FindStringSubmatch(tokens[count])
		if len(matches) != 0 && matches[0] != "" {
			wnd := new(WindOnRWY)
			wnd.Runway = matches[1][1:]
			input := matches[2]
			if count < len(tokens)-1 {
				input += tokens[count+1]
			}
			wind, countused := wind.ParseWind(input)
			wnd.Wind = wind
			RMK.WindOnRWY = append(RMK.WindOnRWY, *wnd)
			count += countused
		}

		if count < len(tokens) && strings.HasPrefix(tokens[count], "QBB") {
			RMK.QBB, _ = strconv.Atoi(tokens[count][3:])
			count++
		}

		for count < len(tokens)-1 && tokens[count+1] == "OBSC" {
			switch tokens[count] {
			case "MT":
				RMK.МТOBSC = true
			case "MAST":
				RMK.MASTOBSC = true
			case "OBST":
				RMK.OBSTOBSC = true
			}
			count += 2
		}
		// may be QFE767/1022 (mmHg/hPa)
		if count < len(tokens) && strings.HasPrefix(tokens[count], "QFE") {
			RMK.QFE, _ = strconv.Atoi(tokens[count][3:6])
			count++
		}
		count++
	}
	return RMK
}
