package main

import (
	"context"
	"time"

	"github.com/yarlson/pin"
)

func main() {
	s := pin.New("Loading...",
		pin.WithSpinnerColor(pin.ColorCyan),
		pin.WithTextColor(pin.ColorYellow),
		pin.WithDoneSymbol('✔'),
		pin.WithDoneSymbolColor(pin.ColorGreen),
		pin.WithPrefix("ftl"),
		pin.WithPrefixColor(pin.ColorMagenta),
		pin.WithSeparatorColor(pin.ColorGray),
	)

	cancel := s.Start(context.Background())
	defer cancel()

	time.Sleep(2 * time.Second)

	s.UpdateMessage("Still working...")
	time.Sleep(2 * time.Second)

	s.Stop("Done!")
}
