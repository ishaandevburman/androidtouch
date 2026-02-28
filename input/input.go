package input

import (
	"context"
	"fmt"

	"github.com/ishaandevburman/androidtouch"
)

type Input struct {
	device *androidtouch.Device
}

func New(device *androidtouch.Device) *Input {
	return &Input{device: device}
}

func (i *Input) Tap(ctx context.Context, x, y int) error {
	_, err := i.device.Run(ctx, "shell", "input", "tap",
		fmt.Sprintf("%d", x),
		fmt.Sprintf("%d", y),
	)
	return err
}

func (i *Input) Swipe(ctx context.Context, x1, y1, x2, y2, duration int) error {
	_, err := i.device.Run(ctx, "shell", "input", "swipe",
		fmt.Sprintf("%d", x1),
		fmt.Sprintf("%d", y1),
		fmt.Sprintf("%d", x2),
		fmt.Sprintf("%d", y2),
		fmt.Sprintf("%d", duration),
	)
	return err
}

func (i *Input) KeyEvent(ctx context.Context, keycode int) error {
	_, err := i.device.Run(ctx, "shell", "input", "keyevent",
		fmt.Sprintf("%d", keycode),
	)
	return err
}