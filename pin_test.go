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
		pin.ColorReset,
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

// TestNonInteractiveStop ensures that in non-interactive mode calling Stop with a message
// prints the message (using fmt.Println).
func TestNonInteractiveStop(t *testing.T) {
	// Temporarily force non-interactive mode.
	originalForceInteractive := pin.ForceInteractive
	pin.ForceInteractive = false
	defer func() { pin.ForceInteractive = originalForceInteractive }()

	ctx, cancel := context.WithCancel(context.Background())
	p := pin.New("NonInteractiveTest")
	_ = p.Start(ctx)
	// Allow some time for p.Start's goroutine to spin (even though it prints nothing).
	time.Sleep(100 * time.Millisecond)

	output := captureOutput(func() {
		p.Stop("Completed in non-interactive")
	})
	cancel()

	expected := "Completed in non-interactive\n"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}

// TestNonInteractiveStopWithoutMessage verifies that calling Stop without a final message
// does not print any output when the terminal is non-interactive.
func TestNonInteractiveStopWithoutMessage(t *testing.T) {
	originalForceInteractive := pin.ForceInteractive
	pin.ForceInteractive = false
	defer func() { pin.ForceInteractive = originalForceInteractive }()

	ctx, cancel := context.WithCancel(context.Background())
	p := pin.New("NonInteractiveTest")
	_ = p.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	output := captureOutput(func() {
		p.Stop()
	})
	cancel()

	if output != "" {
		t.Errorf("Expected no output when no message provided, got %q", output)
	}
}

// TestNonInteractiveStart verifies that Start in non-interactive mode does not
// print any spinner output.
func TestNonInteractiveStart(t *testing.T) {
	originalForceInteractive := pin.ForceInteractive
	pin.ForceInteractive = false
	defer func() { pin.ForceInteractive = originalForceInteractive }()

	ctx, cancel := context.WithCancel(context.Background())
	p := pin.New("NonInteractiveTest")
	cancelFn := p.Start(ctx)
	// Allow some time for Start's goroutine (which prints nothing in non-interactive mode).
	time.Sleep(100 * time.Millisecond)
	cancel()
	cancelFn()

	output := captureOutput(func() {
		// No additional printing should occur.
	})
	if output != "" {
		t.Errorf("Expected no output from Start in non-interactive mode, got %q", output)
	}
}

// TestNonInteractiveFullMessageLogging verifies that in non-interactive mode, the
// initial message, updated message, and final done message are logged in sequence.
func TestNonInteractiveFullMessageLogging(t *testing.T) {
	originalForceInteractive := pin.ForceInteractive
	pin.ForceInteractive = false
	defer func() { pin.ForceInteractive = originalForceInteractive }()

	output := captureOutput(func() {
		ctx, cancel := context.WithCancel(context.Background())
		p := pin.New("Initial")
		_ = p.Start(ctx)
		time.Sleep(50 * time.Millisecond)
		p.UpdateMessage("Updated")
		time.Sleep(50 * time.Millisecond)
		p.Stop("Done")
		cancel()
	})

	// Split the captured output into non-empty lines.
	var lines []string
	for _, l := range strings.Split(output, "\n") {
		if strings.TrimSpace(l) != "" {
			lines = append(lines, l)
		}
	}

	expected := []string{"Initial", "Updated", "Done"}
	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines of output, got %d: %v", len(expected), len(lines), lines)
		return
	}
	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d mismatch: expected %q, got %q", i+1, expected[i], line)
		}
	}
}

// TestFail verifies that the Fail method properly displays a failure message.
func TestFail(t *testing.T) {
	p := pin.New("Working",
		pin.WithFailSymbol('✖'),
		pin.WithFailSymbolColor(pin.ColorRed),
		pin.WithPosition(pin.PositionLeft),
	)

	output := captureOutput(func() {
		cancel := p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Fail("Failed")
		cancel()
	})

	if !strings.Contains(output, "Failed") {
		t.Error("Output should contain the failure message")
	}
	if !strings.Contains(output, "✖") {
		t.Error("Output should contain the failure symbol")
	}
}

