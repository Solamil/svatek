package svatek 

import (
	"fmt"
	"time"
)

// credits to https://kalendar.beda.cz/vypocet-velikonocni-nedele-v-ruznych-programovacich-jazycich
func Velikonoce(rok int) string {
	if rok <= 1583 {
		return ""
	}
	zlateCislo := (rok % 19) + 1
	julEpakta := (11 * zlateCislo) % 30
	stoleti := int(rok / 100) + 1
	slunecniOprava := int(3 * (stoleti - 16) / 4)
	mesicniOprava := int(8 * (stoleti - 15) / 25)
	epakta := (julEpakta - 10 - slunecniOprava + mesicniOprava) % 30
	if epakta < 0 {
		epakta += 30
	}
	tmp := epakta
	if epakta == 24 || (epakta == 25 && zlateCislo > 11) {
		tmp += 1
	}
	pfm := 0 // Paschal Full Moon
	if tmp < 24 {
		pfm = 44 - tmp
	} else {
		pfm = 74 - tmp
	}

	gregOprava := 10 + slunecniOprava
	denTydnePfm := (rok + (int)(rok / 4) - gregOprava + pfm) % 7
	if denTydnePfm < 0 {
		denTydnePfm += 7
	}
	velNedele := pfm + 7 - denTydnePfm
	if velNedele < 32 {
		return fmt.Sprintf("%d. březen", velNedele)	
	} 
	return fmt.Sprintf("%d. duben", velNedele-31)	
}


func Denmatek(year int) time.Time {
	// Second sunday at the month of May
	t := time.Date(year, time.May, 1, 0, 0, 0, 0, time.UTC)

	if t.Weekday() == time.Sunday {
		t = t.AddDate(0, 0, 7)
	} else {
		tmp := 7 - int(t.Weekday()) + 7
		t = t.AddDate(0, 0, tmp)
	}
	return t
}

func Denotcu(year int) time.Time { 
	// Third sunday at the month of June
	t := time.Date(year, time.June, 1, 0, 0, 0, 0, time.UTC)

	if t.Weekday() == time.Sunday {
		t = t.AddDate(0, 0, 7+7)
	} else {
		thirdSunday := 7 - int(t.Weekday()) + 7 + 7
		t = t.AddDate(0, 0, thirdSunday)
	}
	return t
}

