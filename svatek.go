package svatek 

import (
	"time"
)
// Compute the exact day for each year
// credits to https://kalendar.beda.cz/vypocet-velikonocni-nedele-v-ruznych-programovacich-jazycich
func Velikonoce(rok int) time.Time {
	if rok <= 1583 {
		rok = 1584
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
	var t time.Time
	if velNedele < 32 {
		t = time.Date(rok, time.March, velNedele, 0, 0, 0, 0, time.UTC)
	} else {
		t = time.Date(rok, time.April, velNedele-31, 0, 0, 0, 0, time.UTC)
	}
	return t
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

func Summertime(year int, end bool) time.Time {
	month := time.March
	if end {
		month = time.October
	}
	return last_sundayofmonth(year, month)
}

func last_sundayofmonth(year int, month time.Month) time.Time {
	t := time.Date(year, month, 31, 0, 0, 0, 0, time.UTC)
	if t.Weekday() == time.Sunday {
		return t
	}
	t = t.AddDate(0, 0, -int(t.Weekday()))
	return t

}
