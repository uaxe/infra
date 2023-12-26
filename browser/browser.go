package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

var (
	commands = map[string]CommondArg{
		"windows": {"cmd", []string{"/c", "start"}},
		"darwin":  {"open", []string{}},
		"linux":   {"xdg-open", []string{}},
	}
)

type CommondArg struct {
	Cmd  string
	Args []string
}

func OpenBrowserURL(uri string) error {
	runArg, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}
	runArg.Args = append(runArg.Args, uri)
	cmd := exec.Command(runArg.Cmd, runArg.Args...)
	return cmd.Start()
}
