package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kataras/pio"
)

type message struct {
	Datetime string `xml:"Date"`
	Message  string `xml:"Message"`
}

func main() {
	p := pio.NewPrinter("default", os.Stdout).WithMarshalers(pio.Text, pio.XMLIndent)
	// Handle registers a handler
	// a handler is always runs in sync AFTER a print operation.
	// It accepts the complete result as its argument
	// and it's able to start other jobs based on that,
	// i.e: printing by other tool, send errors to the cloud logger and etc...
	// Although is possible to add one or more output, or more than one printer
	// per message type, but this is the easy way of doing these things:
	p.Handle(func(result pio.PrintResult) {
		if result.IsOK() {
			fmt.Printf("original value was: %#v\n", result.Value)
		}
	})

	pio.RegisterPrinter(p) // or just use the p.Println...

	pio.Println(message{
		Datetime: time.Now().Format("2006/01/02 - 15:04:05"),
		Message:  "this is an xml message",
	})

	pio.Println("this is a normal text")

}
