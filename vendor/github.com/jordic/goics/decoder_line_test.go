package goics_test

import (
	//"strings"
	"testing"

	goics "github.com/jordic/goics"
)

var TLines = []string{
	"BEGIN:VEVENT",
	`ATTENDEE;RSVP=TRUE;ROLE=REQ-PARTICIPANT:mailto:jsmith@example.com`,
	`X-FOO;BAR=";hidden=value";FOO=baz:realvalue`,
	`ATTENDEE;ROLE="REQ-PARTICIPANT;foo";PARTSTAT=ACCEPTED;RSVP=TRUE:mailto:foo@bar.com`,
}

func TestDecodeLine(t *testing.T) {
	node := goics.DecodeLine(TLines[0])
	if node.Key != "BEGIN" {
		t.Errorf("Wrong key parsing %s", node.Key)
	}
	if node.Val != "VEVENT" {
		t.Errorf("Wrong key parsing %s", node.Key)
	}
	if node.ParamsLen() != 0 {
		t.Error("No keys")
	}

	node = goics.DecodeLine(TLines[1])
	if node.Key != "ATTENDEE" {
		t.Errorf("Wrong key parsing %s", node.Key)
	}
	if node.Val != "mailto:jsmith@example.com" {
		t.Errorf("Wrong val parsing %s", node.Val)
	}
	if node.ParamsLen() != 2 {
		t.Errorf("Wrong nmber of params %s", node.ParamsLen())
	}
	node = goics.DecodeLine(TLines[2])
	if node.Key != "X-FOO" {
		t.Errorf("Wrong key parsing %s", node.Key)
	}
	if node.ParamsLen() != 2 {
		t.Errorf("Incorrect quoted params count %s, %d", node.ParamsLen(), node.Params)

	}
	node = goics.DecodeLine(TLines[3])
	if node.Key != "ATTENDEE" {
		t.Errorf("Wrong key parsing %s", node.Key)
	}
	if node.ParamsLen() != 3 {
		t.Errorf("Incorrect quoted params count %s, %d", node.ParamsLen(), node.Params)
	}
	if node.Params["ROLE"] != "REQ-PARTICIPANT;foo" {
		t.Errorf("Error extracting quoted params from line %s", node.Params["ROLE"])
	}
}