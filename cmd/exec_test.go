package cmd_test

import (
	"bytes"
	"context"
	"os/exec"
	"testing"

	"github.com/uaxe/infra/cmd"
)

func TestExecCmd(t *testing.T) {
	c := cmd.NewExecCmd(
		context.Background(),
		cmd.WithName("go"),
		cmd.WithArgs("version"),
	)
	stdout, stderr, err := c.Run()
	if err != nil {
		t.Logf("%v, %v", err, string(stderr))
		t.FailNow()
	}
	if len(stderr) > 0 {
		t.Logf("stderr should empty")
		t.FailNow()
	}
	if len(stdout) == 0 {
		t.Logf("stdout is empty")
		t.FailNow()
	}
}

func TestExecCmd_RunWithFunc(t *testing.T) {

	c := cmd.NewExecCmd(
		context.Background(),
		cmd.WithName("go"),
		cmd.WithArgs("version"),
	)

	var stdout, stderr bytes.Buffer
	err := c.RunWithFunc(func(e *exec.Cmd) {
		e.Stdout = &stdout
		e.Stderr = &stderr
	})

	if err != nil {
		t.Logf("%v, %v", err, stderr.String())
		t.FailNow()
	}

	if len(stderr.String()) > 0 {
		t.Logf("stderr should empty")
		t.FailNow()
	}

	if len(stdout.String()) == 0 {
		t.Logf("stdout is empty")
		t.FailNow()
	}

}
