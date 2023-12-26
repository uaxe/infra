package cmd

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func ExecCmd(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

func WriteFileToDir(dir, fname string, r io.Reader) (string, error) {
	if _, err := os.Lstat(dir); nil != err {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	filePath := filepath.Join(dir, fname)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if nil != err {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func MkdirAll(dir string) error {
	if _, err := os.Lstat(dir); err != nil {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
