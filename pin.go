// Package pin provides a customizable CLI spinner for showing progress and status in terminal applications.
//
// Example usage:
//
//	p := pin.New("Loading...",
//	    pin.WithSpinnerColor(ColorCyan),
//	    pin.WithTextColor(ColorYellow),
//	)
//	cancel := p.Start(context.Background())
//	defer cancel()
//	// ... do some work ...
//	p.Stop("Done!")
//
// Example with custom styling:
//
//	p := pin.New("Processing",
//	    WithPrefix("Task"),
//	    WithSeparator("→"),
//	    WithSpinnerColor(ColorBlue),
//	    WithTextColor(ColorCyan),
//	    WithPrefixColor(ColorYellow),
//	)
//	cancel := p.Start(context.Background())
//	defer cancel()
//	// ... do some work ...
//	p.Stop("Completed successfully")
//
// Example with right-side positioning:
//
//	p := pin.New("Uploading", WithPosition(PositionRight))
//	cancel := p.Start(context.Background())
//	defer cancel()
//	// ... do some work ...
//	p.UpdateMessage("Almost done...")
//	// ... do more work ...
//	p.Stop("Upload complete")
package pin

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

// Color represents ANSI color codes for terminal output styling.
// Example usage:
//
//	p := pin.New("Loading...", WithTextColor(ColorGreen))
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
	ColorGray
	ColorWhite
)

// Position represents the position of the spinner relative to the message text.
//
// Example usage:
//
//	p := pin.New("Loading", WithPosition(PositionRight)) // Spinner after the message
type Position int

const (
	PositionLeft  Position = iota // Before the message (default)
	PositionRight                 // After the message
)

// Option is a functional option for configuring a Pin.
type Option func(*Pin)

// WithSpinnerColor sets the color of the spinning animation.
func WithSpinnerColor(color Color) Option {
	return func(p *Pin) {
		p.spinnerColor = color
	}
}

// WithTextColor sets the color of the message text.
func WithTextColor(color Color) Option {
	return func(p *Pin) {
		p.textColor = color
	}
}

// WithDoneSymbol sets the symbol displayed when the spinner completes.
func WithDoneSymbol(symbol rune) Option {
	return func(p *Pin) {
		p.doneSymbol = symbol
	}
}

// WithDoneSymbolColor sets the color of the completion symbol.
func WithDoneSymbolColor(color Color) Option {
	return func(p *Pin) {
		p.doneSymbolColor = color
	}
}

// WithPrefix sets the text displayed before the spinner and message.
func WithPrefix(prefix string) Option {
	return func(p *Pin) {
		p.prefix = prefix
	}
}

// WithPrefixColor sets the color of the prefix text.
func WithPrefixColor(color Color) Option {
	return func(p *Pin) {
		p.prefixColor = color
	}
}

// WithSeparator sets the separator text between prefix and message.
func WithSeparator(separator string) Option {
	return func(p *Pin) {
		p.separator = separator
	}
}

// WithSeparatorColor sets the color of the separator.
func WithSeparatorColor(color Color) Option {
	return func(p *Pin) {
		p.separatorColor = color
	}
}

// WithPosition sets whether the spinner appears before or after the message.
func WithPosition(pos Position) Option {
	return func(p *Pin) {
		p.position = pos
	}
}

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
	position        Position
}

var defaultFrames = []rune{
	'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏',
}

// New creates a new Pin instance with the given message and optional configuration options.
// It sets default styling and applies any provided options.
func New(message string, opts ...Option) *Pin {
	p := &Pin{
		frames:          defaultFrames,
		message:         message,
		stopChan:        make(chan struct{}, 1),
		spinnerColor:    ColorDefault,
		textColor:       ColorDefault,
		doneSymbol:      '✓',
		doneSymbolColor: ColorGreen,
		prefix:          "",
		prefixColor:     ColorDefault,
		separator:       "›",
		separatorColor:  ColorWhite,
		position:        PositionLeft,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Start begins the spinner animation using the provided context.
// It returns a cancel function which, when called, will stop the spinner.
// Note: Canceling the returned function stops the spinner without printing
// a final message. To print a final message, use the Stop() method.
func (p *Pin) Start(ctx context.Context) context.CancelFunc {
	if p.isRunning {
		return func() {}
	}

	if !isTerminal(os.Stdout) {
		p.isRunning = true
		go func() {
			<-ctx.Done()
			p.isRunning = false
		}()
		return func() {}
	}

	p.isRunning = true

	ctx, cancel := context.WithCancel(ctx)
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-p.stopChan:
				return
			case <-ctx.Done():
				p.isRunning = false
				fmt.Print("\r\033[K")
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

	return cancel
}

// Stop halts the spinner animation and optionally displays a final message.
func (p *Pin) Stop(message ...string) {
	if !p.isRunning {
		return
	}

	if !isTerminal(os.Stdout) {
		if len(message) > 0 {
			fmt.Println(message[0])
		}
		return
	}
	p.isRunning = false
	p.stopChan <- struct{}{}

	fmt.Print("\r\033[K")

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

// getSeparatorColorCode returns the color code for the separator, applying an alpha effect.
func (p *Pin) getSeparatorColorCode() string {
	if p.separatorColor == ColorDefault {
		return ""
	}
	return p.separatorColor.getColorCode()
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
	case ColorGray:
		return "\033[90m"
	case ColorWhite:
		return "\033[37m"
	default:
		return ""
	}
}

// isTerminal checks if the provided file descriptor is a terminal.
func isTerminal(f *os.File) bool {
	if ForceInteractive {
		return true
	}
	fi, _ := f.Stat()

	return (fi.Mode() & os.ModeCharDevice) != 0
}

var ForceInteractive bool
