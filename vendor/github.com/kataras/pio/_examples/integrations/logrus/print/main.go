package main

import (
	"github.com/kataras/pio"
	_ "github.com/kataras/pio/_examples/integrations/logrus"
)

func main() {
	pio.Print("This is an info message that will be printed to the logrus' printer")
}
