// Package linux ...
package linux

import "os/exec"

// ShellMV move file by exec shell command
func ShellMV(src, dst string) error {
	cmd := exec.Command("mv", src, dst)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
