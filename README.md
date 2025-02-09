# pin

[![Go Reference](https://pkg.go.dev/badge/github.com/yarlson/pin.svg)](https://pkg.go.dev/github.com/yarlson/pin)

`pin` is a lightweight, customizable terminal spinner library for Go applications. It provides an elegant way to show progress and status in CLI applications with support for colors, custom symbols, and flexible positioning.

![Demo](demo.gif)

## Features

- 🎨 Customizable colors for all elements
- 🔄 Smooth braille-pattern animation
- 🎯 Flexible positioning (left/right of message)
- 💫 Configurable prefix and separator
- 🌓 Separator transparency support
- 🔤 UTF-8 symbol support

## Installation

```bash
go get github.com/yarlson/pin
```

## Quick Start

```go
p := pin.New("Loading...")
p.Start()
// do some work
p.Stop("Done!")
```

## Examples

### Basic Progress Indicator

```go
p := pin.New("Processing data")
p.Start()
// ... do work ...
p.UpdateMessage("Almost there...")
// ... finish work ...
p.Stop("Completed!")
```

### Styled Output

```go
p := pin.New("Uploading")
p.SetPrefix("Transfer")
p.SetSeparator("→")
p.SetSpinnerColor(pin.ColorBlue)
p.SetTextColor(pin.ColorCyan)
p.SetPrefixColor(pin.ColorYellow)
p.Start()
// ... do work ...
p.Stop("Upload complete")
```

### Right-side Spinner

```go
p := pin.New("Downloading")
p.SetPosition(pin.PositionRight)
p.Start()
// ... do work ...
p.Stop("Downloaded")
```

### Custom Styling with Alpha

```go
p := pin.New("Processing")
p.SetPrefix("Task")
p.SetSeparator(":")
p.SetSeparatorColor(pin.ColorWhite)
p.SetSeparatorAlpha(0.7)
p.SetDoneSymbol('✔')
p.SetDoneSymbolColor(pin.ColorGreen)
p.Start()
// ... do work ...
p.Stop("Success")
```

## API Reference

### Creating a New Spinner

```go
p := pin.New("message")
```

### Methods

- `Start()` - Starts the spinner animation
- `Stop(message ...string)` - Stops the spinner with optional final message
- `UpdateMessage(message string)` - Updates the current message

### Customization

- `SetSpinnerColor(color Color)`
- `SetTextColor(color Color)`
- `SetPrefix(prefix string)`
- `SetPrefixColor(color Color)`
- `SetSeparator(separator string)`
- `SetSeparatorColor(color Color)`
- `SetSeparatorAlpha(alpha float32)`
- `SetDoneSymbol(symbol rune)`
- `SetDoneSymbolColor(color Color)`
- `SetPosition(pos Position)`

### Available Colors

- `ColorDefault`
- `ColorBlack`
- `ColorRed`
- `ColorGreen`
- `ColorYellow`
- `ColorBlue`
- `ColorMagenta`
- `ColorCyan`
- `ColorWhite`

## Development

### Running Tests

```bash
go test -v ./...
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) for details

## Acknowledgments

Inspired by elegant CLI spinners and the need for a simple, flexible progress indicator in Go applications.
