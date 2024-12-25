package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func readPID() (int, error) {
	tmp := os.TempDir()
	p := filepath.Join(tmp, pid_file)

	pbytes, err := os.ReadFile(p)
	if err != nil {
		return -1, fmt.Errorf("failed to read pid file: %w", err)
	}
	pid, err := strconv.Atoi(string(pbytes))
	if err != nil {
		return -1, fmt.Errorf("failed to parse pid: %w", err)
	}
	return pid, nil
}

func Kill() (int, error) {
	pid, err := readPID()
	if err != nil {
		return -1, err
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return -1, fmt.Errorf("failed to find process: %w", err)
	}
	err = process.Kill()
	if err != nil {
		return -1, fmt.Errorf("failed to kill process: %w", err)
	}
	return pid, nil
}
