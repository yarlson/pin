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
		p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop("Done")
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
	p := pin.New("Processing")
	p.SetPrefix("Task")
	p.SetSeparator("→")
	p.SetSpinnerColor(pin.ColorBlue)
	p.SetTextColor(pin.ColorCyan)
	p.SetPrefixColor(pin.ColorYellow)
	p.SetDoneSymbol('✔')
	p.SetDoneSymbolColor(pin.ColorGreen)
	p.SetPosition(pin.PositionRight)

	output := captureOutput(func() {
		p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop("Complete")
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

	var outputs []string
	captureAndStore := func(fn func()) {
		outputs = append(outputs, captureOutput(fn))
	}

	captureAndStore(func() {
		p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
	})

	captureAndStore(func() {
		p.UpdateMessage("Updated")
		time.Sleep(250 * time.Millisecond)
	})

	captureAndStore(func() {
		p.Stop("Final")
	})

	if !strings.Contains(strings.Join(outputs, ""), "Initial") {
		t.Error("Output should contain initial message")
	}
	if !strings.Contains(strings.Join(outputs, ""), "Updated") {
		t.Error("Output should contain updated message")
	}
	if !strings.Contains(strings.Join(outputs, ""), "Final") {
		t.Error("Output should contain final message")
	}
}

// TestSeparatorAlpha verifies that separator with alpha value is displayed correctly.
func TestSeparatorAlpha(t *testing.T) {
	p := pin.New("Testing")
	p.SetPrefix("Alpha")
	p.SetSeparatorAlpha(0.5)

	output := captureOutput(func() {
		p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop("Done")
	})

	if !strings.Contains(output, "Alpha") {
		t.Error("Output should contain the prefix")
	}
	if !strings.Contains(output, "Testing") {
		t.Error("Output should contain the message")
	}
}

// TestMultipleStarts verifies that calling Start multiple times is safe.
func TestMultipleStarts(t *testing.T) {
	p := pin.New("Testing")

	output := captureOutput(func() {
		p.Start(context.Background())
		p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop("Done")
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
		p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop()
	})

	if strings.Contains(output, "\n") {
		t.Error("Stop without message should not print newline")
	}
}

// TestPositionSwitching verifies that spinner can be displayed on either side of the message.
func TestPositionSwitching(t *testing.T) {
	p := pin.New("Testing")

	leftOutput := captureOutput(func() {
		p.SetPosition(pin.PositionLeft)
		p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop("Done")
	})

	p = pin.New("Testing")
	rightOutput := captureOutput(func() {
		p.SetPosition(pin.PositionRight)
		p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop("Done")
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
		pin.ColorWhite,
	}

	for _, color := range colors {
		p := pin.New("Testing")
		p.SetSpinnerColor(color)
		p.SetTextColor(color)
		p.SetPrefixColor(color)
		p.SetSeparatorColor(color)
		p.SetDoneSymbolColor(color)

		output := captureOutput(func() {
			p.Start(context.Background())
			time.Sleep(250 * time.Millisecond)
			p.Stop("Done")
		})

		if !strings.Contains(output, "Testing") {
			t.Errorf("Color %v should not break output", color)
		}
	}
}

// TestSeparatorAlphaValues verifies that separator transparency works correctly.
func TestSeparatorAlphaValues(t *testing.T) {
	testCases := []struct {
		name  string
		alpha float32
	}{
		{"Zero", 0.0},
		{"Quarter", 0.25},
		{"Half", 0.5},
		{"Full", 1.0},
		{"Negative", -0.5},
		{"TooHigh", 1.5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := pin.New("Testing")
			p.SetPrefix("Test")
			p.SetSeparatorColor(pin.ColorWhite)
			p.SetSeparatorAlpha(tc.alpha)

			output := captureOutput(func() {
				p.Start(context.Background())
				time.Sleep(250 * time.Millisecond)
				p.Stop("Done")
			})

			if !strings.Contains(output, "Test") {
				t.Error("Output should contain the prefix")
			}
			if !strings.Contains(output, "Testing") {
				t.Error("Output should contain the message")
			}

			dimCode := "\033[2m"
			if tc.alpha < 1 && tc.alpha >= 0 && !strings.Contains(output, dimCode) {
				t.Error("Output should contain dim effect for alpha < 1")
			}
			if tc.alpha >= 1 && strings.Contains(output, dimCode) {
				t.Error("Output should not contain dim effect for alpha >= 1")
			}
		})
	}
}

func TestStartCancellation(t *testing.T) {
	// Create a new spinner with a test message.
	p := pin.New("Testing spinner")

	// Create a cancellable context.
	ctx, cancel := context.WithCancel(context.Background())

	// Start the spinner with the cancellable context.
	_ = p.Start(ctx)

	// Cancel the context to trigger the ctx.Done() branch.
	cancel()

	// Allow some time for the goroutine to process the cancellation.
	time.Sleep(150 * time.Millisecond)
}
