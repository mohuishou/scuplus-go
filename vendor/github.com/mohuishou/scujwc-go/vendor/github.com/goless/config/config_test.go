package config

import "testing"

func Test_Config(t *testing.T) {
	config := New("config_test.json")
	d := config.Get("domain")
	host := config.Get("database.host")
	if d == nil {
		t.Error("domain is null")
	}

	if host == nil {
		t.Error("host is null")
	}
}

func BenchmarkLoops(b *testing.B) {
	config := New("config_test.json")
	key := "domain"

	for i := 0; i < b.N; i++ {
		d := config.Get(key)
		if d == nil {
			b.Error(key + " not found")
		}
	}
}
func BenchmarkLoopsDept(b *testing.B) {
	config := New("config_test.json")
	key := "database.test.test.test"
	for i := 0; i < b.N; i++ {
		host := config.Get(key)
		if host == nil {
			b.Error(key + " not found")
		}
	}
}
