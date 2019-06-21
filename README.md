[![Build Status](https://travis-ci.org/urkk/metar.svg?branch=master)](https://travis-ci.org/urkk/metar)
[![Coverage](https://codecov.io/gh/urkk/metar/branch/master/graph/badge.svg)](https://codecov.io/gh/urkk/metar)
# METAR
METAR (METeorological Aerodrome Report) and TAF (terminal aerodrome forecast) message decoder for use in bots, templates and other data visualization. Raw text messages (as is) are used for decoding.

Based on the format approved by The Federal Service for Hydrometeorology and Environmental Monitoring of Russia.

## Now supported
* In metars
	* Header: station location, date/time, auto/cor/nil
	* Wind (Wind Variability)
	* Visibility
	* Runway Visual Range
	* Type of Weather
	* Clouds
	* Temperature/Dewpoint
	* Altimeter Setting
* Metar supplementary informaton
	* Type of recent weather
	* State of the runway(s) (as *R24/010060*)
	* Wind shear on runway(s) (as *WS R24* or *WS ALL RWY*)
* In tafs
	* Header: station location, date/time, cor/amd/nil/cnl
	* Wind
	* Visibility
	* Type of Weather
	* Clouds
	* Temperatures
* Both - change expected: wind, visibility, type of Weather, clouds

### Supported units
* Wind speed: knots, meters per second or kilometer per hour
* Horizontal visibility: meters or american land miles
* Runway visual range: meters or feet
* QNH pressure: hectopascal or inch of mercury

### Limitations
* no color codes decoded
* only russian style remarks (no RMK from auto station, etc)

## Example
```go
    import (
        "fmt"

        "github.com/urkk/metar"
    )

    msg, err := metar.NewMETAR("URSS 220630Z 02003MPS 9999 -SHRA SCT050CB OVC086 20/16 Q1015 R02/290060 R06/290060 TEMPO -TSRA BKN030CB RMK R06/03002MPS QFE760")
    if err == nil {
        fmt.Printf("%+v\n", msg)
    }
```

# Links

* [Instructional material on metar, speci, taf codes](http://metavia2.ru/help/instruction_METAR_SPECI_TAF.pdf)

 
