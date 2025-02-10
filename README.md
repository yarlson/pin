# pin

[![Go Reference](https://pkg.go.dev/badge/github.com/yarlson/pin.svg)](https://pkg.go.dev/github.com/yarlson/pin)
[![codecov](https://codecov.io/gh/yarlson/pin/branch/main/graph/badge.svg)](https://codecov.io/gh/yarlson/pin)

`pin` is a lightweight, customizable terminal spinner library for Go applications. It provides an elegant way to show progress and status in CLI applications with support for colors, custom symbols, and flexible positioning.

## Features

- 🎨 Customizable colors for all spinner elements via functional options
- 🔄 Smooth braille-pattern animation
- 🎯 Flexible positioning (spinner before/after message)
- 💫 Configurable prefix and separator
- 🔤 UTF-8 symbol support
- ✨ Ability to update the spinner message dynamically

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
p.Start()
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
- `WithPosition(pos Position)` – sets the spinner's position relative to the message.
  
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

## Development

### Running Tests

```bash
go test -v ./...
```

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
