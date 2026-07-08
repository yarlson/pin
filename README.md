# pin

[![Go Reference](https://pkg.go.dev/badge/github.com/yarlson/pin.svg)](https://pkg.go.dev/github.com/yarlson/pin)
[![codecov](https://codecov.io/gh/yarlson/pin/branch/main/graph/badge.svg)](https://codecov.io/gh/yarlson/pin)

`pin` is a lightweight, customizable terminal spinner library for Go applications. It provides an elegant way to show progress and status in CLI applications with support for colors, custom symbols, and flexible positioning.

![Demo](/assets/demo.gif)

## Features

- 🎨 Customizable colors for all spinner elements via functional options
- 🔄 Smooth braille-pattern animation
- 🎯 Flexible positioning (spinner before/after message)
- 💫 Configurable prefix and separator
- 🔤 UTF-8 symbol support
- ✨ Ability to update the spinner message dynamically
- 🖼️ Customizable spinner frames for unique animation effects
- ⚙️ No external dependencies – uses only the Go standard library
- 🚀 Compatible with Go 1.19 and later
- ⏹ Automatically disables animations in non-interactive (piped) environments to prevent output corruption

## Installation

```bash
go get github.com/yarlson/pin
```

## Quick Start

```go
p := pin.New("Loading...",
    pin.WithSpinnerColor(pin.ColorCyan),
    pin.WithTextColor(pin.ColorYellow),
)
cancel := p.Start(context.Background())
defer cancel()
// do some work
p.Stop("Done!")
```

## Custom Output Writer

You can direct spinner output to an alternative writer (for example, `os.Stderr`) using the `WithWriter` option:

```go
p := pin.New("Processing...",
    pin.WithSpinnerColor(pin.ColorCyan),
    pin.WithTextColor(pin.ColorYellow),
    pin.WithWriter(os.Stderr), // output will be written to stderr
)
cancel := p.Start(context.Background())
defer cancel()
// perform your work
p.Stop("Done!")
```

## Non-interactive Behavior

When the spinner detects that the configured writer is not connected to an interactive terminal (for example, when output is piped), it disables animations and outputs messages as plain text. In this mode:

- The **initial message** is printed immediately when the spinner starts.
- Any **updated messages** are printed as soon as you call `UpdateMessage()`.
- The **final done message** is printed when you call `Stop()`.

## Examples

### Basic Progress Indicator

```go
p := pin.New("Processing data")
cancel := p.Start(context.Background())
defer cancel()
// ... do work ...
p.UpdateMessage("Almost there...")
// finish work
p.Stop("Completed!")
```

### Styled Output

```go
p := pin.New("Uploading",
    pin.WithPrefix("Transfer"),
    pin.WithSeparator("→"),
    pin.WithSpinnerColor(pin.ColorBlue),
    pin.WithTextColor(pin.ColorCyan),
    pin.WithPrefixColor(pin.ColorYellow),
)
cancel := p.Start(context.Background())
defer cancel()
// ... do work ...
p.Stop("Upload complete")
```

### Right-side Spinner

```go
p := pin.New("Downloading", pin.WithPosition(pin.PositionRight))
cancel := p.Start(context.Background())
defer cancel()
// ... do work ...
p.Stop("Downloaded")
```

### Custom Styling & Message Updating

```go
p := pin.New("Processing",
    pin.WithPrefix("Task"),
    pin.WithSeparator(":"),
    pin.WithSeparatorColor(pin.ColorWhite),
    pin.WithDoneSymbol('✔'),
    pin.WithDoneSymbolColor(pin.ColorGreen),
)
cancel := p.Start(context.Background())
defer cancel()

// ... do work ...
p.UpdateMessage("Almost done...")
// finish work
p.Stop("Success")
```

### Failure Indicator

You can express a failure state with the spinner using the new `Fail()` method. Customize the failure appearance with `WithFailSymbol`, `WithFailSymbolColor`, and (optionally) `WithFailColor`.

```go
p := pin.New("Deploying",
    pin.WithFailSymbol('✖'),
    pin.WithFailSymbolColor(pin.ColorRed),
)
cancel := p.Start(context.Background())
defer cancel()
// ... perform tasks ...
p.Fail("Deployment failed")
```

## API Reference

### Creating a New Spinner

```go
p := pin.New("message", /* options... */)
```

### Available Options

- `WithSpinnerColor(color Color)` – sets the spinner's animation color.
- `WithTextColor(color Color)` – sets the color of the message text.
- `WithPrefix(prefix string)` – sets text to display before the spinner.
- `WithPrefixColor(color Color)` – sets the color of the prefix text.
- `WithSeparator(separator string)` – sets the separator text between prefix and message.
- `WithSeparatorColor(color Color)` – sets the color of the separator.
- `WithDoneSymbol(symbol rune)` – sets the symbol displayed upon completion.
- `WithDoneSymbolColor(color Color)` – sets the color of the done symbol.
- `WithFailSymbol(symbol rune)` – sets the symbol displayed upon failure.
- `WithFailSymbolColor(color Color)` – sets the color of the failure symbol.
- `WithFailColor(color Color)` – sets the color of the failure message text.
- `WithPosition(pos Position)` – sets the spinner's position relative to the message.
- `WithSpinnerFrames(frames []rune)` – sets the spinner's frames. Nil or empty frame slices keep the default frames.
- `WithWriter(w io.Writer)` – sets a custom writer for spinner output.

### Interactive Mode Override

- `SetForceInteractive(enabled bool)` – forces interactive rendering for non-terminal writers, mainly for tests and demos.
- `ForceInteractive()` – returns whether forced interactive rendering is enabled.

Earlier versions exposed `ForceInteractive` as a mutable package variable. Use these functions instead so concurrent tests and spinner runs stay race-safe.

### Available Colors

- `ColorDefault`
- `ColorBlack`
- `ColorRed`
- `ColorGreen`
- `ColorYellow`
- `ColorBlue`
- `ColorMagenta`
- `ColorCyan`
- `ColorGray`
- `ColorWhite`

## Development & Compatibility

This library is written using only the Go standard library and supports Go version 1.19 and later.

### Running Tests

```bash
go test -v ./...
```

## Prompt

The LLM prompt example in [example/prompt.md](example/prompt.md) shows you how to quickly integrate pin into your codebase.

### Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## License

MIT License – see [LICENSE](LICENSE) for details

## Acknowledgments

Inspired by elegant CLI spinners and the need for a simple, flexible progress indicator in Go applications.
