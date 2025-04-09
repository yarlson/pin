package pin_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/yarlson/pin"
)

func TestNewCreatesSpinner(t *testing.T) {
	message := "Loading..."
	p := pin.New(message)
	if p == nil {
		t.Fatal("Expected a non-nil spinner instance")
	}
	if p.Message() != message {
		t.Fatalf("Expected message %q, got %q", message, p.Message())
	}
}

func TestStartAndCancel(t *testing.T) {
	p := pin.New("Loading...")
	cancel := p.Start(context.Background())
	// Immediately, the spinner should be running.
	if !p.IsRunning() {
		t.Fatal("Expected spinner to be running after Start()")
	}
	// Cancel the spinner.
	cancel()
	// Allow some time for the cancellation to propagate.
	time.Sleep(100 * time.Millisecond)
	if p.IsRunning() {
		t.Fatal("Expected spinner to have stopped after cancellation")
	}
}

func TestStopPrintsMessage(t *testing.T) {
	var buf bytes.Buffer

	// Create a spinner with a custom writer so output can be captured.
	p := pin.New("Processing...", pin.WithWriter(&buf))
	// Start the spinner.
	cancel := p.Start(context.Background())
	// Cancel to simulate spinner stopping (ensuring any background goroutines complete).
	cancel()

	// Now call Stop with a final message.
	p.Stop("Done!")

	output := buf.String()
	if !strings.Contains(output, "Done!") {
		t.Errorf("Expected output to contain final message 'Done!', got %q", output)
	}
	// Also verify spinner is no longer running.
	if p.IsRunning() {
		t.Error("Expected spinner to not be running after Stop()")
	}
}

func TestUpdateMessagePrints(t *testing.T) {
	var buf bytes.Buffer
	// Create a spinner with a custom writer so we can capture output.
	p := pin.New("Initial", pin.WithWriter(&buf))

	// Update the spinner message.
	p.UpdateMessage("Updated")

	output := buf.String()
	if !strings.Contains(output, "Updated") {
		t.Errorf("Expected output to contain 'Updated', got %q", output)
	}
}

func TestSpinnerAnimation(t *testing.T) {
	// Force interactive mode for this test.
	pin.ForceInteractive = true
	defer func() { pin.ForceInteractive = false }()

	var buf bytes.Buffer
	// Create a spinner with a custom writer so we can capture output.
	p := pin.New("Animating", pin.WithWriter(&buf))
	cancel := p.Start(context.Background())
	defer cancel()

	// Let the spinner animate for a short while.
	time.Sleep(150 * time.Millisecond)
	p.UpdateMessage("Updated")
	time.Sleep(150 * time.Millisecond)
	p.Stop("Stopped")

	output := buf.String()
	frames := []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
	found := false
	for _, frame := range frames {
		if strings.Contains(output, string(frame)) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected output to contain one of the spinner frames, got %q", output)
	}
}

func TestFailPrintsFailureMessage(t *testing.T) {
	// Force interactive mode for this test.
	pin.ForceInteractive = true
	defer func() { pin.ForceInteractive = false }()

	var buf bytes.Buffer
	// Create a spinner with a custom writer.
	p := pin.New("Working", pin.WithWriter(&buf))
	cancel := p.Start(context.Background())
	defer cancel()
	// Allow some time for animation to start.
	time.Sleep(150 * time.Millisecond)
	// Call Fail with a failure message.
	p.Fail("Failed")
	output := buf.String()
	if !strings.Contains(output, "Failed") {
		t.Errorf("Expected output to contain 'Failed', got %q", output)
	}
	if !strings.Contains(output, "✖") {
		t.Errorf("Expected output to contain default failure symbol '✖', got %q", output)
	}
}

func TestFailHasPrefix(t *testing.T) {
	// Force interactive mode.
	pin.ForceInteractive = true
	defer func() { pin.ForceInteractive = false }()

	var buf bytes.Buffer
	// Create a spinner with a custom writer and a prefix configuration.
	p := pin.New("Working",
		pin.WithWriter(&buf),
		pin.WithPrefix("TestPrefix"),
		pin.WithSeparator(":"),
	)
	cancel := p.Start(context.Background())
	defer cancel()
	// Let the spinner animate briefly.
	time.Sleep(150 * time.Millisecond)
	// Call Fail with a failure message.
	p.Fail("Error occurred")

	output := buf.String()
	if !strings.Contains(output, "TestPrefix") {
		t.Errorf("Expected output to contain prefix 'TestPrefix', got %q", output)
	}
	if !strings.Contains(output, "Error occurred") {
		t.Errorf("Expected output to contain failure message 'Error occurred', got %q", output)
	}
	if !strings.Contains(output, ":") {
		t.Errorf("Expected output to contain separator ':', got %q", output)
	}
}

