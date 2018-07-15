package api

import "testing"

func TestGet(t *testing.T) {
	t.Log(Get(1))
	t.Log(Add(1))
	t.Log(Get(1))
}
