package goics

import (
	"time"
)

// DateDecode Decodes a date in the distincts formats
func (n *IcsNode) DateDecode() (time.Time, error) {

	// DTEND;VALUE=DATE:20140406
	if val, ok := n.Params["VALUE"]; ok {
		switch {
		case val == "DATE":
			t, err := time.Parse("20060102", n.Val)
			if err != nil {
				return time.Time{}, err
			}
			return t, nil
		case val == "DATE-TIME":
			t, err := time.Parse("20060102T150405", n.Val)
			if err != nil {
				return time.Time{}, err
			}
			return t, nil
		}
	}
	// DTSTART;TZID=Europe/Paris:20140116T120000
	if val, ok := n.Params["TZID"]; ok {
		loc, err := time.LoadLocation(val)
		if err != nil {
			return time.Time{}, err
		}
		t, err := time.ParseInLocation("20060102T150405", n.Val, loc)
		if err != nil {
			return time.Time{}, err
		}
		return t, nil
	}
	//DTSTART:19980119T070000Z utf datetime
	if len(n.Val) == 16 {
		t, err := time.Parse("20060102T150405Z", n.Val)
		if err != nil {
			return time.Time{}, err
		}
		return t, nil
	}

	return time.Time{}, nil
}