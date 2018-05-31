//+build !linux

package devlink

import (
	"fmt"
	"runtime"
)

var (
	// errUnimplemented is returned by all functions on platforms that
	// cannot make use of devlink.
	errUnimplemented = fmt.Errorf("devlink not implemented on %s/%s",
		runtime.GOOS, runtime.GOARCH)
)

var _ osClient = &client{}

// A client is an unimplemented devlink client.
type client struct{}

// newClient always returns an error.
func newClient() (*client, error) {
	return nil, errUnimplemented
}

// Close implements osClient.
func (c *client) Close() error {
	return errUnimplemented
}

// CGroupStats implements osClient.
func (c *client) Devices() ([]*Device, error) {
	return nil, errUnimplemented
}

// PID implements osClient.
func (c *client) Ports() ([]*Port, error) {
	return nil, errUnimplemented
}
