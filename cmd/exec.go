package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func ExecCmd(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var outputBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &outputBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("start fail %s %s %s %+v", err, cmd.Env, name, args)
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-ctx.Done():
		return stderrBuf.Bytes(), fmt.Errorf("kill deadline: %s %s %s",
			name, strings.Join(args, " "), stderrBuf.Bytes())
	case e := <-done:
		if e != nil {
			return stderrBuf.Bytes(), fmt.Errorf("kill wait: %s %s %s %s",
				name, strings.Join(args, " "), e.Error(), stderrBuf.Bytes())
		}
	}
	return outputBuf.Bytes(), nil
}

func WriteFileToDir(dirPath string, r io.Reader, fileID string) (string, error) {
	if _, err := os.Lstat(dirPath); nil != err {
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("MkdirAll Fail|%v", err)
		}
	}

	filePath := fmt.Sprintf("%s/%s", strings.TrimSuffix(dirPath, "/"), fileID)

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
	if _, err := os.Lstat(dir); nil != err {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("mkdir all fail %v", err)
		}
	}
	return nil
}

func WriteTempFile(tempFilePrefix string, r io.Reader) (*os.File, error) {
	file, err := os.CreateTemp(os.TempDir(), tempFilePrefix)
	if nil != err {
		return nil, err
	}
	_, err = io.Copy(file, r)
	if err != nil {
		return nil, err
	}
	return file, nil
}
