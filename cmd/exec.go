package cmd

import (
	"bytes"
	"context"
	"os/exec"
)

type (
	ExecCmd struct {
		*option
	}

	option struct {
		ctx  context.Context
		name string
		dir  string
		args []string
		env  []string
	}

	OptionFunc func(opt *option)
)

func WithName(name string) OptionFunc {
	return func(opt *option) {
		opt.name = name
	}
}

func WithDir(dir string) OptionFunc {
	return func(opt *option) {
		opt.dir = dir
	}
}

func WithEnv(env ...string) OptionFunc {
	return func(opt *option) {
		opt.env = env
	}
}

func WithArgs(args ...string) OptionFunc {
	return func(opt *option) {
		opt.args = args
	}
}

func defaultOption(ctx context.Context) *option {
	return &option{
		ctx:  ctx,
		name: "",
		dir:  "",
		args: []string{},
		env:  []string{},
	}
}

func NewExecCmd(ctx context.Context, opts ...OptionFunc) *ExecCmd {
	c := &ExecCmd{option: defaultOption(ctx)}
	for i := range opts {
		opts[i](c.option)
	}
	return c
}

func (c *ExecCmd) Run() ([]byte, []byte, error) {
	cmd := exec.CommandContext(c.ctx, c.name, c.args...)
	cmd.Dir = c.dir
	cmd.Env = c.env
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

func (c *ExecCmd) CombinedOutput() ([]byte, error) {
	cmd := exec.CommandContext(c.ctx, c.name, c.args...)
	cmd.Dir = c.dir
	cmd.Env = c.env
	return cmd.CombinedOutput()
}

func (c *ExecCmd) RunWithFunc(f func(*exec.Cmd)) error {
	cmd := exec.CommandContext(c.ctx, c.name, c.args...)
	cmd.Dir = c.dir
	cmd.Env = c.env
	f(cmd)
	return cmd.Run()
}

func (c *ExecCmd) Cmd() *exec.Cmd {
	cmd := exec.CommandContext(c.ctx, c.name, c.args...)
	cmd.Dir = c.dir
	cmd.Env = c.env
	return cmd
}
