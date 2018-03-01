package code

import "testing"

func TestSet(t *testing.T) {
	err := Set("ex@example.com", "123456")
	t.Log(err)
	t.Log(Get("ex@example.com"))
}