// TestFailRightPosition verifies failure output with spinner positioned to the right of the text.
func TestFailRightPosition(t *testing.T) {
	p := pin.New("Working",
		pin.WithFailSymbol('✖'),
		pin.WithFailSymbolColor(pin.ColorRed),
		pin.WithPosition(pin.PositionRight),
	)

	output := captureOutput(func() {
		cancel := p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Fail("Failed")
		cancel()
	})

	if !strings.Contains(output, "Failed") {
		t.Error("Output should contain the failure message")
	}
}

// TestNonInteractiveFail ensures that in non-interactive mode, calling Fail with a message prints the message.
func TestNonInteractiveFail(t *testing.T) {
	originalForceInteractive := pin.ForceInteractive
	pin.ForceInteractive = false
	defer func() { pin.ForceInteractive = originalForceInteractive }()

	ctx, cancel := context.WithCancel(context.Background())
	p := pin.New("NonInteractiveTest")
	_ = p.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	output := captureOutput(func() {
		p.Fail("Failed non-interactive")
	})
	cancel()

	expected := "Failed non-interactive\n"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}

// TestNonInteractiveFailWithoutMessage verifies that calling Fail without a final message produces no output in non-interactive mode.
func TestNonInteractiveFailWithoutMessage(t *testing.T) {
	originalForceInteractive := pin.ForceInteractive
	pin.ForceInteractive = false
	defer func() { pin.ForceInteractive = originalForceInteractive }()

	ctx, cancel := context.WithCancel(context.Background())
	p := pin.New("NonInteractiveTest")
	_ = p.Start(ctx)
	time.Sleep(100 * time.Millisecond)

	output := captureOutput(func() {
		p.Fail()
	})
	cancel()

	if output != "" {
		t.Errorf("Expected no output when no message provided, got %q", output)
	}
}

// TestFailNotRunning verifies that calling Fail when the spinner is not running produces no output.
func TestFailNotRunning(t *testing.T) {
	p := pin.New("Not Running")

	output := captureOutput(func() {
		p.Fail("Should not output")
	})

	if output != "" {
		t.Errorf("Expected no output when Fail is called on a non-running spinner, got %q", output)
	}
}

// TestFailHasPrefix verifies that when a prefix is provided, the Fail method prints the prefix along with the failure message.
func TestFailHasPrefix(t *testing.T) {
	p := pin.New("Working",
		pin.WithFailSymbol('✖'),
		pin.WithFailSymbolColor(pin.ColorRed),
		pin.WithPrefix("TestPrefix"),
		pin.WithPrefixColor(pin.ColorBlue),
		pin.WithSeparator(":"),
		pin.WithSeparatorColor(pin.ColorWhite),
		pin.WithPosition(pin.PositionLeft),
	)

	output := captureOutput(func() {
		cancel := p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Fail("Error occurred")
		cancel()
	})

	if !strings.Contains(output, "TestPrefix") {
		t.Error("Expected output to contain prefix 'TestPrefix'")
	}

	if !strings.Contains(output, "Error occurred") {
		t.Error("Expected output to contain failure message 'Error occurred'")
	}
}

// TestIsTerminalNonFile tests the branch where the writer is not an *os.File.
func TestIsTerminalNonFile(t *testing.T) {
	originalForceInteractive := pin.ForceInteractive
	pin.ForceInteractive = false
	defer func() { pin.ForceInteractive = originalForceInteractive }()

	dummy := &bytes.Buffer{}
	p := pin.New("Test Message", pin.WithWriter(dummy))

	output := captureOutput(func() {
		p.UpdateMessage("NonFileWriterTest")
	})
	if !strings.Contains(output, "NonFileWriterTest") {
		t.Error("Expected update message to be printed due to non-*os.File writer")
	}
}

