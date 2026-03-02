package androidtouch

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Device struct {
	Serial  string
	ctx     context.Context
	timeout time.Duration
	adbPath string
}

func (d *Device) NewSession() (*Session, error) {
	s := &Session{
		device: d,
	}

	if err := s.start(); err != nil {
		return nil, err
	}

	return s, nil
}

func NewDevice(serial string) (*Device, error) {
	bin, err := exec.LookPath("adb")
	if err != nil {
		return nil, fmt.Errorf("adb not found in PATH")
	}

	d := &Device{
		Serial:  serial,
		ctx:     context.Background(),
		timeout: 3 * time.Second,
		adbPath: bin,
	}

	state, err := d.Run("get-state")
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(state, "device") {
		return nil, fmt.Errorf("device not ready, state: %q", state)
	}

	return d, nil
}

func (d *Device) adbArgs(args ...string) []string {
	if d.Serial != "" {
		return append([]string{"-s", d.Serial}, args...)
	}
	return args
}

func (d *Device) Run(args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(d.ctx, d.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, d.adbPath, d.adbArgs(args...)...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

func (d *Device) WithContext(ctx context.Context) *Device {
	return &Device{
		Serial:  d.Serial,
		ctx:     ctx,
		timeout: d.timeout,
	}
}

func (d *Device) SetTimeout(t time.Duration) {
	d.timeout = t
}

func (d *Device) Screencap() ([]byte, error) {
	ctx, cancel := context.WithTimeout(d.ctx, d.timeout)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		d.adbPath,
		d.adbArgs("exec-out", "screencap", "-p")...,
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

func (d *Device) ScreencapToFile(path string) error {
	data, err := d.Screencap()
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}


func (d *Device) WaitForStableFrame(
	threshold float64,
	checkInterval time.Duration,
	timeout time.Duration,
) ([]byte, error) {

	deadline := time.Now().Add(timeout)

	var prev []byte

	for {
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for stable frame")
		}

		current, err := d.Screencap()
		if err != nil {
			return nil, err
		}

		if prev != nil {
			diff := frameDiffPercent(prev, current)
			if diff <= threshold {
				return current, nil
			}
		}

		prev = current
		time.Sleep(checkInterval)
	}
}

func frameDiffPercent(a, b []byte) float64 {
	if len(a) != len(b) {
		return 100.0
	}

	var diff int
	for i := range a {
		if a[i] != b[i] {
			diff++
		}
	}

	return float64(diff) * 100 / float64(len(a))
}