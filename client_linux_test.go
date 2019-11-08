//+build linux

package devlink

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/genetlink/genltest"
	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
	"github.com/mdlayher/netlink/nltest"
	"golang.org/x/sys/unix"
)

func TestLinuxClientEmptyResponse(t *testing.T) {
	const (
		bus    = "pci"
		device = "0000:01:00.0"
	)
	tests := []struct {
		name string
		fn   func(t *testing.T, c *client)
		msgs []genetlink.Message
	}{
		{
			name: "devices",
			fn: func(t *testing.T, c *client) {
				devices, err := c.Devices()
				if err != nil {
					t.Fatalf("failed to get devices: %v", err)
				}

				if diff := cmp.Diff(0, len(devices)); diff != "" {
					t.Fatalf("unexpected number of devices (-want +got):\n%s", diff)
				}
			},
		},
		{
			name: "ports",
			fn: func(t *testing.T, c *client) {
				ports, err := c.Ports()
				if err != nil {
					t.Fatalf("failed to get ports: %v", err)
				}

				if diff := cmp.Diff(0, len(ports)); diff != "" {
					t.Fatalf("unexpected number of ports (-want +got):\n%s", diff)
				}
			},
		},
		{
			name: "dpipe_tables",
			fn: func(t *testing.T, c *client) {
				dev := Device{bus, device}
				tables, err := c.DpipeTables(&dev)
				if err != nil {
					t.Fatalf("failed to get DPIPE tables: %v", err)
				}

				if diff := cmp.Diff(0, len(tables)); diff != "" {
					t.Fatalf("unexpected number of DPIPE tables (-want +got):\n%s", diff)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient(t, func(_ genetlink.Message, _ netlink.Message) ([]genetlink.Message, error) {
				return tt.msgs, nil
			})
			defer c.Close()

			tt.fn(t, c)
		})
	}
}

func TestLinuxClientDevicesOK(t *testing.T) {
	const (
		bus     = "pci"
		deviceA = "0000:01:00.0"
		deviceB = "0000:02:00.0"
	)

	tests := []struct {
		name    string
		msgs    []genetlink.Message
		devices []*Device
	}{
		{
			name: "one",
			msgs: []genetlink.Message{{
				Data: nltest.MustMarshalAttributes([]netlink.Attribute{
					{
						Type: unix.DEVLINK_ATTR_BUS_NAME,
						Data: nlenc.Bytes(bus),
					},
					{
						Type: unix.DEVLINK_ATTR_DEV_NAME,
						Data: nlenc.Bytes(deviceA),
					},
				}),
			}},
			devices: []*Device{{
				Bus:    bus,
				Device: deviceA,
			}},
		},
		{
			name: "multiple",
			msgs: []genetlink.Message{
				{
					Data: nltest.MustMarshalAttributes([]netlink.Attribute{
						{
							Type: unix.DEVLINK_ATTR_BUS_NAME,
							Data: nlenc.Bytes(bus),
						},
						{
							Type: unix.DEVLINK_ATTR_DEV_NAME,
							Data: nlenc.Bytes(deviceA),
						},
					}),
				},
				{
					Data: nltest.MustMarshalAttributes([]netlink.Attribute{
						{
							Type: unix.DEVLINK_ATTR_BUS_NAME,
							Data: nlenc.Bytes(bus),
						},
						{
							Type: unix.DEVLINK_ATTR_DEV_NAME,
							Data: nlenc.Bytes(deviceB),
						},
					}),
				},
			},
			devices: []*Device{
				{
					Bus:    bus,
					Device: deviceA,
				},
				{
					Bus:    bus,
					Device: deviceB,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const (
				cmd   = unix.DEVLINK_CMD_GET
				flags = netlink.Request | netlink.Dump
			)

			fn := func(_ genetlink.Message, _ netlink.Message) ([]genetlink.Message, error) {
				return tt.msgs, nil
			}

			c := testClient(t, genltest.CheckRequest(familyID, cmd, flags, fn))
			defer c.Close()

			devices, err := c.Devices()
			if err != nil {
				t.Fatalf("failed to get devices: %v", err)
			}

			if diff := cmp.Diff(tt.devices, devices); diff != "" {
				t.Fatalf("unexpected devices (-want +got):\n%s", diff)
			}
		})
	}
}

const familyID = 20

func testClient(t *testing.T, fn genltest.Func) *client {
	family := genetlink.Family{
		ID:      familyID,
		Version: unix.DEVLINK_GENL_VERSION,
		Name:    unix.DEVLINK_GENL_NAME,
	}

	conn := genltest.Dial(genltest.ServeFamily(family, fn))

	c, err := initClient(conn)
	if err != nil {
		t.Fatalf("failed to open client: %v", err)
	}

	return c
}
