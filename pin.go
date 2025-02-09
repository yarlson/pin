// Package pin provides a customizable CLI spinner for showing progress and status in terminal applications.
//
// Example usage:
//
//	p := pin.New("Loading...")
//	p.Start()
//	// ... do some work ...
//	p.Stop("Done!")
//
// Example with custom styling:
//
//	p := pin.New("Processing")
//	p.SetPrefix("Task")
//	p.SetSeparator("→")
//	p.SetSpinnerColor(pin.ColorBlue)
//	p.SetTextColor(pin.ColorCyan)
//	p.SetPrefixColor(pin.ColorYellow)
//	p.Start()
//	// ... do some work ...
//	p.Stop("Completed successfully")
//
// Example with right-side positioning:
//
//	p := pin.New("Uploading")
//	p.SetPosition(pin.PositionRight)
//	p.Start()
//	// ... do some work ...
//	p.UpdateMessage("Almost done...")
//	// ... do more work ...
//	p.Stop("Upload complete")
package pin

import (
	"fmt"
	"sync"
	"time"
)

// Color represents ANSI color codes for terminal output styling.
//
// Example usage:
//
//	p.SetTextColor(pin.ColorGreen)
//	p.SetSpinnerColor(pin.ColorBlue)
type Color int

const (
	ColorDefault Color = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

// Position represents the position of the spinner relative to the message text.
//
// Example usage:
//
//	p.SetPosition(pin.PositionRight) // Places spinner after the message
//	p.SetPosition(pin.PositionLeft)  // Places spinner before the message (default)
type Position int

const (
	PositionLeft  Position = iota // Before the message (default)
	PositionRight                 // After the message
)

// Pin represents an animated terminal spinner with customizable appearance and behavior.
// It supports custom colors, symbols, prefixes, and positioning.
//
// Basic usage:
//
//	p := pin.New("Loading")
//	p.Start()
//	time.Sleep(2 * time.Second)
//	p.Stop("Done")
//
// Advanced usage:
//
//	p := pin.New("Processing")
//	p.SetPrefix("Status")
//	p.SetSeparator(":")
//	p.SetSeparatorColor(pin.ColorWhite)
//	p.SetSeparatorAlpha(0.7)
//	p.SetSpinnerColor(pin.ColorCyan)
//	p.SetTextColor(pin.ColorYellow)
//	p.Start()
//
//	// Update message during operation
//	p.UpdateMessage("Still working...")
//
//	// Complete with success
//	p.SetDoneSymbolColor(pin.ColorGreen)
//	p.Stop("Completed!")
type Pin struct {
	frames          []rune
	current         int
	message         string
	messageMu       sync.RWMutex
	stopChan        chan struct{}
	isRunning       bool
	spinnerColor    Color
	textColor       Color
	doneSymbol      rune
	doneSymbolColor Color
	prefix          string
	prefixColor     Color
	separator       string
	separatorColor  Color
	separatorAlpha  float32
	position        Position
}

// braille patterns for spinning animation
var defaultFrames = []rune{
	'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏',
}

// getColorCode returns the ANSI color code for the given color
func (c Color) getColorCode() string {
	switch c {
	case ColorBlack:
		return "\033[30m"
	case ColorRed:
		return "\033[31m"
	case ColorGreen:
		return "\033[32m"
	case ColorYellow:
		return "\033[33m"
	case ColorBlue:
		return "\033[34m"
	case ColorMagenta:
		return "\033[35m"
	case ColorCyan:
		return "\033[36m"
	case ColorWhite:
		return "\033[37m"
	default:
		return ""
	}
}

// New creates a new Pin instance with the given message.
// The pin starts with default styling and left-side positioning.
func New(message string) *Pin {
	return &Pin{
		frames:          defaultFrames,
		message:         message,
		stopChan:        make(chan struct{}),
		spinnerColor:    ColorDefault,
		textColor:       ColorDefault,
		doneSymbol:      '✓',
		doneSymbolColor: ColorGreen,
		prefix:          "",
		prefixColor:     ColorDefault,
		separator:       "›",
		separatorColor:  ColorWhite,
		separatorAlpha:  0.5,
		position:        PositionLeft,
	}
}

// SetSpinnerColor sets the color of the spinning animation.
func (p *Pin) SetSpinnerColor(color Color) {
	p.spinnerColor = color
}

// SetTextColor sets the color of the message text.
func (p *Pin) SetTextColor(color Color) {
	p.textColor = color
}

// SetDoneSymbol sets the symbol displayed when the spinner completes.
func (p *Pin) SetDoneSymbol(symbol rune) {
	p.doneSymbol = symbol
}

// SetDoneSymbolColor sets the color of the completion symbol.
func (p *Pin) SetDoneSymbolColor(color Color) {
	p.doneSymbolColor = color
}

// SetPrefix sets the text displayed before the spinner and message.
func (p *Pin) SetPrefix(prefix string) {
	p.prefix = prefix
}

// SetPrefixColor sets the color of the prefix text.
func (p *Pin) SetPrefixColor(color Color) {
	p.prefixColor = color
}

// SetSeparator sets the separator text between prefix and message.
func (p *Pin) SetSeparator(separator string) {
	p.separator = separator
}

// SetSeparatorColor sets the color of the separator.
func (p *Pin) SetSeparatorColor(color Color) {
	p.separatorColor = color
}

// SetSeparatorAlpha sets the opacity of the separator between 0.0 and 1.0.
func (p *Pin) SetSeparatorAlpha(alpha float32) {
	if alpha < 0 {
		alpha = 0
	}
	if alpha > 1 {
		alpha = 1
	}
	p.separatorAlpha = alpha
}

// getSeparatorColorCode returns the color code with alpha applied
func (p *Pin) getSeparatorColorCode() string {
	if p.separatorColor == ColorDefault {
		return ""
	}

	// Convert regular color to dim color (alpha effect) if alpha is less than 1
	if p.separatorAlpha < 1 {
		return "\033[2m" + p.separatorColor.getColorCode()
	}
	return p.separatorColor.getColorCode()
}

// SetPosition sets whether the spinner appears before or after the message.
func (p *Pin) SetPosition(pos Position) {
	p.position = pos
}

// Start begins the spinner animation.
func (p *Pin) Start() {
	if p.isRunning {
		return
	}
	p.isRunning = true

	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-p.stopChan:
				return
			case <-ticker.C:
				spinnerColorCode := p.spinnerColor.getColorCode()
				textColorCode := p.textColor.getColorCode()
				prefixColorCode := p.prefixColor.getColorCode()
				separatorColorCode := p.getSeparatorColorCode()
				reset := "\033[0m"

				prefixPart := ""
				if p.prefix != "" {
					prefixPart = fmt.Sprintf("%s%s%s %s%s%s ",
						prefixColorCode, p.prefix, reset,
						separatorColorCode, p.separator, reset)
				}

				p.messageMu.RLock()
				message := p.message
				p.messageMu.RUnlock()

				var format string
				var args []interface{}

				if p.position == PositionLeft {
					format = "\r%s%s%c%s %s%s%s"
					args = []interface{}{
						prefixPart,
						spinnerColorCode, p.frames[p.current], reset,
						textColorCode, message, reset,
					}
				} else {
					format = "\r%s%s%s%s %s%c%s "
					args = []interface{}{
						prefixPart,
						textColorCode, message, reset,
						spinnerColorCode, p.frames[p.current], reset,
					}
				}

				fmt.Printf(format, args...)
				p.current = (p.current + 1) % len(p.frames)
			}
		}
	}()
}

