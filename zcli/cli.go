package zcli

import (
	"fmt"
	"os"
)

func NewCli(name, description, version string) *Cli {
	result := &Cli{
		version:        version,
		bannerFunction: defaultBannerFunction,
	}
	result.rootCommand = NewCommand(name, description)
	result.rootCommand.setApp(result)
	result.rootCommand.setParentCommandPath("")
	return result
}

func defaultBannerFunction(c *Cli) string {
	version := ""
	if len(c.Version()) > 0 {
		version = " " + c.Version()
	}
	return fmt.Sprintf("%s%s - %s", c.Name(), version, c.ShortDescription())
}

type Cli struct {
	version        string
	rootCommand    *Command
	defaultCommand *Command
	preRunCommand  func(*Cli) error
	bannerFunction func(*Cli) string
	errorHandler   func(string, error) error
}

type Action func() error

func (c *Cli) ShortDescription() string {
	return c.rootCommand.shortdescription
}

func (c *Cli) Version() string {
	return c.version
}

func (c *Cli) Name() string {
	return c.rootCommand.name
}

func (c *Cli) SetBannerFunction(fn func(*Cli) string) {
	c.bannerFunction = fn
}

func (c *Cli) SetErrorFunction(fn func(string, error) error) {
	c.errorHandler = fn
}

func (c *Cli) AddCommand(command *Command) {
	c.rootCommand.AddCommand(command)
}

func (c *Cli) PrintBanner() {
	fmt.Println(c.bannerFunction(c))
	fmt.Println("")
}

func (c *Cli) PrintHelp() {
	c.rootCommand.PrintHelp()
}

func (c *Cli) Run(args ...string) error {
	if c.preRunCommand != nil {
		err := c.preRunCommand(c)
		if err != nil {
			return err
		}
	}
	if len(args) == 0 {
		args = os.Args[1:]
	}
	return c.rootCommand.run(args)
}

func (c *Cli) DefaultCommand(defaultCommand *Command) *Cli {
	c.defaultCommand = defaultCommand
	return c
}

func (c *Cli) NewSubCommand(name, description string) *Command {
	return c.rootCommand.NewSubCommand(name, description)
}

func (c *Cli) NewSubCommandInheritFlags(name, description string) *Command {
	return c.rootCommand.NewSubCommandInheritFlags(name, description)
}

func (c *Cli) PreRun(callback func(*Cli) error) {
	c.preRunCommand = callback
}

func (c *Cli) BoolFlag(name, description string, variable *bool) *Cli {
	c.rootCommand.BoolFlag(name, description, variable)
	return c
}

func (c *Cli) StringFlag(name, description string, variable *string) *Cli {
	c.rootCommand.StringFlag(name, description, variable)
	return c
}

func (c *Cli) IntFlag(name, description string, variable *int) *Cli {
	c.rootCommand.IntFlag(name, description, variable)
	return c
}

func (c *Cli) AddFlags(flags any) *Cli {
	c.rootCommand.AddFlags(flags)
	return c
}

func (c *Cli) Action(callback Action) *Cli {
	c.rootCommand.Action(callback)
	return c
}

func (c *Cli) LongDescription(longdescription string) *Cli {
	c.rootCommand.LongDescription(longdescription)
	return c
}

func (c *Cli) OtherArgs() []string {
	return c.rootCommand.flags.Args()
}

func (c *Cli) NewSubCommandFunction(name string, description string, test any) *Cli {
	c.rootCommand.NewSubCommandFunction(name, description, test)
	return c
}

func (c *Cli) CommandExists(name string) (*Command, bool) {
	result, exists := c.rootCommand.CommandExists(name)
	return result, exists
}
