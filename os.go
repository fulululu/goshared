package goshared

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// ShellMV move file by exec shell command (Unix/Linux only)
func ShellMV(src, dst string) error {
	cmd := exec.Command("mv", src, dst)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

// IOMV Solve problem that os.Rename() give error "invalid cross-device link"
// in Docker container with Volumes(container FS and volume FS are different).
// Note that you need have permissions on both FS.
func IOMV(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}

// NewSignalCtx generate context which work with SIGINT and SIGTERM
func NewSignalCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sCh := make(chan os.Signal, 1)
		signal.Notify(sCh, syscall.SIGINT, syscall.SIGTERM)
		<-sCh
		cancel()
	}()
	return ctx
}

// CheckFileExistence ...
func CheckFileExistence(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
