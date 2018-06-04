//+build linux

package devlink

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mdlayher/devlink/internal/dlh"
	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/genetlink/genltest"
	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
	"github.com/mdlayher/netlink/nltest"
)

func TestLinuxClientIsNotExist(t *testing.T) {
	tests := []struct {
		name string
		fn   func(c *client) error
		msgs []genetlink.Message
	}{
		{
			name: "devices",
			fn: func(c *client) error {
				_, err := c.Devices()
				return err
			},
		},
		{
			name: "ports",
			fn: func(c *client) error {
				_, err := c.Ports()
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient(t, func(_ genetlink.Message, _ netlink.Message) ([]genetlink.Message, error) {
				return tt.msgs, nil
			})
			defer c.Close()

			if err := tt.fn(c); !os.IsNotExist(err) {
				t.Fatalf("expected is not exist, but got: %v", err)
			}
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
						Type: dlh.AttrBusName,
						Data: nlenc.Bytes(bus),
					},
					{
						Type: dlh.AttrDevName,
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
							Type: dlh.AttrBusName,
							Data: nlenc.Bytes(bus),
						},
						{
							Type: dlh.AttrDevName,
							Data: nlenc.Bytes(deviceA),
						},
					}),
				},
				{
					Data: nltest.MustMarshalAttributes([]netlink.Attribute{
						{
							Type: dlh.AttrBusName,
							Data: nlenc.Bytes(bus),
						},
						{
							Type: dlh.AttrDevName,
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
				cmd   = dlh.CmdGet
				flags = netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump
			)

			fn := func(_ genetlink.Message, _ netlink.Message) ([]genetlink.Message, error) {
				return tt.msgs, nil
			}

			c := testClient(t, checkRequest(cmd, flags, fn))
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

func checkRequest(command uint8, flags netlink.HeaderFlags, fn genltest.Func) genltest.Func {
	return func(greq genetlink.Message, nreq netlink.Message) ([]genetlink.Message, error) {
		if want, got := command, greq.Header.Command; command != 0 && want != got {
			return nil, fmt.Errorf("unexpected generic netlink header command: %d, want: %d", got, want)
		}

		if want, got := flags, nreq.Header.Flags; flags != 0 && want != got {
			return nil, fmt.Errorf("unexpected netlink header flags: %s, want: %s", got, want)
		}

		return fn(greq, nreq)
	}
}

func testClient(t *testing.T, fn genltest.Func) *client {
	family := genetlink.Family{
		ID:      20,
		Version: dlh.GenlVersion,
		Name:    dlh.GenlName,
	}

	conn := genltest.Dial(genltest.ServeFamily(family, fn))

	c, err := initClient(conn)
	if err != nil {
		t.Fatalf("failed to open client: %v", err)
	}

	return c
}
