[![Build Status](https://travis-ci.org/urkk/metar.svg?branch=master)](https://travis-ci.org/urkk/metar)
[![Coverage](https://codecov.io/gh/urkk/metar/branch/master/graph/badge.svg)](https://codecov.io/gh/urkk/metar)
# METAR
METAR (METeorological Aerodrome Report) and TAF (terminal aerodrome forecast) message decoder for use in bots, templates and other data visualization.

Suitable for the recognition messages from russians airports, and in the ex-USSR and European countries with some limitations. Based on the format approved by The Federal Service for Hydrometeorology and Environmental Monitoring of Russia.

### Limitations
* visibility only in meters, not miles.
* no color codes decoded
* only russians remarks (no RMK from auto station, etc)

## Example
```go
    import (
        "fmt"

        "github.com/urkk/metar"
    )

    msg := metar.NewMETAR("URSS 220630Z 02003MPS 9999 -SHRA SCT050CB OVC086 20/16 Q1015 R02/290060 R06/290060 TEMPO -TSRA BKN030CB RMK R06/03002MPS QFE760")
    fmt.Printf("%+v\n", msg)
```

# Links

* [Instructional material on metar, speci, taf codes](http://metavia2.ru/help/instruction_METAR_SPECI_TAF.pdf)

 