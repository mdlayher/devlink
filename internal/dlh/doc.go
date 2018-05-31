// Package dlh is an auto-generated package which contains constants and
// types used to access devlink information using generic netlink.
package dlh

// Pull the latest devlink.h from the kernel for code generation.
//go:generate wget https://raw.githubusercontent.com/torvalds/linux/master/include/uapi/linux/devlink.h

// Generate Go source from C constants.
//go:generate c-for-go -out ../ -nocgo dlh.yml

// Clean up build artifacts.
//go:generate rm -rf devlink.h _obj/
