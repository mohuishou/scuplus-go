package goics_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	goics "github.com/jordic/goics"
)

func TestComponentCreation(t *testing.T) {

	c := goics.NewComponent()
	c.SetType("VCALENDAR")
	c.AddProperty("CALSCAL", "GREGORIAN")
	c.AddProperty("PRODID", "-//tmpo.io/src/goics")

	if c.Properties["CALSCAL"] != "GREGORIAN" {
		t.Error("Error setting property")
	}

	m := goics.NewComponent()
	m.SetType("VEVENT")
	m.AddProperty("UID", "testing")

	c.AddComponent(m)

	if len(c.Elements) != 1 {
		t.Error("Error adding a component")
	}

}

type EventTest struct {
	component goics.Componenter
}

func (evt *EventTest) EmitICal() goics.Componenter {
	return evt.component
}

func TestWritingSimple(t *testing.T) {
	c := goics.NewComponent()
	c.SetType("VCALENDAR")
	c.AddProperty("CALSCAL", "GREGORIAN")

	ins := &EventTest{
		component: c,
	}

	w := &bytes.Buffer{}
	enc := goics.NewICalEncode(w)
	enc.Encode(ins)

	result := &bytes.Buffer{}
	fmt.Fprintf(result, "BEGIN:VCALENDAR"+goics.CRLF)
	fmt.Fprintf(result, "CALSCAL:GREGORIAN"+goics.CRLF)
	fmt.Fprintf(result, "END:VCALENDAR"+goics.CRLF)

	res := bytes.Compare(w.Bytes(), result.Bytes())
	if res != 0 {
		t.Errorf("%s!=%s %d", w, result, res)

	}

}

func TestFormatDateFieldFormat(t *testing.T) {
	rkey := "DTEND;VALUE=DATE"
	rval := "20140406"
	ti := time.Date(2014, time.April, 06, 0, 0, 0, 0, time.UTC)
	key, val := goics.FormatDateField("DTEND", ti)
	if rkey != key {
		t.Error("Expected", rkey, "Result", key)
	}
	if rval != val {
		t.Error("Expected", rval, "Result", val)
	}
}

func TestFormatDateTimeFieldFormat(t *testing.T) {
	rkey := "X-MYDATETIME;VALUE=DATE-TIME"
	rval := "20120901T130000"
	ti := time.Date(2012, time.September, 01, 13, 0, 0, 0, time.UTC)
	key, val := goics.FormatDateTimeField("X-MYDATETIME", ti)
	if rkey != key {
		t.Error("Expected", rkey, "Result", key)
	}
	if rval != val {
		t.Error("Expected", rval, "Result", val)
	}
}

func TestDateTimeFormat(t *testing.T) {
	rkey := "DTSTART"
	rval := "19980119T070000Z"
	ti := time.Date(1998, time.January, 19, 07, 0, 0, 0, time.UTC)
	key, val := goics.FormatDateTime("DTSTART", ti)
	if rkey != key {
		t.Error("Expected", rkey, "Result", key)
	}
	if rval != val {
		t.Error("Expected", rval, "Result", val)
	}
}

var shortLine = `asdf defined is a test\n\r`

func TestLineWriter(t *testing.T) {

	w := &bytes.Buffer{}

	result := &bytes.Buffer{}
	fmt.Fprintf(result, shortLine)

	encoder := goics.NewICalEncode(w)
	encoder.WriteLine(shortLine)

	res := bytes.Compare(w.Bytes(), result.Bytes())

	if res != 0 {
		t.Errorf("%s!=%s", w, result)
	}

}

var longLine = `As returned by NewWriter, a Writer writes records terminated by thisisat test that is expanded in multi lines` + goics.CRLF

func TestLineWriterLongLine(t *testing.T) {

	w := &bytes.Buffer{}

	result := &bytes.Buffer{}
	fmt.Fprintf(result, "As returned by NewWriter, a Writer writes records terminated by thisisat ")
	fmt.Fprintf(result, goics.CRLFSP)
	fmt.Fprintf(result, "test that is expanded in multi lines")
	fmt.Fprintf(result, goics.CRLF)

	encoder := goics.NewICalEncode(w)
	encoder.WriteLine(longLine)

	res := bytes.Compare(w.Bytes(), result.Bytes())

	if res != 0 {
		t.Errorf("%s!=%s %d", w, result, res)
	}
}

func Test2ongLineWriter(t *testing.T) {
	goics.LineSize = 10

	w := &bytes.Buffer{}

	result := &bytes.Buffer{}
	fmt.Fprintf(result, "12345678")
	fmt.Fprintf(result, goics.CRLF)
	fmt.Fprintf(result, " 2345678")
	fmt.Fprintf(result, goics.CRLF)
	fmt.Fprintf(result, " 2345678")

	var str = `1234567823456782345678`
	encoder := goics.NewICalEncode(w)
	encoder.WriteLine(str)

	res := bytes.Compare(w.Bytes(), result.Bytes())

	if res != 0 {
		t.Errorf("%s!=%s %d", w, result, res)
	}

}
