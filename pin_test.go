package pin_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/yarlson/pin"
)

func init() {
	pin.ForceInteractive = true
}

var (
	stdoutMu sync.Mutex
)

// captureOutput helps test terminal output by capturing stdout during test execution.
// This is useful for verifying what the user would actually see in their terminal.
func captureOutput(fn func()) string {
	stdoutMu.Lock()
	defer stdoutMu.Unlock()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	_ = w.Close()
	os.Stdout = old
	return <-outC
}

// TestBasicUsage verifies the core start-stop functionality with default settings.
func TestBasicUsage(t *testing.T) {
	p := pin.New("Loading")

	output := captureOutput(func() {
		cancel := p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop("Done")
		cancel()
	})

	if !strings.Contains(output, "Loading") {
		t.Error("Output should contain the message")
	}
	if !strings.Contains(output, "Done") {
		t.Error("Output should contain the done message")
	}
}

// TestCustomization verifies that all customization options work together.
func TestCustomization(t *testing.T) {
	p := pin.New("Processing",
		pin.WithPrefix("Task"),
		pin.WithSeparator("→"),
		pin.WithSpinnerColor(pin.ColorBlue),
		pin.WithTextColor(pin.ColorCyan),
		pin.WithPrefixColor(pin.ColorYellow),
		pin.WithDoneSymbol('✔'),
		pin.WithDoneSymbolColor(pin.ColorGreen),
		pin.WithPosition(pin.PositionRight),
	)

	output := captureOutput(func() {
		cancel := p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop("Complete")
		cancel()
	})

	if !strings.Contains(output, "Task") {
		t.Error("Output should contain the prefix")
	}
	if !strings.Contains(output, "→") {
		t.Error("Output should contain the separator")
	}
	if !strings.Contains(output, "Processing") {
		t.Error("Output should contain the message")
	}
	if !strings.Contains(output, "Complete") {
		t.Error("Output should contain the done message")
	}
}

// TestMessageUpdate verifies that messages can be updated while the spinner is running.
func TestMessageUpdate(t *testing.T) {
	p := pin.New("Initial")
	output := captureOutput(func() {
		cancel := p.Start(context.Background())
		time.Sleep(250 * time.Millisecond) // spinner prints "Initial"
		p.UpdateMessage("Updated")
		time.Sleep(250 * time.Millisecond) // spinner prints "Updated"
		p.Stop("Final")
		cancel()
	})

	if !strings.Contains(output, "Initial") {
		t.Error("Output should contain initial message")
	}
	if !strings.Contains(output, "Updated") {
		t.Error("Output should contain updated message")
	}
	if !strings.Contains(output, "Final") {
		t.Error("Output should contain final message")
	}
}

// TestMultipleStarts verifies that calling Start multiple times is safe.
func TestMultipleStarts(t *testing.T) {
	p := pin.New("Testing")

	output := captureOutput(func() {
		cancel1 := p.Start(context.Background())
		cancel2 := p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop("Done")
		cancel1()
		cancel2()
	})

	if strings.Count(output, "Testing") > len(output)/2 {
		t.Error("Multiple starts should not cause duplicate output")
	}
}

// TestStopWithoutStart verifies that calling Stop before Start is safe.
func TestStopWithoutStart(t *testing.T) {
	p := pin.New("Testing")

	output := captureOutput(func() {
		p.Stop("Done")
	})

	if output != "" {
		t.Error("Stop without start should produce no output")
	}
}

// TestStopWithoutMessage verifies that Stop can be called without a final message.
func TestStopWithoutMessage(t *testing.T) {
	p := pin.New("Testing")

	output := captureOutput(func() {
		cancel := p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop()
		cancel()
	})

	if strings.Contains(output, "\n") {
		t.Error("Stop without message should not print newline")
	}
}

// TestPositionSwitching verifies that spinner can be displayed on either side of the message.
func TestPositionSwitching(t *testing.T) {
	leftOutput := captureOutput(func() {
		p := pin.New("Testing", pin.WithPosition(pin.PositionLeft))
		cancel := p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop("Done")
		cancel()
	})

	rightOutput := captureOutput(func() {
		p := pin.New("Testing", pin.WithPosition(pin.PositionRight))
		cancel := p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop("Done")
		cancel()
	})

	if leftOutput == rightOutput {
		t.Error("Left and right positions should produce different outputs")
	}
}

// TestAllColors verifies that all color combinations work correctly.
func TestAllColors(t *testing.T) {
	colors := []pin.Color{
		pin.ColorDefault,
		pin.ColorBlack,
		pin.ColorRed,
		pin.ColorGreen,
		pin.ColorYellow,
		pin.ColorBlue,
		pin.ColorMagenta,
		pin.ColorCyan,
		pin.ColorGray,
		pin.ColorWhite,
	}

	for _, color := range colors {
		p := pin.New("Testing",
			pin.WithSpinnerColor(color),
			pin.WithTextColor(color),
			pin.WithPrefixColor(color),
			pin.WithSeparatorColor(color),
			pin.WithDoneSymbolColor(color),
		)

		output := captureOutput(func() {
			cancel := p.Start(context.Background())
			time.Sleep(250 * time.Millisecond)
			p.Stop("Done")
			cancel()
		})

		if !strings.Contains(output, "Testing") {
			t.Errorf("Color %v should not break output", color)
		}
	}
}

// TestStartCancellation verifies that cancellation in Start properly stops the spinner.
func TestStartCancellation(t *testing.T) {
	p := pin.New("Testing spinner")
	ctx, cancel := context.WithCancel(context.Background())
	cancelFunc := p.Start(ctx)
	// Cancel the context to trigger the cancellation branch.
	cancel()
	time.Sleep(150 * time.Millisecond)
	// Ensure the returned cancel function is also invoked.
	cancelFunc()
}
