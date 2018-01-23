package main

import (
	"os"

	"github.com/kataras/pio"
)

// Hijackers can be used to intercept a print operation (per printer),
// they can cancel the print or they can marshal the value
// and make use of the []byte to decide whenever to cancel the print operation,
// if a marshaler is already callled then the printer will not marshal it again
// so it costs nothing.
//
// Below we'll see a simple example, which skips all integer values.
func main() {
	pio.Register("default", os.Stdout).Marshal(pio.Text).Hijack(func(ctx *pio.Ctx) {
		// check if the given value is integer
		// if yes, then cancel the print from that "default" printer.
		if _, ok := ctx.Value.(int); ok {
			ctx.Cancel()
			return
		}

		// ctx.Printer -> gives access to the current printer,
		// at this case the, named as, "default" Printer.

		// ctx.Next() -> to continue to the next hijacker, if available.
	})

	// this should not:
	pio.Print(42)

	pio.Print("this should be the only printed")

	// this should not:
	pio.Print(93)
}
