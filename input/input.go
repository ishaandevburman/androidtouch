package input

import (
	"strconv"
	"time"

	"github.com/ishaandevburman/androidtouch"
)

type Input struct {
	session *androidtouch.Session
}

func New(session *androidtouch.Session) *Input {
	return &Input{session: session}
}

func (i *Input) Tap(x, y int) error {
	return i.session.Run(
		"input tap " + strconv.Itoa(x) + " " + strconv.Itoa(y),
	)
}

func (i *Input) TapSync(x, y int) error {
	return i.session.RunSync(
		"input tap " + strconv.Itoa(x) + " " + strconv.Itoa(y),
	)
}

func (i *Input) TapWithDelay(x, y int, d time.Duration) error {
	if err := i.TapSync(x, y); err != nil {
		return err
	}
	time.Sleep(d)
	return nil
}

func (i *Input) Swipe(x1, y1, x2, y2, duration int) error {
	return i.session.Run(
		"input swipe " +
			strconv.Itoa(x1) + " " +
			strconv.Itoa(y1) + " " +
			strconv.Itoa(x2) + " " +
			strconv.Itoa(y2) + " " +
			strconv.Itoa(duration),
	)
}

func (i *Input) KeyEvent(keycode int) error {
	return i.session.Run(
		"input keyevent " + strconv.Itoa(keycode),
	)
}
