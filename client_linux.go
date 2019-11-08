//+build linux

package devlink

import (
	"fmt"

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
func (c *client) Close() error { return c.c.Close() }

// Devices implements osClient.
func (c *client) Devices() ([]*Device, error) {
	msgs, err := c.execute(unix.DEVLINK_CMD_GET, netlink.Dump, nil)
	if err != nil {
		return nil, err
	}

	return parseDevices(msgs)
}

// Ports implements osClient.
func (c *client) Ports() ([]*Port, error) {
	msgs, err := c.execute(unix.DEVLINK_CMD_PORT_GET, netlink.Dump, nil)
	if err != nil {
		return nil, err
	}

	return parsePorts(msgs)
}

// DpipeTables implements osClient.
func (c *client) DpipeTables(dev *Device) ([]*DpipeTable, error) {
	if dev == nil {
		return nil, fmt.Errorf("invalid argument")
	}
	encoder := netlink.NewAttributeEncoder()
	encoder.String(unix.DEVLINK_ATTR_BUS_NAME, dev.Bus)
	encoder.String(unix.DEVLINK_ATTR_DEV_NAME, dev.Device)
	data, err := encoder.Encode()
	if err != nil {
		return nil, err
	}

	msgs, err := c.execute(unix.DEVLINK_CMD_DPIPE_TABLE_GET, netlink.Acknowledge, data)
	if err != nil {
		return nil, err
	}

	return parseDpipeTables(msgs)
}

// execute executes the specified command with additional header flags. The
// netlink.Request header flag is automatically set.
func (c *client) execute(cmd uint8, flags netlink.HeaderFlags, data []byte) ([]genetlink.Message, error) {
	return c.c.Execute(
		genetlink.Message{
			Header: genetlink.Header{
				Command: cmd,
				Version: unix.DEVLINK_GENL_VERSION,
			},
			Data: data,
		},
		// Always pass the genetlink family ID and request flag.
		c.family.ID,
		netlink.Request|flags,
	)
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

// parseDpipeTables parses DPIPE tables from a slice of generic netlink messages.
func parseDpipeTables(msgs []genetlink.Message) ([]*DpipeTable, error) {
	var bus, dev string
	if len(msgs) == 0 {
		// No devlink response found.
		return nil, nil
	}

	ts := make([]*DpipeTable, 0, len(msgs))
	for _, m := range msgs {
		ad, err := netlink.NewAttributeDecoder(m.Data)
		if err != nil {
			return nil, err
		}

		for ad.Next() {
			switch ad.Type() {
			case unix.DEVLINK_ATTR_BUS_NAME:
				bus = ad.String()
			case unix.DEVLINK_ATTR_DEV_NAME:
				dev = ad.String()
			case unix.DEVLINK_ATTR_DPIPE_TABLES:
				tablesData := ad.Bytes()
				adTables, err := netlink.NewAttributeDecoder(tablesData)
				if err != nil {
					continue
				}
				for adTables.Next() {
					if adTables.Type() == unix.DEVLINK_ATTR_DPIPE_TABLE {
						tableData := adTables.Bytes()
						adTable, err := netlink.NewAttributeDecoder(tableData)
						if err != nil {
							continue
						}
						var t DpipeTable
						t.Bus = bus
						t.Device = dev
						for adTable.Next() {
							switch adTable.Type() {
							case unix.DEVLINK_ATTR_DPIPE_TABLE_NAME:
								t.Name = adTable.String()
							case unix.DEVLINK_ATTR_DPIPE_TABLE_SIZE:
								t.Size = adTable.Uint64()
							case unix.DEVLINK_ATTR_DPIPE_TABLE_COUNTERS_ENABLED:
								t.CountersEnabled = adTable.Uint8() != 0
							}
						}
						ts = append(ts, &t)
					}
				}
			}

		}
	}
	return ts, nil
}
