package goics_test

import (
	"strings"
	"testing"
	"time"

	goics "github.com/jordic/goics"
)

type Event struct {
	Start, End  time.Time
	ID, Summary string
}

type Events []Event

func (e *Events) ConsumeICal(c *goics.Calendar, err error) error {
	for _, el := range c.Events {
		node := el.Data
		dtstart, err := node["DTSTART"].DateDecode()
		if err != nil {
			return err
		}
		dtend, err := node["DTEND"].DateDecode()
		if err != nil {
			return err
		}
		d := Event{
			Start:   dtstart,
			End:     dtend,
			ID:      node["UID"].Val,
			Summary: node["SUMMARY"].Val,
		}
		*e = append(*e, d)
	}
	return nil
}

func TestConsumer(t *testing.T) {
	d := goics.NewDecoder(strings.NewReader(testConsumer))
	consumer := Events{}
	err := d.Decode(&consumer)
	if err != nil {
		t.Error("Unable to consume ics")
	}
	if len(consumer) != 1 {
		t.Error("Incorrect length decoding container", len(consumer))
	}

	if consumer[0].Start != time.Date(2014, time.April, 06, 0, 0, 0, 0, time.UTC) {
		t.Error("Expected", consumer[0].Start)
	}
	if consumer[0].ID != "-kpd6p8pqal11-n66f1wk1tw76@xxxx.com" {
		t.Errorf("Error decoding text")
	}
}

var testConsumer = `BEGIN:VCALENDAR
PRODID;X-RICAL-TZSOURCE=TZINFO:-//test//EN
CALSCALE:GREGORIAN
VERSION:2.0
BEGIN:VEVENT
DTEND;VALUE=DATE:20140506
DTSTART;VALUE=DATE:20140406
UID:-kpd6p8pqal11-n66f1wk1tw76@xxxx.com
DESCRIPTION:CHECKIN:  01/05/2014\nCHECKOUT: 06/05/2014\nNIGHTS:   5\nPHON
 E:    \nEMAIL:    (no se ha facilitado ningún correo electrónico)\nPRO
 PERTY: Apartamento xxx 6-8 pax en Centro\n
SUMMARY:Luigi Carta (FYSPZN)
LOCATION:Apartamento xxx 6-8 pax en Centro
END:VEVENT
END:VCALENDAR
`