package util

import "os/exec"

func LookPath(name string) string {
	path, err := exec.LookPath(name)
	if err != nil {
		return name
	}
	return path
}
