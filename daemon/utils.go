package daemon

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/gluster/glusterd2/errors"
)

func WritePidToFile(pid int, path string) error {
	pidFileDir := filepath.Dir(path)

	err := os.MkdirAll(pidFileDir, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	pidInBytes := []byte(strconv.Itoa(pid))
	err = ioutil.WriteFile(path, pidInBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func ReadPidFromFile(path string) (int, error) {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return -1, err
	}

	pid, err := strconv.Atoi(string(bytes.TrimSpace(content)))
	if err != nil {
		return -1, err
	}

	return pid, nil
}

func GetProcess(pid int) (*os.Process, error) {

	process, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}

	// From https://golang.org/pkg/os/#FindProcess:
	// On Unix systems, FindProcess always succeeds and returns a Process
	// for the given pid, regardless of whether the process exists.
	//
	// Sending signal 0 can be used to check for the existence of a process ID
	// Refer `man 2 kill`
	err = process.Signal(syscall.Signal(0))
	if err == nil {
		return process, nil
	}

	return nil, errors.ErrProcessNotFound
}
