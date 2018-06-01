//+build linux

package devlink_test

import (
	"os"
	"testing"

	"github.com/mdlayher/devlink"
)

func TestLinuxClientIntegration(t *testing.T) {
	c, err := devlink.New()
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("skipping, devlink is not available on this system")
		}

		t.Fatalf("failed to open client: %v", err)
	}
	defer c.Close()

	// TODO(mdlayher): expand upon this.

	t.Run("devices", func(t *testing.T) {
		devices, err := c.Devices()
		if err != nil {
			t.Fatalf("failed to get devices: %v", err)
		}

		for _, d := range devices {
			t.Logf("device: %+v", d)
		}
	})

	t.Run("ports", func(t *testing.T) {
		ports, err := c.Ports()
		if err != nil {
			t.Fatalf("failed to get ports: %v", err)
		}

		for _, p := range ports {
			t.Logf("port: %+v", p)
		}
	})
}
