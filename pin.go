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
//
// Example with failure:
//
//	p := pin.New("Processing",
//	    WithFailSymbol('✖'),
//	    WithFailSymbolColor(ColorRed),
//	)
//	cancel := p.Start(context.Background())
//	defer cancel()
//	// ... do some work ...
//	p.Fail("Error occurred")
//
// Example with custom output writer:
//
//	p := pin.New("Saving Data",
//	    WithSpinnerColor(ColorMagenta),
//	    WithWriter(os.Stderr), // send output to stderr
//	)
//	cancel := p.Start(context.Background())
//	defer cancel()
//	// ... do some work ...
//	p.Stop("Saved!")
package pin

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
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
	ColorReset
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

// WithFailSymbol sets the symbol displayed when the spinner fails.
func WithFailSymbol(symbol rune) Option {
	return func(p *Pin) {
		p.failSymbol = symbol
	}
}

// WithFailSymbolColor sets the color of the failure symbol.
func WithFailSymbolColor(color Color) Option {
	return func(p *Pin) {
		p.failSymbolColor = color
	}
}

// WithFailColor sets the color of the failure message text.
// If not set, the failure message is printed using the spinner's text color.
func WithFailColor(color Color) Option {
	return func(p *Pin) {
		p.failColor = color
	}
}

// WithSpinnerFrames sets the frames for the spinner.
// If not set, defaults to the braille symbols. The frames are used from
// beginning to end and then start at the beginning (frames[0]) again
func WithSpinnerFrames(frames []rune) Option {
	return func(p *Pin) {
		p.frames = frames
	}
}

