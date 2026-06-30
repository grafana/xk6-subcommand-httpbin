package httpbin

import "go.k6.io/k6/v2/subcommand"

func init() {
	subcommand.RegisterExtension("httpbin", newSubcommand)
}
