package main

import (
	"github.com/kataras/pio"
	"os"
)

func main() {
	p := pio.NewTextPrinter("color", os.Stdout)
	p.Println(pio.Blue("this is a blue text"))
	p.Println(pio.Gray("this is a gray text"))
	p.Println(pio.Red("this is a red text"))
	p.Println(pio.Purple("this is a purple text"))
	p.Println(pio.Yellow("this is a yellow text"))
	p.Println(pio.Green("this is a green text"))
}
