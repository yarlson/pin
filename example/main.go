package main

import (
	"context"
	"time"

	"github.com/yarlson/pin"
)

func main() {
	s := pin.New("Loading...")

	s.SetSpinnerColor(pin.ColorCyan)
	s.SetTextColor(pin.ColorYellow)

	s.SetDoneSymbol('✔')
	s.SetDoneSymbolColor(pin.ColorGreen)

	s.SetPrefix("ftl")
	s.SetPrefixColor(pin.ColorMagenta)

	cancel := s.Start(context.Background())
	defer cancel()

	time.Sleep(2 * time.Second)

	s.UpdateMessage("Still working...")
	time.Sleep(2 * time.Second)

	s.Stop("Done!")
}
