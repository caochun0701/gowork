// +build !darwin,!linux,!windows darwin,!cgo linux,!cgo windows,!cgo

package instance

import "libbeat/logp"

func setupMetrics(name string) error {
	logp.Warn("Metrics not implemented for this OS.")
	return nil
}
