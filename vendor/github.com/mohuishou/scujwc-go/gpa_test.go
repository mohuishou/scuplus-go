package scujwc

import (
	"testing"
)

func TestGPA(t *testing.T) {
	g, err := j.GPA()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(g)
}

func TestGPAAll(t *testing.T) {
	g, err := j.GPAAll()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(g)
}

func TestGPANotPass(t *testing.T) {
	g, err := j.GPANotPass()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(g)
}
