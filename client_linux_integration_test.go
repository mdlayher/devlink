//+build linux

package devlink_test

import (
	"testing"

	"github.com/mdlayher/devlink"
)

func TestLinuxClientIntegration(t *testing.T) {
	c, err := devlink.New()
	if err != nil {
		t.Fatalf("failed to open client: %v", err)
	}
	defer c.Close()

	// TODO(mdlayher): expand upon this.

	t.Run("devices", func(t *testing.T) {
		if _, err := c.Devices(); err != nil {
			t.Fatalf("failed to get devices: %v", err)
		}
	})

	t.Run("ports", func(t *testing.T) {
		if _, err := c.Ports(); err != nil {
			t.Fatalf("failed to get devices: %v", err)
		}
	})
}
