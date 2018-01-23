/*
Package logrus registers the global logrus logger to the pio ecosystem,
install logrus first:


	$ go get github.com/


Example Code:


	package main

	import (
		"github.com/kataras/pio"
		_ "github.com/kataras/pio/_examples/integrations/logrus"
	)

	func main() {
		pio.Print("This is an info message that will be printed to the logrus' output")
	}

*/
package logrus

import (
	"github.com/kataras/pio"
	"github.com/sirupsen/logrus"
)

// Name of this printer.
const Name = "logrus"

func init() {
	Register(logrus.Infof)
}

// sirupsen/logrus/entry.go#96
// func (h *hook) Fire(e *logrus.Entry)
// no f. it. It doesn't work as expected
// they don't made it to be able to stop a
// specific print operation or change its buffer
// because the entry's buffer is being initialized
// after the hook's Fire, so we can't
// create it with a specific hook, do it
// with the pio's own way:

// Register registers a specific logrus printf-compatible function signature
// to the pio registry.
//
// pio can take only one by-design because it is not based on any log levels
// but, you can extend it by calling its Hijack function
// to determinate what to log.
func Register(printf func(string, ...interface{})) *pio.Printer {
	return pio.Register("logrus", pio.Wrap(printf)).Marshal(pio.Text)
}

// Remove removes the logrus printer from the pio.
func Remove() {
	pio.Remove("logrus")
}
