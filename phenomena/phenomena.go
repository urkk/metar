package phenomena

import (
	"regexp"
)

// Possible combinations of codes and modifiers. Not all codes can be combined. So the VC can only be applied to SH, TS, FG, VA, BLDU, BLSA, BLSN, PO, FC, SS, DS. It's impossible to have a "light" tornado (-FC), ground only fog (MIFG), etc. See documenation.
var PossiblePhemonenaCode = map[string]bool{"-DZ": true, "-RA": true, "-SN": true, "-SG": true, "-PL": true, "UP": true, "-DZRA": true, "-RADZ": true, "-SNDZ": true, "-SGDZ": true, "-PLDZ": true, "-DZSN": true, "-RASN": true, "-SNRA": true, "-SGRA": true, "-PLRA": true, "FZUP": true, "-DZSG": true, "-RASG": true, "-SNSG": true, "-SGSN": true, "-PLSN": true, "-DZPL": true, "-RAPL": true, "-SNPL": true, "-SGPL": true, "-PLSG": true, "-DZRASN": true, "-RADZSN": true, "-SNDZRA": true, "-SGDZRA": true, "-PLDZRA": true, "-DZRASG": true, "-RADZSG": true, "-SNRADZ": true, "-SGRASN": true, "-PLRASN": true, "-DZRAPL": true, "-RADZPL": true, "-SNRASG": true, "-SGPLSN": true, "-PLSNRA": true, "-DZSNRA": true, "-RASNDZ": true, "-SNRAPL": true, "-SGSNRA": true, "-PLRADZ": true, "-DZSGRA": true, "-RASNSG": true, "-SNPLRA": true, "-SGRADZ": true, "-PLSNSG": true, "-DZPLRA": true, "-RASNPL": true, "-SNPLSG": true, "-SGSNPL": true, "-PLSGSN": true, "-RASGSN": true, "-SNSGRA": true, "-RASGDZ": true, "-SNSGPL": true, "PL": true, "-RAPLDZ": true, "PLDZ": true, "-RAPLSN": true, "PLRA": true, "DZ": true, "RA": true, "SN": true, "SG": true, "PLSN": true, "DZRA": true, "RADZ": true, "SNDZ": true, "SGDZ": true, "PLSG": true, "DZSN": true, "RASN": true, "SNRA": true, "SGRA": true, "PLDZRA": true, "DZSG": true, "RASG": true, "SNSG": true, "SGSN": true, "PLRASN": true, "DZPL": true, "RAPL": true, "SNPL": true, "SGPL": true, "PLSNRA": true, "DZRASN": true, "RADZSN": true, "SNDZRA": true, "SGDZRA": true, "PLRADZ": true, "DZRASG": true, "RADZSG": true, "SNRADZ": true, "SGRASN": true, "PLSNSG": true, "DZRAPL": true, "RADZPL": true, "SNRASG": true, "SGPLSN": true, "PLSGSN": true, "DZSNRA": true, "RASNDZ": true, "SNRAPL": true, "SGSNRA": true, "+PL": true, "DZSGRA": true, "RASNSG": true, "SNPLRA": true, "SGRADZ": true, "+PLDZ": true, "DZPLRA": true, "RASNPL": true, "SNPLSG": true, "SGSNPL": true, "+PLRA": true, "RASGSN": true, "SNSGRA": true, "+PLSN": true, "RASGDZ": true, "SNSGPL": true, "+PLSG": true, "RAPLDZ": true, "+PLDZRA": true, "RAPLSN": true, "+PLRASN": true, "+DZ": true, "+RA": true, "+SN": true, "+SG": true, "+PLSNRA": true, "+DZRA": true, "+RADZ": true, "+SNDZ": true, "+SGDZ": true, "+PLRADZ": true, "+DZSN": true, "+RASN": true, "+SNRA": true, "+SGRA": true, "+PLSNSG": true, "+DZSG": true, "+RASG": true, "+SNSG": true, "+SGSN": true, "+PLSGSN": true, "+DZPL": true, "+RAPL": true, "+SNPL": true, "+SGPL": true, "SHUP": true, "+DZRASN": true, "+RADZSN": true, "+SNDZRA": true, "+SGDZRA": true, "TSUP": true, "+DZRASG": true, "+RADZSG": true, "+SNRADZ": true, "+SGRASN": true, "TS": true, "+DZRAPL": true, "+RADZPL": true, "+SNRASG": true, "+SGPLSN": true, "VCTS": true, "+DZSNRA": true, "+RASNDZ": true, "+SNRAPL": true, "+SGSNRA": true, "+DZSGRA": true, "+RASNSG": true, "+SNPLRA": true, "+SGRADZ": true, "+DZPLRA": true, "+RASNPL": true, "+SNPLSG": true, "+SGSNPL": true, "+RASGSN": true, "+SNSGRA": true, "+RASGDZ": true, "+SNSGPL": true, "+RAPLDZ": true, "+RAPLSN": true, "-SHRA": true, "-SHSN": true, "-SHGR": true, "-SHGS": true, "-SHRASN": true, "-SHSNRA": true, "-SHGRRA": true, "-SHGSRA": true, "-SHRAGR": true, "-SHSNGR": true, "-SHGRSN": true, "-SHGSSN": true, "-SHRAGS": true, "-SHSNGS": true, "-SHRASNGR": true, "-SHSNRAGR": true, "-SHGRRASN": true, "-SHGSRASN": true, "-SHRAGRSN": true, "-SHSNGRRA": true, "-SHGRSNRA": true, "-SHGSSNRA": true, "-SHRASNGS": true, "-SHSNRAGS": true, "-SHRAGSSN": true, "-SHSNGSRA": true, "SHRA": true, "SHSN": true, "SHGR": true, "SHGS": true, "SHRASN": true, "SHSNRA": true, "SHGRRA": true, "SHGSRA": true, "SHRAGR": true, "SHSNGR": true, "SHGRSN": true, "SHGSSN": true, "SHRAGS": true, "SHSNGS": true, "SHRASNGR": true, "SHSNRAGR": true, "SHGRRASN": true, "SHGSRASN": true, "SHRAGRSN": true, "SHSNGRRA": true, "SHGRSNRA": true, "SHGSSNRA": true, "SHRASNGS": true, "SHSNRAGS": true, "SHRAGSSN": true, "SHSNGSRA": true, "+SHRA": true, "+SHSN": true, "+SHGR": true, "+SHGS": true, "+SHRASN": true, "+SHSNRA": true, "+SHGRRA": true, "+SHGSRA": true, "+SHRAGR": true, "+SHSNGR": true, "+SHGRSN": true, "+SHGSSN": true, "+SHRAGS": true, "+SHSNGS": true, "+SHRASNGR": true, "+SHSNRAGR": true, "+SHGRRASN": true, "+SHGSRASN": true, "+SHRAGRSN": true, "+SHSNGRRA": true, "+SHGRSNRA": true, "+SHGSSNRA": true, "+SHRASNGS": true, "+SHSNRAGS": true, "+SHRAGSSN": true, "+SHSNGSRA": true, "-TSRA": true, "-TSSN": true, "-TSGR": true, "-TSGS": true, "-TSRASN": true, "-TSSNRA": true, "-TSGRRA": true, "-TSGSRA": true, "-TSRAGR": true, "-TSSNGR": true, "-TSGRSN": true, "-TSGSSN": true, "-TSRAGS": true, "-TSSNGS": true, "-TSRASNGR": true, "-TSSNRAGR": true, "-TSGRRASN": true, "-TSGSRASN": true, "-TSRAGRSN": true, "-TSSNGRRA": true, "-TSGRSNRA": true, "-TSGSSNRA": true, "-TSRASNGS": true, "-TSSNRAGS": true, "-TSRAGSSN": true, "-TSSNGSRA": true, "TSRA": true, "TSSN": true, "TSGR": true, "TSGS": true, "TSRASN": true, "TSSNRA": true, "TSGRRA": true, "TSGSRA": true, "TSRAGR": true, "TSSNGR": true, "TSGRSN": true, "TSGSSN": true, "TSRAGS": true, "TSSNGS": true, "TSRASNGR": true, "TSSNRAGR": true, "TSGRRASN": true, "TSGSRASN": true, "TSRAGRSN": true, "TSSNGRRA": true, "TSGRSNRA": true, "TSGSSNRA": true, "TSRASNGS": true, "TSSNRAGS": true, "TSRAGSSN": true, "TSSNGSRA": true, "+TSRA": true, "+TSSN": true, "+TSGR": true, "+TSGS": true, "+TSRASN": true, "+TSSNRA": true, "+TSGRRA": true, "+TSGSRA": true, "+TSRAGR": true, "+TSSNGR": true, "+TSGRSN": true, "+TSGSSN": true, "+TSRAGS": true, "+TSSNGS": true, "+TSRASNGR": true, "+TSSNRAGR": true, "+TSGRRASN": true, "+TSGSRASN": true, "+TSRAGRSN": true, "+TSSNGRRA": true, "+TSGRSNRA": true, "+TSGSSNRA": true, "+TSRASNGS": true, "+TSSNRAGS": true, "+TSRAGSSN": true, "+TSSNGSRA": true, "-FZDZ": true, "-FZRA": true, "-FZDZRA": true, "-FZRADZ": true, "FZDZ": true, "FZRA": true, "FZDZRA": true, "FZRADZ": true, "+FZDZ": true, "+FZRA": true, "+FZDZRA": true, "+FZRADZ": true, "-DS": true, "DS": true, "+DS": true, "VCDS": true, "-SS": true, "SS": true, "+SS": true, "VCSS": true, "FG": true, "FC": true, "PO": true, "VA": true, "VCFG": true, "+FC": true, "VCPO": true, "VCVA": true, "VCFC": true, "VCSH": true, "BLSA": true, "BLDU": true, "BLSN": true, "DRSA": true, "DRDU": true, "DRSN": true, "SA": true, "DU": true, "VCBLSA": true, "VCBLDU": true, "VCBLSN": true, "MIFG": true, "PRFG": true, "BCFG": true, "FZFG": true, "BR": true, "HZ": true, "FU": true, "SQ": true, "NSW": true}

