// Code generated by "stringer -type=PortType -output=string.go"; DO NOT EDIT.

package devlink

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unknown-0]
	_ = x[Auto-1]
	_ = x[Ethernet-2]
	_ = x[InfiniBand-3]
}

const _PortType_name = "UnknownAutoEthernetInfiniBand"

var _PortType_index = [...]uint8{0, 7, 11, 19, 29}

func (i PortType) String() string {
	if i < 0 || i >= PortType(len(_PortType_index)-1) {
		return "PortType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PortType_name[_PortType_index[i]:_PortType_index[i+1]]
}
