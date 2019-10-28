//+build linux

package devlink

import (
	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
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
	f, err := c.GetFamily(unix.DEVLINK_GENL_NAME)
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
			Command: unix.DEVLINK_CMD_GET,
			Version: unix.DEVLINK_GENL_VERSION,
		},
	}

	flags := netlink.Request | netlink.Dump

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
			Command: unix.DEVLINK_CMD_PORT_GET,
			Version: unix.DEVLINK_GENL_VERSION,
		},
	}

	flags := netlink.Request | netlink.Dump

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
		return nil, nil
	}

	ps := make([]*Port, 0, len(msgs))
	for _, m := range msgs {
		ad, err := netlink.NewAttributeDecoder(m.Data)
		if err != nil {
			return nil, err
		}

		var p Port
		for ad.Next() {
			switch ad.Type() {
			case unix.DEVLINK_ATTR_BUS_NAME:
				p.Bus = ad.String()
			case unix.DEVLINK_ATTR_DEV_NAME:
				p.Device = ad.String()
			case unix.DEVLINK_ATTR_PORT_INDEX:
				p.Port = int(ad.Uint32())
			case unix.DEVLINK_ATTR_PORT_TYPE:
				p.Type = PortType(ad.Uint16())
			// Allow netdev/ibdev name to share the same "Name" field.
			case unix.DEVLINK_ATTR_PORT_NETDEV_NAME, unix.DEVLINK_ATTR_PORT_IBDEV_NAME:
				p.Name = ad.String()
			}
		}

		if err := ad.Err(); err != nil {
			return nil, err
		}

		ps = append(ps, &p)
	}

	return ps, nil
}
