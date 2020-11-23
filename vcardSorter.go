package main

import (
	"errors"
	"strings"
	"time"

	"github.com/mapaiva/vcard-go"
)

// VCards type is used for sorting
type VCards []vcard.VCard

// Len returns array length
func (c VCards) Len() int {
	return len(c)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (c VCards) Less(i, j int) bool {
	var birthDayI, birthDayJ time.Time
	pBD, err := parseBirthDay(c[i])
	if err != nil {
		birthDayI = time.Date(9999, time.January, 0, 0, 0, 0, 0, time.UTC)
	} else {
		birthDayI = time.Date(1900, pBD.Month(), pBD.Day(), 1, 0, 0, 0, time.UTC)
	}

	pBD, err = parseBirthDay(c[j])
	if err != nil {
		birthDayJ = time.Date(9999, time.January, 0, 0, 0, 0, 0, time.UTC)
	} else {
		birthDayJ = time.Date(1900, pBD.Month(), pBD.Day(), 1, 0, 0, 0, time.UTC)
	}

	return birthDayJ.After(birthDayI)
}

// Swap swaps the elements with indexes i and j.
func (c VCards) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// parseBirthDay parses VCard birthday to time.Time
func parseBirthDay(card vcard.VCard) (time.Time, error) {
	var parsedBirthDay time.Time
	bd := card.BirthDay
	if bd != "" {
		// check the different date formats which are used in VCARDs
		bdTime, err := time.Parse("20060102", bd)
		if err != nil {
			bdTime, err = time.Parse("2006-01-02", bd)
			if err != nil {
				if strings.HasPrefix(bd, "--") {
					// year of birth unknown
					bd = strings.TrimPrefix(bd, "--")
					bd = "0001" + bd
					bdTime, err = time.Parse("20060102", bd)
					if err != nil {
						return parsedBirthDay, err
					}

				} else {
					return parsedBirthDay, errors.New("BirthDay has unknown format")
				}
			}
		}
		return bdTime, nil
	}
	return parsedBirthDay, errors.New("No BirthDay found")
}
