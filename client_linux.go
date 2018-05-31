//+build linux

package devlink

import (
	"os"

	"github.com/mdlayher/devlink/internal/dlh"
	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nlenc"
)

var _ osClient = &client{}

// A client is a Linux-specific devlink client.
type client struct {
	c      *genetlink.Conn
	family genetlink.Family
}

// newClient opens a connection to the devlink family using generic netlink.
func newClient() (*client, error) {
	c, err := genetlink.Dial(nil)
	if err != nil {
		return nil, err
	}

	return initClient(c)
}

// initClient is the internal client constructor used in some tests.
func initClient(c *genetlink.Conn) (*client, error) {
	f, err := c.GetFamily(dlh.GenlName)
	if err != nil {
		_ = c.Close()
		return nil, err
	}

	return &client{
		c:      c,
		family: f,
	}, nil
}

// Close implements osClient.
func (c *client) Close() error {
	return c.c.Close()
}

// Devices implements osClient.
func (c *client) Devices() ([]*Device, error) {
	msg := genetlink.Message{
		Header: genetlink.Header{
			Command: dlh.CmdGet,
			Version: dlh.GenlVersion,
		},
	}

	flags := netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump

	msgs, err := c.c.Execute(msg, c.family.ID, flags)
	if err != nil {
		return nil, err
	}

	return parseDevices(msgs)
}

// Ports implements osClient.
func (c *client) Ports() ([]*Port, error) {
	msg := genetlink.Message{
		Header: genetlink.Header{
			Command: dlh.CmdPortGet,
			Version: dlh.GenlVersion,
		},
	}

	flags := netlink.HeaderFlagsRequest | netlink.HeaderFlagsDump

	msgs, err := c.c.Execute(msg, c.family.ID, flags)
	if err != nil {
		return nil, err
	}

	return parsePorts(msgs)
}

// parseDevices parses Devices from a slice of generic netlink messages.
func parseDevices(msgs []genetlink.Message) ([]*Device, error) {
	// It appears that a Device is just a subset of the attributes found in
	// a Port, so we just call the port parsing function to avoid duplication.
	ports, err := parsePorts(msgs)
	if err != nil {
		return nil, err
	}

	ds := make([]*Device, 0, len(msgs))
	for _, p := range ports {
		ds = append(ds, &Device{
			Bus:    p.Bus,
			Device: p.Device,
		})
	}

	return ds, nil
}

// parsePorts parses Ports from a slice of generic netlink messages.
func parsePorts(msgs []genetlink.Message) ([]*Port, error) {
	if len(msgs) == 0 {
		// No devlink response found.
		return nil, os.ErrNotExist
	}

	ps := make([]*Port, 0, len(msgs))
	for _, m := range msgs {
		attrs, err := netlink.UnmarshalAttributes(m.Data)
		if err != nil {
			return nil, err
		}

		var p Port
		for _, a := range attrs {
			switch a.Type {
			case dlh.AttrBusName:
				p.Bus = nlenc.String(a.Data)
			case dlh.AttrDevName:
				p.Device = nlenc.String(a.Data)
			case dlh.AttrPortIndex:
				p.Port = int(nlenc.Uint32(a.Data))
			case dlh.AttrPortType:
				p.Type = PortType(nlenc.Uint16(a.Data))
			// Allow netdev/ibdev name to share the same "Name" field.
			case dlh.AttrPortNetdevName, dlh.AttrPortIbdevName:
				p.Name = nlenc.String(a.Data)
			}
		}

		ps = append(ps, &p)
	}

	return ps, nil
}
