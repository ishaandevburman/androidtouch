package androidtouch

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
	"sync"
)

const sentinel = "__ANDROIDTOUCH_DONE__"

type Session struct {
	device *Device

	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader

	mu sync.Mutex
}

func (s *Session) start() error {
	cmd := exec.Command(s.device.adbPath, s.device.adbArgs("shell")...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// IMPORTANT: drain stderr to avoid deadlock
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Drain stderr in background
	go io.Copy(io.Discard, stderrPipe)

	s.cmd = cmd
	s.stdin = stdin
	s.stdout = bufio.NewReader(stdoutPipe)

	return nil
}

func (s *Session) restart() error {
	if s.cmd != nil && s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
		_ = s.cmd.Wait()
	}
	return s.start()
}

func (s *Session) Run(command string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.write(command); err == nil {
		return nil
	}

	// Restart and retry once
	if err := s.restart(); err != nil {
		return err
	}

	return s.write(command)
}

func (s *Session) write(command string) error {
	_, err := s.stdin.Write([]byte(command + "\n"))
	return err
}

func (s *Session) RunSync(command string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.runSyncInternal(command); err == nil {
		return nil
	}

	// Restart and retry once
	if err := s.restart(); err != nil {
		return err
	}

	return s.runSyncInternal(command)
}

func (s *Session) runSyncInternal(command string) error {
	full := command + "; echo " + sentinel + "\n"

	if _, err := s.stdin.Write([]byte(full)); err != nil {
		return err
	}

	for {
		line, err := s.stdout.ReadString('\n')
		if err != nil {
			return err
		}

		if strings.Contains(line, sentinel) {
			break
		}
	}

	return nil
}