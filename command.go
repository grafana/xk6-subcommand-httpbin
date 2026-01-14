// Package httpbin implements the "httpbin" subcommand, which runs a local HTTP server for testing purposes.
package httpbin

import (
	_ "embed"
	"net/http"
	"time"

	"github.com/mccutchen/go-httpbin/v2/httpbin"
	"github.com/spf13/cobra"
	"go.k6.io/k6/cmd/state"
	"go.k6.io/k6/subcommand"
)

func init() {
	subcommand.RegisterExtension("httpbin", newCommand)
}

//go:embed help.txt
var help string

func newCommand(gs *state.GlobalState) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "httpbin",
		Short: "A HTTP server for testing purposes",
		Long:  help,
	}

	flags := cmd.Flags()

	address := flags.StringP("bind", "b", "localhost:5454", "Address for the HTTP server")

	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		mux := http.NewServeMux()

		mux.Handle("/", httpbin.New().Handler())

		server := &http.Server{
			Addr:              *address,
			Handler:           mux,
			ReadHeaderTimeout: 2 * time.Second,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       15 * time.Second,
		}

		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				gs.Logger.Fatalf("Server ListenAndServe failed: %v", err)
			}
		}()

		<-cmd.Context().Done()

		return server.Shutdown(cmd.Context())
	}

	return cmd
}
