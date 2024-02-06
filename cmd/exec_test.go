package cmd_test

import (
	"context"
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
