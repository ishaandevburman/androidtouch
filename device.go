package androidtouch

import (
	"bytes"
	"context"
	"os/exec"
)

type Device struct {
	Serial string // optional
}

func NewDevice(serial string) *Device {
	return &Device{Serial: serial}
}

func (d *Device) adbArgs(args ...string) []string {
	if d.Serial != "" {
		return append([]string{"-s", d.Serial}, args...)
	}
	return args
}

func (d *Device) Run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "adb", d.adbArgs(args...)...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	return out.String(), err
}