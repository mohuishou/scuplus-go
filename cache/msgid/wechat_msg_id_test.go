package msgid

import "testing"

func TestSet(t *testing.T) {
	err := Set(1, []string{"1234"})
	t.Log(err)
	err = Set(1, []string{"12345", "123456"})
	t.Log(err)
	t.Log(Get(1))
	t.Log(Get(1))
	t.Log(Get(1))
}