// TestIsTerminalStatError tests the branch where writer.Stat() returns an error.
func TestIsTerminalStatError(t *testing.T) {
	originalForceInteractive := pin.ForceInteractive
	pin.ForceInteractive = false
	defer func() { pin.ForceInteractive = originalForceInteractive }()

	r, _, _ := os.Pipe()
	r.Close()

	p := pin.New("Test Message", pin.WithWriter(r))

	output := captureOutput(func() {
		p.UpdateMessage("StatErrorTest")
	})
	if !strings.Contains(output, "StatErrorTest") {
		t.Error("Expected update message to be printed due to Stat error branch")
	}
}

// TestForceInteractiveSuppressUpdate tests that when ForceInteractive is true,
// UpdateMessage does not print anything.
func TestForceInteractiveSuppressUpdate(t *testing.T) {
	originalForceInteractive := pin.ForceInteractive
	pin.ForceInteractive = true
	defer func() { pin.ForceInteractive = originalForceInteractive }()

	dummy := &bytes.Buffer{}
	p := pin.New("Test Message", pin.WithWriter(dummy))

	output := captureOutput(func() {
		p.UpdateMessage("ForceInteractiveTest")
	})
	if output != "" {
		t.Error("Expected no output when ForceInteractive is true")
	}
}

// TestFailWithCustomFailColor verifies that setting a custom fail color overrides the default text color for the failure message.
func TestFailWithCustomFailColor(t *testing.T) {
	originalForceInteractive := pin.ForceInteractive
	pin.ForceInteractive = true
	defer func() { pin.ForceInteractive = originalForceInteractive }()

	p := pin.New("Working",
		pin.WithFailSymbol('✖'),
		pin.WithFailSymbolColor(pin.ColorRed),
		pin.WithFailColor(pin.ColorMagenta),
		pin.WithTextColor(pin.ColorYellow),
		pin.WithPosition(pin.PositionLeft),
	)

	output := captureOutput(func() {
		cancel := p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Fail("Failed with custom color")
		cancel()
	})

	expectedFailMsgColor := "\033[35m"
	if !strings.Contains(output, expectedFailMsgColor) {
		t.Errorf("Expected output to contain the fail color ANSI code %q, got: %q", expectedFailMsgColor, output)
	}
}

// TestFailWithoutCustomFailColorUsesTextColor verifies that when no custom fail color is set,
// the failure message uses the spinner's text color.
func TestFailWithoutCustomFailColorUsesTextColor(t *testing.T) {
	originalForceInteractive := pin.ForceInteractive
	pin.ForceInteractive = true
	defer func() { pin.ForceInteractive = originalForceInteractive }()

	p := pin.New("Working",
		pin.WithFailSymbol('✖'),
		pin.WithFailSymbolColor(pin.ColorRed),
		pin.WithTextColor(pin.ColorBlue),
		pin.WithPosition(pin.PositionLeft),
	)

	output := captureOutput(func() {
		cancel := p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Fail("Failed using text color")
		cancel()
	})

	expectedTextColorCode := "\033[34m"
	if !strings.Contains(output, expectedTextColorCode) {
		t.Errorf("Expected output to contain text color ANSI code %q, got: %q", expectedTextColorCode, output)
	}
}

// TestAllColors verifies that custom spinner configs work correctly
// Note that since spinner frames only show up on proper terminals and thus
// can't be captured, we can't really verify that they were emitted.
func TestSpinnerFrames(t *testing.T) {
	framesets := []string{
		".oO0Oo",
		"|/-\\",
	}

	for _, frames := range framesets {
		p := pin.New("Testing",
			pin.WithSpinnerFrames([]rune(frames)),
		)
		cancel := p.Start(context.Background())
		time.Sleep(250 * time.Millisecond)
		p.Stop("Done")
		cancel()
	}
}