func TestPrefixAndSeparatorColors(t *testing.T) {
	// Force interactive mode.
	pin.ForceInteractive = true
	defer func() { pin.ForceInteractive = false }()

	var buf bytes.Buffer

	// Define custom prefix, separator, and their colors.
	prefix := "MyPrefix"
	separator := ">"
	prefixColor := pin.ColorCyan
	separatorColor := pin.ColorWhite

	// Create a spinner with these custom options.
	p := pin.New("TestMessage",
		pin.WithWriter(&buf),
		pin.WithPrefix(prefix),
		pin.WithPrefixColor(prefixColor),
		pin.WithSeparator(separator),
		pin.WithSeparatorColor(separatorColor),
	)

	// Start the spinner and then invoke Fail to print the final output.
	cancel := p.Start(context.Background())
	defer cancel()
	time.Sleep(150 * time.Millisecond)
	p.Fail("Failure occurred")

	output := buf.String()

	if !strings.Contains(output, prefixColor.String()) {
		t.Errorf("Expected output to contain prefix color %q, got %q", prefixColor, output)
	}
	if !strings.Contains(output, separatorColor.String()) {
		t.Errorf("Expected output to contain separator color %q, got %q", separatorColor, output)
	}
	if !strings.Contains(output, prefix) {
		t.Errorf("Expected output to contain prefix %q, got %q", prefix, output)
	}
	if !strings.Contains(output, separator) {
		t.Errorf("Expected output to contain separator %q, got %q", separator, output)
	}
}

func TestStopDisplaysDoneSymbol(t *testing.T) {
	// Force interactive mode.
	pin.ForceInteractive = true
	defer func() { pin.ForceInteractive = false }()

	var buf bytes.Buffer
	doneSymbol := '✓'
	doneSymbolColor := pin.ColorGreen

	// Create a spinner configured with custom done symbol and done symbol color.
	p := pin.New("Processing",
		pin.WithWriter(&buf),
		pin.WithDoneSymbol(doneSymbol),
		pin.WithDoneSymbolColor(doneSymbolColor),
	)

	cancel := p.Start(context.Background())
	defer cancel()
	time.Sleep(150 * time.Millisecond)
	p.Stop("Completed")
	output := buf.String()
	if !strings.Contains(output, string(doneSymbol)) {
		t.Errorf("Expected output to contain done symbol %q, got %q", string(doneSymbol), output)
	}
	if !strings.Contains(output, "Completed") {
		t.Errorf("Expected output to contain final message 'Completed', got %q", output)
	}
	if !strings.Contains(output, doneSymbolColor.String()) {
		t.Errorf("Expected output to contain done symbol color %q, got %q", doneSymbolColor, output)
	}
}

func TestWithCustomSpinnerFrames(t *testing.T) {
	// Force interactive mode.
	pin.ForceInteractive = true
	defer func() { pin.ForceInteractive = false }()

	var buf bytes.Buffer
	// Define custom frames (e.g. a simple sequence: a, b, c).
	customFrames := []rune{'a', 'b', 'c'}

	// Create a spinner with custom frames using the new option.
	p := pin.New("CustomFrames", pin.WithWriter(&buf), pin.WithSpinnerFrames(customFrames))

	// Start the spinner to trigger the animation.
	cancel := p.Start(context.Background())
	defer cancel()
	time.Sleep(200 * time.Millisecond)
	p.Stop("Finished")

	output := buf.String()
	frameFound := false
	// Check that at least one of the custom frames appears in the captured output.
	for _, frame := range customFrames {
		if strings.Contains(output, string(frame)) {
			frameFound = true
			break
		}
	}
	if !frameFound {
		t.Errorf("Expected output to contain one of the custom spinner frames %q, got %q", customFrames, output)
	}
}

func TestFailWithCustomFailColor(t *testing.T) {
	// Force interactive mode.
	pin.ForceInteractive = true
	defer func() { pin.ForceInteractive = false }()

	var buf bytes.Buffer
	customFailColor := pin.ColorRed

	// Create a spinner with a custom failure color.
	p := pin.New("Working",
		pin.WithWriter(&buf),
		pin.WithFailColor(customFailColor),
	)
	cancel := p.Start(context.Background())
	defer cancel()
	time.Sleep(150 * time.Millisecond)
	p.Fail("Failure occurred")
	output := buf.String()
	if !strings.Contains(output, customFailColor.String()) {
		t.Errorf("Expected output to contain custom fail color %q, got %q", customFailColor, output)
	}
	if !strings.Contains(output, "Failure occurred") {
		t.Errorf("Expected output to contain failure message 'Failure occurred', got %q", output)
	}
}

func TestPositionSwitching(t *testing.T) {
	// Force interactive mode.
	pin.ForceInteractive = true
	defer func() { pin.ForceInteractive = false }()

	var bufLeft, bufRight bytes.Buffer

	// Create spinner with PositionLeft (default behavior).
	spinnerLeft := pin.New("TestPos",
		pin.WithWriter(&bufLeft),
		pin.WithPosition(pin.PositionLeft),
	)
	cancelLeft := spinnerLeft.Start(context.Background())
	time.Sleep(150 * time.Millisecond)
	spinnerLeft.Stop("Left Done")
	cancelLeft()

	// Create spinner with PositionRight.
	spinnerRight := pin.New("TestPos",
		pin.WithWriter(&bufRight),
		pin.WithPosition(pin.PositionRight),
	)
	cancelRight := spinnerRight.Start(context.Background())
	time.Sleep(150 * time.Millisecond)
	spinnerRight.Stop("Right Done")
	cancelRight()

	// The outputs should differ because the frame is placed in a different position.
	if bufLeft.String() == bufRight.String() {
		t.Errorf("Expected different outputs for left and right spinner positions, but both outputs were:\n%q", bufLeft.String())
	}
}

