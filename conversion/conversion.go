package conversion

import (
	"math"
)

// converts kilometres per hour to meters per second
func KphToMps(kph float64) float64 {
	return kph / 3.6
}

// converts knots to meters per second
func KtsToMps(kts float64) float64 {
	return kts / 1.944
}

// converts meters per second to knots
func MpsToKts(m float64) float64 {
	return m * 1.94384
}

// converts statute miles to meters
func SMileToM(sm int) int {
	return int(math.Round(float64(sm) * 1609.344))
}

// converts feet to meters (rounded to 10 meters)
func FtToM(ft int) int {
	return int(math.Round(float64(ft)*0.3048/10) * 10)
}

// converts inch of mercury to hectopascal
func InHgTohPa(inHg float64) int {
	return int(math.Round(inHg * 33.86389))
}

// converts hectopascal to mm of mercury
func HPaToMmHg(hPa int) int {
	return int(math.Round(float64(hPa) * 0.75006375541921))
}

// converts mm of mercury to hectopascal
func MmHgToHPa(mm int) int {
	return int(math.Round(float64(mm) * 1.333223684))
}

// converts direction in degrees to points of the compass
func DirectionToCardinalDirection(dir int) string {
	index := int(math.Round(float64(dir%360) / 45))
	return map[int]string{0: "N", 1: "NE", 2: "E", 3: "SE", 4: "S", 5: "SW", 6: "W", 7: "NW", 8: "N"}[index]
}

// calculates the relative humidity of the dew point and temperature
func CalcRelativeHumidity(temp, dewpoint int) int {
	// see https://www.vaisala.com/sites/default/files/documents/Humidity_Conversion_Formulas_B210973EN-F.pdf
	// used constants in temperature range -20...+50°C
	m := 7.591386
	tn := 240.7263
	rh := 100 * math.Pow(10, m*(float64(dewpoint)/(float64(dewpoint)+tn)-(float64(temp)/(float64(temp)+tn))))
	return int(math.Round(rh))

}
