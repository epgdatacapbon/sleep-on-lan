package main

import (
	"bytes"
	"fmt"
)

type version struct {
	ApplicationName     string
	Major, Minor, Patch int
}

// Version string
var Version = version{"Sleep On Lan", 1, 0, 5}

func (v version) Version() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch))
	return buf.String()
}

func (v version) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s Version %d.%d.%d", v.ApplicationName, v.Major, v.Minor, v.Patch))
	return buf.String()
}
