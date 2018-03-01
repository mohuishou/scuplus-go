package msgid

import "testing"

func TestSet(t *testing.T) {
	err := Set(1, "1234")
	t.Log(err)
	err = Set(1, "12345")
	t.Log(err)
	t.Log(Get(1))
	t.Log(Get(1))
	t.Log(Get(1))
}
