// Package httpbin implements the "httpbin" subcommand, which runs a local HTTP server for testing purposes.
package httpbin

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mccutchen/go-httpbin/v2/httpbin"
	"github.com/spf13/cobra"
	"go.k6.io/k6/cmd/state"
)

//go:embed help.md
var help string

const (
	readHeaderTimeout = 2 * time.Second
	readTimeout       = 5 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 15 * time.Second
	shutdownTimeout   = 5 * time.Second
)

// newSubcommand creates a new "httpbin" subcommand.
func newSubcommand(gs *state.GlobalState) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "httpbin",
		Short: "A HTTP server for testing purposes",
		Long:  help,
	}

	flags := cmd.Flags()

	address := flags.StringP("bind", "b", "localhost:5454", "Address for the HTTP server (host:port)")

	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		mux := http.NewServeMux()

		mux.Handle("/", httpbin.New().Handler())

		server := &http.Server{
			Addr:              *address,
			Handler:           mux,
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
		}

		gs.Logger.WithField("address", *address).Info("Starting httpbin server")
		gs.Logger.Infof("Visit http://%s for available endpoints", *address)

		errChan := make(chan error, 1)

		go func() {
			err := server.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				errChan <- fmt.Errorf("server failed: %w", err)
			}
		}()

		select {
		case err := <-errChan:
			return err
		case <-cmd.Context().Done():
			gs.Logger.Info("Shutting down httpbin server...")

			shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()

			err := server.Shutdown(shutdownCtx)
			if err != nil {
				gs.Logger.WithError(err).Warn("Server shutdown encountered an error")
			}

			gs.Logger.Info("Server stopped")

			return nil
		}
	}

	return cmd
}
