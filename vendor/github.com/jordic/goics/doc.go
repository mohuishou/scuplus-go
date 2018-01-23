/*
Package goics is a toolkit for encoding and decoding ics/Ical/icalendar files.

This is a work in progress project, that will try to incorporate as many exceptions and variants of the format.

This is a toolkit because user has to define the way it needs the data. The idea is builded with something similar to the consumer/provider pattern.

We want to decode a stream of vevents from a .ics file into a custom type Events

	type Event struct {
		Start, End  time.Time
		Id, Summary string
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
				Id:      node["UID"].Val,
				Summary: node["SUMMARY"].Val,
			}
			*e = append(*e, d)
		}
		return nil
	}

Our custom type, will need to implement ICalConsumer interface, where,
the type will pick up data from the format.
The decoding process will be somthing like this:


	d := goics.NewDecoder(strings.NewReader(testConsumer))
	consumer := Events{}
	err := d.Decode(&consumer)


I have choosed this model, because, this format is a pain and also I don't like a lot the reflect package.

For encoding objects to iCal format, something similar has to be done:

The object emitting elements for the encoder, will have to implement the ICalEmiter, returning a Component structure to be encoded.
This also had been done, because every package could require to encode vals and keys their way. Just for encoding time, I found more than
three types of lines.

	type EventTest struct {
		component goics.Componenter
	}

	func (evt *EventTest) EmitICal() goics.Componenter {
		return evt.component
	}


The Componenter, is an interface that every Component that can be encoded to ical implements.

	c := goics.NewComponent()
	c.SetType("VCALENDAR")
	c.AddProperty("CALSCAL", "GREGORIAN")
	c.AddProperty("PRODID", "-//tmpo.io/src/goics")

	m := goics.NewComponent()
	m.SetType("VEVENT")
	m.AddProperty("UID", "testing")
	c.AddComponent(m)

Properties, had to be stored as strings, the conversion from origin type to string format, must be done,
on the emmiter. There are some helpers for date conversion and on the future I will add more, for encoding
params on the string, and also for handling lists and recurrent events.

A simple example not functional used for testing:

	c := goics.NewComponent()
	c.SetType("VCALENDAR")
	c.AddProperty("CALSCAL", "GREGORIAN")

	ins := &EventTest{
		component: c,
	}

	w := &bytes.Buffer{}
	enc := goics.NewICalEncode(w)
	enc.Encode(ins)


*/
package goics
