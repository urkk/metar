package conversion

import (
	"math"
)

// KphToMps - converts kilometres per hour to meters per second
func KphToMps(kph int) float64 {
	return float64(kph) / 3.6
}

// KphToKts - converts kilometres per hour to knots
func KphToKts(kph int) float64 {
	return float64(kph) / 1.852
}

// KtsToMps - converts knots to meters per second
func KtsToMps(kts float64) float64 {
	return kts / 1.94384
}

// MpsToKts - converts meters per second to knots
func MpsToKts(m float64) float64 {
	return m * 1.94384
}

// SMileToM - converts statute miles to meters
func SMileToM(sm float64) int {
	return int(math.Round(sm * 1609.344))
}

// FtToM - converts feet to meters (rounded to 10 meters)
func FtToM(ft int) int {
	return int(math.Round(float64(ft)*0.3048/10) * 10)
}

// MToFt - converts metres to feet (rounded to 10)
func MToFt(m int) int {
	return int(math.Round(float64(m)*3.28084/10) * 10)
}

// MToSMile - converts metres to statute miles
func MToSMile(m int) float64 {
	return float64(m) * 0.00062137119223733
}

// FtToSMile - converts feet to statute miles
func FtToSMile(m int) float64 {
	return float64(m) / 5280
}

// SMileToFt - converts statute miles to feet
func SMileToFt(m float64) int {
	return int(m * 5280)
}

// InHgTohPa - converts inch of mercury to hectopascal
func InHgTohPa(inHg float64) int {
	return int(math.Round(inHg * 33.86389))
}

// HPaToMmHg - converts hectopascal to mm of mercury
func HPaToMmHg(hPa int) int {
	return int(math.Round(float64(hPa) * 0.75006375541921))
}

// MmHgToHPa - converts mm of mercury to hectopascal
func MmHgToHPa(mm int) int {
	return int(math.Round(float64(mm) * 1.333223684))
}

// DirectionToCardinalDirection - converts direction in degrees to points of the compass
func DirectionToCardinalDirection(dir int) string {
	index := int(math.Round(float64(dir%360) / 45))
	return map[int]string{0: "N", 1: "NE", 2: "E", 3: "SE", 4: "S", 5: "SW", 6: "W", 7: "NW", 8: "N"}[index]
}

// CalcRelativeHumidity - calculates the relative humidity of the dew point and temperature
func CalcRelativeHumidity(temp, dewpoint int) int {
	// see https://www.vaisala.com/sites/default/files/documents/Humidity_Conversion_Formulas_B210973EN-F.pdf
	// used constants in temperature range -20...+50°C
	m := 7.591386
	tn := 240.7263
	rh := 100 * math.Pow(10, m*(float64(dewpoint)/(float64(dewpoint)+tn)-(float64(temp)/(float64(temp)+tn))))
	return int(math.Round(rh))

}
