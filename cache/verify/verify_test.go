package verify

import "testing"

func TestGet(t *testing.T) {
	t.Log(Get(2, "my"))
	t.Log(Set(1, "my", true))
	t.Log(Get(1, "my"))
	t.Log(Set(1, "my2", false))
}
