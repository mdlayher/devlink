// Package devlink provides access to Linux's devlink interface.
package devlink

import (
	"io"
)

// A Client provides access to Linux devlink information.
type Client struct {
	c osClient
}

// New creates a new Client.
func New() (*Client, error) {
	c, err := newClient()
	if err != nil {
		return nil, err
	}

	return &Client{
		c: c,
	}, nil
}

// Close releases resources used by a Client.
func (c *Client) Close() error {
	return c.c.Close()
}

// Devices retrieves all devlink devices on this system.
func (c *Client) Devices() ([]*Device, error) {
	return c.c.Devices()
}

// Ports retrieves all devlink ports attached to devices on this system.
func (c *Client) Ports() ([]*Port, error) {
	return c.c.Ports()
}

// An osClient is the operating system-specific implementation of Client.
type osClient interface {
	io.Closer
	Devices() ([]*Device, error)
	Ports() ([]*Port, error)
}

// A Device is a devlink device.
type Device struct {
	Bus    string
	Device string
}

//go:generate stringer -type=PortType -output=string.go

// A PortType is the operating mode of a devlink port.
type PortType int

// Possible PortType values.
const (
	Unknown PortType = iota
	Auto
	Ethernet
	InfiniBand
)

// A Port is a devlink port which is attached to a devlink device.
type Port struct {
	Bus    string
	Device string
	Port   int
	Type   PortType
	Name   string
}
