package qutil

import (
	"os/exec"
	"runtime"
)

func Open(url string) error {
	var (
		cmd  string
		args []string
	)

	switch runtime.GOOS {
	case "windows":
		cmd, args = "cmd", []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		// "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
