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
	cmd := exec.Command(name, args...)
	var outputBuf bytes.Buffer
	cmd.Stdout = &outputBuf
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("Start Fail|%s|%s|%s|%+v", err, cmd.Env, name, args)
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("Kill Deadline: %s %s %s", name, strings.Join(args, " "), cmd.Process.Kill().Error())
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("Kill Wait: %s %s %s", name, strings.Join(args, " "), err.Error())
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

func MkdirAll(dirPath string) error {
	if _, err := os.Lstat(dirPath); nil != err {
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("MkdirAll Fail|%v", err)
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
