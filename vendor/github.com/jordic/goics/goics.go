// Copyright 2015 jordi collell j@tmpo.io. All rights reserved
// Package goics implements an ical encoder and decoder.
// First release will include decode and encoder of Event types

package goics

// ICalConsumer is the realy important part of the decoder lib
// The decoder is organized around the Provider/Consumer pattern.
// the decoder acts as a consummer producing IcsNode's and
// Every data type that wants to receive data, must implement
// the consumer pattern.
type ICalConsumer interface {
	ConsumeICal(d *Calendar, err error) error
}

// ICalEmiter must be implemented in order to allow objects to be serialized
// It should return a *goics.Calendar and optional a map of fields and
// their serializers, if no serializer is defined, it will serialize as
// string..
type ICalEmiter interface {
	EmitICal() Componenter
}

// Componenter defines what should be a component that can be rendered with
// others components inside and some properties
// CALENDAR >> VEVENT ALARM VTODO
type Componenter interface {
	Write(w *ICalEncode)
	AddComponent(c Componenter)
	SetType(t string)
	AddProperty(string, string)
}

// Calendar holds the base struct for a Component VCALENDAR
type Calendar struct {
	Data   map[string]*IcsNode // map of every property found on ics file
	Events []*Event            // slice of events founds in file
}

// Event holds the base struct for a Event Component during decoding
type Event struct {
	Data   map[string]*IcsNode
	Alarms []*map[string]*IcsNode
	// Stores multiple keys for the same property... ( a list )
	List map[string][]*IcsNode
}
