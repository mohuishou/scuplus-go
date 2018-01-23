package main

import (
	"io"
	"os"
	"time"

	"github.com/kataras/pio"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		delay = 2 * time.Second
		times = 5
	)

	pio.Register("default", os.Stdout)

	reader, writer := io.Pipe()
	logrus.SetOutput(writer)
	cancel := pio.Scan(reader, true)

	i := 1
	for range time.Tick(delay) {
		logrus.Printf("[%d] Printing %d", i, time.Now().Second())
		if i == times {
			break
		}
		i++
	}

	<-time.After(1 * time.Second)
	cancel()
}
