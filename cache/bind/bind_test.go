package bind

import "testing"

func Test_getExpTime(t *testing.T) {
	t.Log(getExpTime())
}

func TestGet(t *testing.T) {
	t.Log(Get(1, "jwc"))
	t.Log(Add(1, "jwc"))
	t.Log(Get(1, "jwc"))
}
