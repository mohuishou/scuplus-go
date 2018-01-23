package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kataras/pio"
	_ "github.com/kataras/pio/_examples/integrations/logrus"
)

func main() {
	var (
		delay = 2 * time.Second
		times = 5
	)

	pio.Get("logrus").AddOutput(os.Stdout)

	i := 1
	for range time.Tick(delay) {
		pio.Print(fmt.Sprintf("[%d] Printing %d\n", i, time.Now().Second()))
		if i == times {
			break
		}
		i++

	}
	<-time.After(1 * time.Second)
}