// Stop halts the spinner animation and optionally displays a final message.
func (p *Pin) Stop(message ...string) {
	if !p.isRunning {
		return
	}
	p.isRunning = false
	p.stopChan <- struct{}{}

	// Clear the entire line and return cursor to start
	fmt.Print("\r\033[K")

	// If a final message was provided, display it with the done symbol
	if len(message) > 0 {
		prefixColorCode := p.prefixColor.getColorCode()
		symbolColorCode := p.doneSymbolColor.getColorCode()
		textColorCode := p.textColor.getColorCode()
		separatorColorCode := p.getSeparatorColorCode()
		reset := "\033[0m"

		prefixPart := ""
		if p.prefix != "" {
			prefixPart = fmt.Sprintf("%s%s%s %s%s%s ",
				prefixColorCode, p.prefix, reset,
				separatorColorCode, p.separator, reset)
		}

		var format string
		var args []interface{}

		if p.position == PositionLeft {
			format = "%s%s%c%s %s%s%s\n"
			args = []interface{}{
				prefixPart,
				symbolColorCode, p.doneSymbol, reset,
				textColorCode, message[0], reset,
			}
		} else {
			format = "%s%s%s%s %s%c%s \n"
			args = []interface{}{
				prefixPart,
				textColorCode, message[0], reset,
				symbolColorCode, p.doneSymbol, reset,
			}
		}

		fmt.Printf(format, args...)
	}
}

// UpdateMessage changes the message shown next to the spinner.
func (p *Pin) UpdateMessage(message string) {
	p.messageMu.Lock()
	p.message = message
	p.messageMu.Unlock()
}