// Not all phenomena are acceptable as recent. No indication of intensity.
var PossibleRecentPhemonenaCode = map[string]bool{"REFZDZ": true, "REFZRA": true, "REDZ": true, "RERA": true, "RESHRA": true, "RERASN": true, "RESN": true, "RESG": true, "RESHGR": true, "RESHGS": true, "REBLSN": true, "RESS": true, "REDS": true, "RETSRA": true, "RETSSN": true, "RETSGR": true, "RETSGS": true, "RETS": true, "REFC": true, "REVA": true, "REPL": true, "REUP": true, "REFZUP": true, "RETSUP": true, "RESHUP": true}

// A list of unique codes to convert groups such as TSRASNGR to Thunderstorm with rain (TSRA), snow (SN), hail (GR)
var UniqueCodes = map[string]bool{"BCFG": true, "BLDU": true, "BLSA": true, "BLSN": true, "BR": true, "DRDU": true, "DRSA": true, "DRSN": true, "DS": true, "DU": true, "DZ": true, "FC": true, "FG": true, "FU": true, "FZDZ": true, "FZFG": true, "FZRA": true, "GS": true, "GR": true, "HZ": true, "IC": true, "MIFG": true, "PO": true, "PL": true, "PRFG": true, "RA": true, "SA": true, "SG": true, "SH": true, "SHGR": true, "SHGS": true, "SHRA": true, "SHSN": true, "SN": true, "SQ": true, "SS": true, "TS": true, "TSGR": true, "TSGS": true, "TSRA": true, "TSSN": true, "UP": true, "SHUP": true, "TSUP": true, "FZUP": true, "VA": true, "NSW": true}

type Intensity string

const (
	Moderate Intensity = ""
	Light              = "-"
	Heavy              = "+"
)

type Phenomena struct {
	Vicinity     bool
	Intensity    Intensity
	Abbreviation string
}

func ParsePhenomena(token string) (p *Phenomena) {
	if _, found := PossiblePhemonenaCode[token]; found {
		p = new(Phenomena)
		regex := regexp.MustCompile(`(\+|-)?(VC)?([A-Z]{2,8})`)
		matches := regex.FindStringSubmatch(token)
		p.Intensity = Intensity(matches[1])
		p.Vicinity = matches[2] == "VC"
		p.Abbreviation = matches[3]
	}
	return p
}

// Recent phenomena can't be in the vicinity and have qualifier
func ParseRecentPhenomena(token string) (p *Phenomena) {
	if _, found := PossibleRecentPhemonenaCode[token]; found {
		p = new(Phenomena)
		// remove RE
		p.Abbreviation = token[2:]
	}
	return p
}