func TestNonInteractiveStart(t *testing.T) {
	// Use a bytes.Buffer to simulate a non-interactive writer.
	var buf bytes.Buffer

	// Create a spinner with the custom writer.
	p := pin.New("Non-interactive Message", pin.WithWriter(&buf))

	// Call Start; since buf is not *os.File (or os.Stdout), it will be treated as non-interactive.
	cancel := p.Start(context.Background())
	// Allow a short delay to ensure the message is printed.
	time.Sleep(100 * time.Millisecond)
	// Cancel the spinner.
	cancel()

	output := buf.String()
	expected := "Non-interactive Message\n"
	if output != expected {
		t.Errorf("Expected output %q, got %q", expected, output)
	}
}

func TestNonInteractiveFullMessageLogging(t *testing.T) {
	// Use a bytes.Buffer to simulate non-interactive output.
	var buf bytes.Buffer

	// Create a spinner with the custom writer (non-interactive mode).
	p := pin.New("Initial", pin.WithWriter(&buf))

	// Start the spinner.
	cancel := p.Start(context.Background())

	// Allow a short time for the initial message to be printed.
	time.Sleep(50 * time.Millisecond)

	// Update the spinner's message.
	p.UpdateMessage("Updated")

	// Allow a short time for the update to be printed.
	time.Sleep(50 * time.Millisecond)

	// Stop the spinner with a final message.
	p.Stop("Done")
	cancel()

	// Split the captured output into lines (ignoring empty lines).
	var lines []string
	for _, l := range strings.Split(buf.String(), "\n") {
		if strings.TrimSpace(l) != "" {
			lines = append(lines, l)
		}
	}

	expected := []string{"Initial", "Updated", "Done"}
	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines of output, got %d: %v", len(expected), len(lines), lines)
		return
	}
	for i, line := range expected {
		if lines[i] != line {
			t.Errorf("Line %d mismatch: expected %q, got %q", i+1, line, lines[i])
		}
	}
}

func TestWithSpinnerColorAndTextColor(t *testing.T) {
	// Force interactive mode so that the animation branch is executed.
	pin.ForceInteractive = true
	defer func() { pin.ForceInteractive = false }()

	var buf bytes.Buffer

	// Define desired spinner and text colors.
	expectedSpinnerColor := pin.ColorCyan
	expectedTextColor := pin.ColorYellow

	// Create a spinner with the custom spinner and text color options.
	s := pin.New("TestMessage",
		pin.WithWriter(&buf),
		pin.WithSpinnerColor(expectedSpinnerColor),
		pin.WithTextColor(expectedTextColor),
	)

	// Start the spinner to trigger the animation loop.
	cancel := s.Start(context.Background())
	// Allow some time for multiple animation ticks.
	time.Sleep(250 * time.Millisecond)
	s.Stop("Done")
	cancel()

	output := buf.String()

	// Verify that the output contains the expected spinner color and text color.
	if !strings.Contains(output, expectedSpinnerColor.String()) {
		t.Errorf("Output does not contain the expected spinner color %q. Output: %q", expectedSpinnerColor, output)
	}

	if !strings.Contains(output, expectedTextColor.String()) {
		t.Errorf("Output does not contain the expected text color %q. Output: %q", expectedTextColor, output)
	}
}

func TestCustomFailSymbolAndColor(t *testing.T) {
	// Force interactive mode to exercise the animated branch.
	pin.ForceInteractive = true
	defer func() { pin.ForceInteractive = false }()

	var buf bytes.Buffer

	customFailSymbol := 'X'                // Custom failure symbol.
	customFailSymbolColor := pin.ColorBlue // Custom failure symbol color.

	// Create a spinner with custom fail symbol and fail symbol color.
	p := pin.New("Working",
		pin.WithWriter(&buf),
		pin.WithFailSymbol(customFailSymbol),
		pin.WithFailSymbolColor(customFailSymbolColor),
	)

	cancel := p.Start(context.Background())
	defer cancel()

	// Allow some time for the spinner to animate.
	time.Sleep(150 * time.Millisecond)
	// Trigger the failure.
	p.Fail("Operation failed")

	output := buf.String()

	// Verify the custom failure symbol appears in the output.
	if !strings.Contains(output, string(customFailSymbol)) {
		t.Errorf("Output does not contain custom fail symbol %q. Output: %q", string(customFailSymbol), output)
	}

	// Verify the custom failure symbol color code appears.
	if !strings.Contains(output, customFailSymbolColor.String()) {
		t.Errorf("Output does not contain custom fail symbol color %q. Output: %q", customFailSymbolColor.String(), output)
	}

	// Also, verify that the failure message is present.
	if !strings.Contains(output, "Operation failed") {
		t.Errorf("Output does not contain failure message 'Operation failed'. Output: %q", output)
	}
}
