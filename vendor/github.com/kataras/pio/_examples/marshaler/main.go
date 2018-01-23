package main

import (
	"os"
	"time"

	"github.com/kataras/pio"
)

type message struct {
	From string `json:"printer_name"`
	// fields should be exported, as you already know.
	Order    int    `json:"order"`
	Datetime string `json:"date"`
	Message  string `json:"message"`
}

func main() {
	// showcase json
	println("-----------")
	printWith("json", pio.JSONIndent)

	// showcase xml
	println("-----------")
	printWith("xml", pio.XMLIndent)

	// show case text
	println("-----------")
	pio.Register("text", os.Stderr).Marshal(pio.Text)
	pio.Println("this is a text message, from text printer that has been registered inline")

	print("-----------")
}

func printWith(printerName string, marshaler pio.Marshaler) {
	// this "json" printer is responsible to print
	// text (string) and json structs.
	// use `pio#JSON` or `pio#JSONIdent` for nicer print format
	// and `encoding/json#Marshal` is the same thing,
	// pio is fully compatible with standard marshal functions.
	p := pio.NewPrinter(printerName, os.Stderr).
		Marshal(marshaler)

	p.Println(message{
		From:     printerName,
		Order:    1,
		Datetime: time.Now().Format("2006/01/02 - 15:04:05"),
		Message:  "This is our structed error log message",
	})

	p.Println(message{
		From:     printerName,
		Order:    2,
		Datetime: time.Now().Format("2006/01/02 - 15:04:05"),
		Message:  "This is our second structed error log message",
	})

	// this will print only xml because we use a single printer to print
	//
	// p := pio.Register("xml", os.Stderr).WithMarshalers(pio.XMLIndent)

	// p.Print(message{
	// 	Order:    3,
	// 	Datetime: time.Now().Format("2006/01/02 - 15:04:05"),
	// 	Message:  "This is our second structed error log message",
	// })
}
