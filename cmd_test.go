package httpbin //nolint:testpackage

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus/hooks/test"
	"go.k6.io/k6/cmd/state"
)

// getAvailablePort returns an available port on localhost.
func getAvailablePort(tb testing.TB) (int, error) {
	tb.Helper()

	lc := net.ListenConfig{}

	listener, err := lc.Listen(context.Background(), "tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = listener.Close()
	}()

	addr := listener.Addr().(*net.TCPAddr) //nolint:forcetypeassert

	return addr.Port, nil
}

// newTestGlobalState creates a GlobalState suitable for testing with suppressed logging.
func newTestGlobalState() *state.GlobalState {
	logger, _ := test.NewNullLogger()

	return &state.GlobalState{
		Logger: logger,
	}
}

// TestNewSubcommand tests the basic subcommand creation.
func TestNewSubcommand(t *testing.T) {
	t.Parallel()

	gs := newTestGlobalState()
	cmd := newSubcommand(gs)

	if cmd == nil {
		t.Fatal("NewSubcommand returned nil")
	}

	if cmd.Use != "httpbin" {
		t.Errorf("expected Use to be 'httpbin', got '%s'", cmd.Use)
	}

	if cmd.Short != "A HTTP server for testing purposes" {
		t.Errorf("unexpected Short description: %s", cmd.Short)
	}

	if cmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

// TestCommandFlags tests that the command has the expected flags.
func TestCommandFlags(t *testing.T) {
	t.Parallel()

	gs := newTestGlobalState()
	cmd := newSubcommand(gs)

	// Check bind flag exists
	bindFlag := cmd.Flags().Lookup("bind")
	if bindFlag == nil {
		t.Fatal("bind flag not found")
	}

	if bindFlag.Shorthand != "b" {
		t.Errorf("expected bind flag shorthand to be 'b', got '%s'", bindFlag.Shorthand)
	}

	if bindFlag.DefValue != "localhost:5454" {
		t.Errorf("expected default bind address to be 'localhost:5454', got '%s'", bindFlag.DefValue)
	}
}

// TestServerStartsAndStops tests that the server can start and stop gracefully.
func TestServerStartsAndStops(t *testing.T) {
	t.Parallel()

	gs := newTestGlobalState()
	cmd := newSubcommand(gs)

	// Use a random available port
	port := "localhost:0" // 0 means random available port

	err := cmd.Flags().Set("bind", port)
	if err != nil {
		t.Fatalf("failed to set bind flag: %v", err)
	}

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Run the command in a goroutine
	errChan := make(chan error, 1)

	go func() {
		errChan <- cmd.ExecuteContext(ctx)
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Cancel the context to trigger shutdown
	cancel()

	// Wait for command to finish
	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("command returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server did not stop within timeout")
	}
}

// TestServerRespondsToRequests tests that the server actually serves httpbin endpoints.
//
//nolint:funlen,tparallel
func TestServerRespondsToRequests(t *testing.T) {
	t.Parallel()

	gs := newTestGlobalState()
	cmd := newSubcommand(gs)

	// Get an available port
	port, err := getAvailablePort(t)
	if err != nil {
		t.Fatalf("failed to get available port: %v", err)
	}

	testPort := fmt.Sprintf("localhost:%d", port)

	err = cmd.Flags().Set("bind", testPort)
	if err != nil {
		t.Fatalf("failed to set bind flag: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(cancel)

	// Run the command in a goroutine
	errChan := make(chan error, 1)

	go func() {
		errChan <- cmd.ExecuteContext(ctx)
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Test various endpoints
	tests := []struct {
		name           string
		endpoint       string
		method         string
		expectedStatus int
	}{
		{"Root", "/", "GET", http.StatusOK},
		{"GET endpoint", "/get", "GET", http.StatusOK},
		{"IP endpoint", "/ip", "GET", http.StatusOK},
		{"Headers endpoint", "/headers", "GET", http.StatusOK},
		{"User-Agent endpoint", "/user-agent", "GET", http.StatusOK},
		{"Status 200", "/status/200", "GET", http.StatusOK},
		{"Status 404", "/status/404", "GET", http.StatusNotFound},
		{"POST endpoint", "/post", "POST", http.StatusOK},
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("http://%s%s", testPort, tt.endpoint)

			req, err := http.NewRequestWithContext(ctx, tt.method, url, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}

			t.Cleanup(func() {
				_ = resp.Body.Close()
			})

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}

	// Shutdown
	cancel()

	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("command returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server did not stop within timeout")
	}
}

// TestServerHandlesInvalidAddress tests error handling for invalid addresses.
func TestServerHandlesInvalidAddress(t *testing.T) {
	t.Parallel()

	gs := newTestGlobalState()
	cmd := newSubcommand(gs)

	// Use an invalid address
	err := cmd.Flags().Set("bind", "invalid:address:format")
	if err != nil {
		t.Fatalf("failed to set bind flag: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Run the command - should fail quickly
	err = cmd.ExecuteContext(ctx)
	if err == nil {
		t.Error("expected error for invalid address, got nil")
	}
}

// TestEmbeddedHelp tests that the embedded help text is loaded.
func TestEmbeddedHelp(t *testing.T) {
	t.Parallel()

	if help == "" {
		t.Error("embedded help string is empty")
	}

	// Check for some expected content in help
	expectedStrings := []string{
		"httpbin",
		"endpoint",
		"HTTP",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(help, expected) {
			t.Errorf("help text does not contain expected string: %s", expected)
		}
	}
}

// TestMultipleInstances tests that multiple instances can't bind to the same port.
func TestMultipleInstances(t *testing.T) {
	t.Parallel()

	// Get an available port to use for both instances
	port, err := getAvailablePort(t)
	if err != nil {
		t.Fatalf("failed to get available port: %v", err)
	}

	testPort := fmt.Sprintf("localhost:%d", port)

	// Start first instance
	gs1 := newTestGlobalState()
	cmd1 := newSubcommand(gs1)

	err = cmd1.Flags().Set("bind", testPort)
	if err != nil {
		t.Fatalf("failed to set bind flag: %v", err)
	}

	ctx1, cancel1 := context.WithCancel(context.Background())

	defer cancel1()

	errChan1 := make(chan error, 1)

	go func() {
		errChan1 <- cmd1.ExecuteContext(ctx1)
	}()

	// Wait for first server to start
	time.Sleep(200 * time.Millisecond)

	// Try to start second instance on same port
	gs2 := newTestGlobalState()
	cmd2 := newSubcommand(gs2)

	err = cmd2.Flags().Set("bind", testPort)
	if err != nil {
		t.Fatalf("failed to set bind flag: %v", err)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel2()

	err = cmd2.ExecuteContext(ctx2)
	if err == nil {
		t.Error("expected error when binding to already-used port, got nil")
	}

	// Cleanup first server.
	cancel1()
	<-errChan1
}

// BenchmarkServerStartStop benchmarks the server start/stop cycle.
func BenchmarkServerStartStop(b *testing.B) {
	logger, _ := test.NewNullLogger()

	for b.Loop() {
		gs := &state.GlobalState{Logger: logger}
		cmd := newSubcommand(gs)

		// Get an available port for each iteration.
		port, err := getAvailablePort(b)
		if err != nil {
			b.Fatalf("failed to get available port: %v", err)
		}

		portStr := fmt.Sprintf("localhost:%d", port)

		err = cmd.Flags().Set("bind", portStr)
		if err != nil {
			b.Fatalf("failed to set bind flag: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())

		errChan := make(chan error, 1)

		go func() {
			errChan <- cmd.ExecuteContext(ctx)
		}()

		time.Sleep(50 * time.Millisecond)
		cancel()

		<-errChan
	}
}
