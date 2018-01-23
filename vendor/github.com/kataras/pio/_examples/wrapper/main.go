package main

import (
	"fmt"
	"time"

	"github.com/kataras/pio"
	"github.com/sirupsen/logrus"
)

/*
	go get -u github.com/sirupsen/logrus
*/

func init() {
	// take an output from a print function
	output := pio.Wrap(logrus.Errorf)
	// register a new printer with name "logrus"
	// which will be able to read text and print as string.
	pio.Register("logrus", output).Marshal(pio.Text)

	// p := pio.Register("logrus", output).Marshal(pio.Text)
	// p.Print("using the logrus printer only")
}

func main() {
	for i := 1; i <= 5; i++ {
		<-time.After(time.Second)
		pio.Print(fmt.Sprintf("[%d] This is an error message that will be printed to the logrus' printer", i))
	}
}
