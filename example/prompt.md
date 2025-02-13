You are developing or refactoring a Go project that uses the `pin` library for displaying interactive CLI spinners. The following is a comprehensive description of the library, its public API, and usage examples designed for new projects or refactored code.

---

## Library Overview:

The `pin` library is a lightweight and customizable terminal spinner for Go applications. It provides an elegant progress indicator with support for:

- Custom colors
- Dynamic message updates
- Flexible positioning (spinner before or after the message)
- Custom symbols for success or failure states
- Automatic adjustment in non-interactive environments (animations are disabled when output is piped)

## Installation:

To install the library, run:

```bash
go get github.com/yarlson/pin
```

## Public API:

1. **Creating a New Spinner:**

   - **Constructor:**
     ```go
     func New(message string, opts ...Option) *Pin
     ```
     _Description:_ Initializes a new spinner with a base message and an optional list of functional options for customization.

2. **Controlling the Spinner:**

   - **Start:**

     ```go
     func (p *Pin) Start(ctx context.Context) context.CancelFunc
     ```

     _Description:_ Begins the spinner animation using the provided context. Returns a cancellation function that can be called to stop the spinner.

   - **Stop:**

     ```go
     func (p *Pin) Stop(message ...string)
     ```

     _Description:_ Stops the spinner and outputs an optional final message indicating success or normal termination.

   - **Fail:**

     ```go
     func (p *Pin) Fail(message ...string)
     ```

     _Description:_ Stops the spinner and displays a failure message with a failure-specific symbol and color.

   - **UpdateMessage:**
     ```go
     func (p *Pin) UpdateMessage(message string)
     ```
     _Description:_ Dynamically updates the spinner's displayed message while it is still active.

3. **Functional Options for Customization:**
   These functions return an `Option` that customizes various aspects of the spinner.

   - ```go
     func WithSpinnerColor(color Color) Option
     ```
     _Description:_ Sets the color of the spinner's animation.
   - ```go
     func WithTextColor(color Color) Option
     ```
     _Description:_ Sets the color of the message text.
   - ```go
     func WithDoneSymbol(symbol rune) Option
     ```
     _Description:_ Sets the symbol displayed when the spinner stops successfully.
   - ```go
     func WithDoneSymbolColor(color Color) Option
     ```
     _Description:_ Sets the color of the done symbol.
   - ```go
     func WithPrefix(prefix string) Option
     ```
     _Description:_ Adds a prefix before the spinner and message.
   - ```go
     func WithPrefixColor(color Color) Option
     ```
     _Description:_ Sets the color of the prefix text.
   - ```go
     func WithSeparator(separator string) Option
     ```
     _Description:_ Defines the separator between the prefix and the main message text.
   - ```go
     func WithSeparatorColor(color Color) Option
     ```
     _Description:_ Sets the color for the separator.
   - ```go
     func WithPosition(pos Position) Option
     ```
     _Description:_ Determines the spinner's placement relative to the message text. Use `PositionLeft` (default) or `PositionRight`.
   - ```go
     func WithSpinnerFrames(frames []rune) Option
     ```
     _Description:_ Sets the spinner's frames for custom animations.
   - ```go
     func WithFailSymbol(symbol rune) Option
     ```
     _Description:_ Sets the symbol shown when the spinner indicates a failure.
   - ```go
     func WithFailSymbolColor(color Color) Option
     ```
     _Description:_ Sets the color for the failure symbol.
   - ```go
     func WithFailColor(color Color) Option
     ```
     _Description:_ Sets the color of the failure message text.
   - ```go
     func WithWriter(w io.Writer) Option
     ```
     _Description:_ Redirects spinner output to a custom writer such as `os.Stderr`.

4. **Public Constants:**

   **Colors:**

   ```go
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
   ```

   _Description:_ These constants represent ANSI colors for styling elements of the spinner (text, symbols, animation).

   **Positions:**

   ```go
   const (
       PositionLeft  Position = iota // Spinner appears before the message (default)
       PositionRight                 // Spinner appears after the message
   )
   ```

   _Description:_ These constants allow you to specify the spinner's placement relative to the message text.

## Usage Examples:

1. **Basic Spinner Usage:**

   ```go:example/basic.go
   package main

   import (
       "context"
       "time"
       "github.com/yarlson/pin"
   )

   func main() {
       // Create a spinner with default settings.
       p := pin.New("Loading...",
           pin.WithSpinnerColor(pin.ColorCyan),
           pin.WithTextColor(pin.ColorYellow),
       )

       // Start the spinner.
       cancel := p.Start(context.Background())
       defer cancel()

       // Simulate work.
       time.Sleep(3 * time.Second)

       // Stop the spinner with a final message.
       p.Stop("Done!")
   }
   ```

2. **Advanced Customization with Dynamic Updates:**

   ```go:example/advanced.go
   package main

   import (
       "context"
       "time"
       "github.com/yarlson/pin"
   )

   func main() {
       // Initialize spinner with custom options.
       p := pin.New("Processing",
           pin.WithSpinnerColor(pin.ColorBlue),
           pin.WithTextColor(pin.ColorCyan),
           pin.WithPrefix("Task"),
           pin.WithPrefixColor(pin.ColorYellow),
           pin.WithSeparator("->"),
           pin.WithPosition(pin.PositionRight),
           pin.WithDoneSymbol('✔'),
           pin.WithDoneSymbolColor(pin.ColorGreen),
       )

       // Start the spinner.
       ctx, cancel := context.WithCancel(context.Background())
       defer cancel()
       p.Start(ctx)

       // Update the spinner message while processing.
       time.Sleep(2 * time.Second)
       p.UpdateMessage("Still processing...")
       time.Sleep(2 * time.Second)

       // Stop the spinner with a success message.
       p.Stop("Success!")
   }
   ```

3. **Handling Failure States:**

   ```go:example/fail.go
   package main

   import (
       "context"
       "time"
       "github.com/yarlson/pin"
   )

   func main() {
       // Configure spinner with failure indicators.
       p := pin.New("Deploying",
           pin.WithFailSymbol('✖'),
           pin.WithFailSymbolColor(pin.ColorRed),
           pin.WithFailColor(pin.ColorYellow),
       )

       // Start the spinner.
       ctx, cancel := context.WithCancel(context.Background())
       defer cancel()
       p.Start(ctx)

       // Simulate a failure scenario.
       time.Sleep(2 * time.Second)
       p.Fail("Deployment failed")
   }
   ```

4. **Specifying a Custom Output Destination:**

   ```go:example/custom_writer.go
   package main

   import (
       "context"
       "os"
       "time"
       "github.com/yarlson/pin"
   )

   func main() {
       // Direct spinner output to os.Stderr.
       p := pin.New("Saving Data",
           pin.WithSpinnerColor(pin.ColorMagenta),
           pin.WithWriter(os.Stderr),
       )

       // Start the spinner.
       ctx, cancel := context.WithCancel(context.Background())
       defer cancel()
       p.Start(ctx)

       // Simulate work.
       time.Sleep(3 * time.Second)
       p.Stop("Saved!")
   }
   ```

Usage Context:

This description is intended for integrating or refactoring the pin spinner in your Go project. It details every aspect of the API (public functions and constants) and provides real-world usage examples to simplify the implementation. Use this comprehensive guide to ensure a consistent and interactive CLI experience in your application.
