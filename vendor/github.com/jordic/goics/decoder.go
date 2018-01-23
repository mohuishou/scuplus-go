package goics

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

const (
	keySep    = ":"
	vBegin    = "BEGIN"
	vCalendar = "VCALENDAR"
	vEnd      = "END"
	vEvent    = "VEVENT"

	maxLineRead = 65
)

// Errors
var (
	ErrCalendarNotFound = errors.New("vCalendar not found")
	ErrParseEndCalendar = errors.New("wrong format END:VCALENDAR not Found")
)

type decoder struct {
	scanner      *bufio.Scanner
	err          error
	Calendar     *Calendar
	currentEvent *Event
	nextFn       stateFn
	prevFn       stateFn
	current      string
	buffered     string
	line         int
}

type stateFn func(*decoder)
// NewDecoder creates an instance of de decoder
func NewDecoder(r io.Reader) *decoder {
	d := &decoder{
		scanner:  bufio.NewScanner(r),
		nextFn:   decodeInit,
		line:     0,
		buffered: "",
	}
	return d
}

func (d *decoder) Decode(c ICalConsumer) error {
	d.next()
	if d.Calendar == nil {
		d.err = ErrCalendarNotFound
		d.Calendar = &Calendar{}
	}
	// If theres no error but, nextFn is not reset
	// last element not closed
	if d.nextFn != nil && d.err == nil {
		d.err = ErrParseEndCalendar
	}
	if d.err != nil {
		return d.err
	}

	d.err = c.ConsumeICal(d.Calendar, d.err)
	return d.err
}

// Lines processed. If Decoder reports an error.
// Error
func (d *decoder) Lines() int {
	return d.line
}

// Current Line content
func (d *decoder) CurrentLine() string {
	return d.current
}

// Advances a new line in the decoder
// And calls the next stateFunc
// checks if next line is continuation line
func (d *decoder) next() {
	// If there's not buffered line
	if d.buffered == "" {
		res := d.scanner.Scan()
		if true != res {
			d.err = d.scanner.Err()
			return
		}
		d.line++
		d.current = d.scanner.Text()
	} else {
		d.current = d.buffered
		d.buffered = ""
	}

	if len(d.current) > 65 {
		isContinuation := true
		for isContinuation == true {
			res := d.scanner.Scan()
			d.line++
			if true != res {
				d.err = d.scanner.Err()
				return
			}
			line := d.scanner.Text()
			if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
				d.current = d.current + line[1:]
			} else {
				// If is not a continuation line, buffer it, for the
				// next call.
				d.buffered = line
				isContinuation = false
			}
		}
	}

	if d.nextFn != nil {
		d.nextFn(d)
	}
}

func decodeInit(d *decoder) {
	node := DecodeLine(d.current)
	if node.Key == vBegin && node.Val == vCalendar {
		d.Calendar = &Calendar{
			Data: make(map[string]*IcsNode),
		}
		d.prevFn = decodeInit
		d.nextFn = decodeInsideCal
		d.next()
		return
	}
	d.next()
}

func decodeInsideCal(d *decoder) {
	node := DecodeLine(d.current)
	switch {
	case node.Key == vBegin && node.Val == vEvent:
		d.currentEvent = &Event{
			Data: make(map[string]*IcsNode),
			List: make(map[string][]*IcsNode),
		}
		d.nextFn = decodeInsideEvent
		d.prevFn = decodeInsideCal
	case node.Key == vEnd && node.Val == vCalendar:
		d.nextFn = nil
	default:
		d.Calendar.Data[node.Key] = node
	}
	d.next()
}

func decodeInsideEvent(d *decoder) {

	node := DecodeLine(d.current)
	if node.Key == vEnd && node.Val == vEvent {
		// Come back to parent node
		d.nextFn = d.prevFn
		d.Calendar.Events = append(d.Calendar.Events, d.currentEvent)
		d.next()
		return
	}
	//@todo handle Valarm
	//@todo handle error if we found a startevent without closing pass one
	// #2 handle multiple equal keys. ej. Attendee
	// List keys already set
	if _, ok := d.currentEvent.List[node.Key]; ok {
		d.currentEvent.List[node.Key] = append(d.currentEvent.List[node.Key], node)
	} else {
		// Multiple key detected...
		if val, ok := d.currentEvent.Data[node.Key]; ok {
			d.currentEvent.List[node.Key] = []*IcsNode{val, node}
			delete(d.currentEvent.Data, node.Key)
		} else {
			d.currentEvent.Data[node.Key] = node
		}
	}
	d.next()

}