// WithWriter sets a custom io.Writer for spinner output.
func WithWriter(w io.Writer) Option {
	return func(p *Pin) {
		p.out = w
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
//
// You can also indicate failure using the Fail method:
//
//	p := pin.New("Deploying",
//	    WithFailSymbol('✖'),
//	    WithFailSymbolColor(ColorRed),
//	)
//	p.Start()
//	// ... error occurred ...
//	p.Fail("Deployment failed")
type Pin struct {
	frames          []rune
	current         int
	message         string
	messageMu       sync.RWMutex
	stopChan        chan struct{}
	isRunning       int32
	spinnerColor    Color
	textColor       Color
	doneSymbol      rune
	doneSymbolColor Color
	failSymbol      rune
	failSymbolColor Color
	failColor       Color
	prefix          string
	prefixColor     Color
	separator       string
	separatorColor  Color
	position        Position
	out             io.Writer
	wg              sync.WaitGroup
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
		failSymbol:      '✖',
		failSymbolColor: ColorRed,
		failColor:       ColorDefault,
		prefix:          "",
		prefixColor:     ColorDefault,
		separator:       "›",
		separatorColor:  ColorWhite,
		position:        PositionLeft,
		out:             os.Stdout,
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
	if p.IsRunning() {
		return func() {}
	}

	if !isTerminal(p.out) {
		ctx, cancel := context.WithCancel(ctx)
		p.setRunning(true)
		p.messageMu.RLock()
		msg := p.message
		p.messageMu.RUnlock()
		_, _ = fmt.Fprintln(p.out, msg)
		go func() {
			<-ctx.Done()
			p.setRunning(false)
		}()
		return cancel
	}

	p.setRunning(true)

	ctx, cancel := context.WithCancel(ctx)
	ticker := time.NewTicker(100 * time.Millisecond)
	p.wg.Add(1)
	go func() {
		defer ticker.Stop()
		defer p.wg.Done()
		for {
			select {
			case <-p.stopChan:
				return
			case <-ctx.Done():
				p.setRunning(false)
				_, _ = fmt.Fprint(p.out, "\r\033[K")
				return
			case <-ticker.C:
				prefixPart := p.buildPrefixPart()

				p.messageMu.RLock()
				message := p.message
				p.messageMu.RUnlock()

				var format string
				var args []interface{}

				if p.position == PositionLeft {
					format = "\r\033[K%s%s%c%s %s%s%s"
					args = []interface{}{
						prefixPart,
						p.spinnerColor, p.frames[p.current], ColorReset,
						p.textColor, message, ColorReset,
					}
				} else {
					format = "\r\033[K%s%s%s%s %s%c%s "
					args = []interface{}{
						prefixPart,
						p.textColor, message, ColorReset,
						p.textColor, p.frames[p.current], ColorReset,
					}
				}

				_, _ = fmt.Fprintf(p.out, format, args...)
				p.current = (p.current + 1) % len(p.frames)
			}
		}
	}()

	return cancel
}

// Stop halts the spinner animation and optionally displays a final message.
func (p *Pin) Stop(message ...string) {
	if !p.IsRunning() {
		return
	}

	if p.handleNonTerminal(message...) {
		return
	}

	p.setRunning(false)
	p.stopChan <- struct{}{}
	p.wg.Wait()

	_, _ = fmt.Fprint(p.out, "\r\033[K")

	if len(message) > 0 {
		p.printResult(message[0], p.doneSymbol, p.doneSymbolColor)
	}
}

// Fail halts the spinner animation and displays a failure message.
// This method is similar to Stop but uses a distinct symbol and color scheme to indicate an error state.
func (p *Pin) Fail(message ...string) {
	if !p.IsRunning() {
		return
	}

	if p.handleNonTerminal(message...) {
		return
	}

	p.setRunning(false)
	p.stopChan <- struct{}{}
	p.wg.Wait()

	fmt.Print("\r\033[K")

	if len(message) > 0 {
		p.printResult(message[0], p.failSymbol, p.failSymbolColor)
	}
}

// UpdateMessage changes the message shown next to the spinner.
func (p *Pin) UpdateMessage(message string) {
	if !p.IsRunning() {
		return
	}

	p.messageMu.Lock()
	p.message = message
	p.messageMu.Unlock()
	if !isTerminal(p.out) {
		_, _ = fmt.Fprintln(p.out, message)
	}
}

// String returns the ANSI color code for the given color
func (c Color) String() string {
	switch c {
	case ColorReset:
		return "\033[0m"
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

// isTerminal checks if the provided writer is a terminal.
func isTerminal(w io.Writer) bool {
	if ForceInteractive {
		return true
	}

	// Ensure the writer is an *os.File
	f, ok := w.(*os.File)
	if !ok {
		return false
	}

	fi, err := f.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != 0
}

var ForceInteractive bool

// buildPrefixPart constructs the prefix string (including colors) if a prefix is set.
func (p *Pin) buildPrefixPart() string {
	if p.prefix == "" {
		return ""
	}
	return fmt.Sprintf("%s%s%s %s%s%s ", p.prefixColor, p.prefix, ColorReset, p.separatorColor, p.separator, ColorReset)
}

// printResult prints the final message along with a symbol using the appropriate formatting.
func (p *Pin) printResult(msg string, symbol rune, symbolColor Color) {
	var msgColorCode Color
	if symbol == p.failSymbol && p.failColor != ColorDefault {
		msgColorCode = p.failColor
	} else {
		msgColorCode = p.textColor
	}
	prefixPart := p.buildPrefixPart()

	if p.position == PositionLeft {
		format := "%s%s%c%s %s%s%s\n"
		_, _ = fmt.Fprintf(p.out, format, prefixPart, symbolColor, symbol, ColorReset, msgColorCode, msg, ColorReset)
	} else {
		format := "%s%s%s%s %s%c%s\n"
		_, _ = fmt.Fprintf(p.out, format, prefixPart, msgColorCode, msg, ColorReset, symbolColor, symbol, ColorReset)
	}
}

// handleNonTerminal checks if stdout is non-terminal.
// If yes, it prints a plain message (if provided) and returns true.
func (p *Pin) handleNonTerminal(message ...string) bool {
	if !isTerminal(p.out) {
		if len(message) > 0 {
			_, _ = fmt.Fprintln(p.out, message[0])
		}
		p.setRunning(false)
		return true
	}
	return false
}

// Message returns the current spinner message.
func (p *Pin) Message() string {
	return p.message
}

// IsRunning returns whether the spinner is active.
func (p *Pin) IsRunning() bool {
	return atomic.LoadInt32(&p.isRunning) == 1
}

// setRunning sets the running state of the spinner.
func (p *Pin) setRunning(running bool) {
	var val int32
	if running {
		val = 1
	}
	atomic.StoreInt32(&p.isRunning, val)
}
