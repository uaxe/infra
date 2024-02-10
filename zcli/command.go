package zcli

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type Command struct {
	name              string
	commandPath       string
	shortdescription  string
	longdescription   string
	app               *Cli
	subCommands       []*Command
	subCommandsMap    map[string]*Command
	longestSubcommand int
	action            Action
	flags             *flag.FlagSet
	flagCount         int
	helpFlag          bool
	hidden            bool
}

func NewCommand(name string, description string) *Command {
	result := &Command{
		name:             name,
		shortdescription: description,
		subCommandsMap:   make(map[string]*Command),
		hidden:           false,
	}
	return result
}

func (c *Command) NewSubCommand(name, description string) *Command {
	result := NewCommand(name, description)
	c.AddCommand(result)
	return result
}

func (c *Command) Action(callback Action) *Command {
	c.action = callback
	return c
}

func (c *Command) Hidden() {
	c.hidden = true
}

func (c *Command) parseFlags(args []string) error {
	tmp := os.Stderr
	os.Stderr = nil
	err := c.flags.Parse(args)
	os.Stderr = tmp
	return err
}

func (c *Command) isHidden() bool {
	return c.hidden
}

func (c *Command) isDefaultCommand() bool {
	return c.app.defaultCommand == c
}

func (c *Command) PrintHelp() {
	c.app.PrintBanner()

	commandTitle := c.commandPath
	if c.shortdescription != "" {
		commandTitle += " - " + c.shortdescription
	}

	if c.commandPath != c.name {
		fmt.Println(commandTitle)
	}
	if c.longdescription != "" {
		fmt.Println(c.longdescription + "\n")
	}
	if len(c.subCommands) > 0 {
		fmt.Println("Available commands:")
		fmt.Println("")
		for _, subcommand := range c.subCommands {
			if subcommand.isHidden() {
				continue
			}
			spacer := strings.Repeat(" ", 3+c.longestSubcommand-len(subcommand.name))
			isDefault := ""
			if subcommand.isDefaultCommand() {
				isDefault = "[default]"
			}
			fmt.Printf("   %s%s%s %s\n", subcommand.name, spacer, subcommand.shortdescription, isDefault)
		}
		fmt.Println("")
	}
	if c.flagCount > 0 {
		fmt.Println("Flags:")
		fmt.Println()
		c.flags.SetOutput(os.Stdout)
		c.flags.PrintDefaults()
		c.flags.SetOutput(os.Stderr)
	}
	fmt.Println()
}

func (c *Command) run(args []string) error {

	if len(args) > 0 {

		subcommand := c.subCommandsMap[args[0]]
		if subcommand != nil {
			return subcommand.run(args[1:])
		}

		err := c.parseFlags(args)
		if err != nil {
			if c.app.errorHandler != nil {
				return c.app.errorHandler(c.commandPath, err)
			}
			return fmt.Errorf("Error: %w\nSee '%s --help' for usage", err, c.commandPath)
		}

		if c.helpFlag {
			c.PrintHelp()
			return nil
		}
	}

	if c.action != nil {
		return c.action()
	}

	if c.app.defaultCommand != nil {
		if c.app.defaultCommand != c {
			if len(args) == 0 {
				return c.app.defaultCommand.run(args)
			}
		}
	}

	c.PrintHelp()

	return nil
}

func (c *Command) setApp(app *Cli) {
	c.app = app
}

func (c *Command) setParentCommandPath(parentCommandPath string) {
	if parentCommandPath != "" {
		c.commandPath += parentCommandPath + " "
	}
	c.commandPath += c.name

	c.flags = flag.NewFlagSet(c.commandPath, flag.ContinueOnError)
	c.BoolFlag("help", "Get help on the '"+strings.ToLower(c.commandPath)+"' command.", &c.helpFlag)
}

func (c *Command) AddCommand(command *Command) {
	command.setApp(c.app)
	command.setParentCommandPath(c.commandPath)
	name := command.name
	c.subCommands = append(c.subCommands, command)
	c.subCommandsMap[name] = command
	if len(name) > c.longestSubcommand {
		c.longestSubcommand = len(name)
	}
}

func (c *Command) BoolFlag(name, description string, variable *bool) *Command {
	c.flags.BoolVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

func (c *Command) StringFlag(name, description string, variable *string) *Command {
	c.flags.StringVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

func (c *Command) IntFlag(name, description string, variable *int) *Command {
	c.flags.IntVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

func (c *Command) AddFlags(optionStruct any) *Command {

	t := reflect.TypeOf(optionStruct)

	if t.Kind() != reflect.Ptr {
		panic("AddFlags() requires a pointer to a struct")
	}
	if t.Elem().Kind() != reflect.Struct {
		panic("AddFlags() requires a pointer to a struct")
	}

	v := reflect.ValueOf(optionStruct).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Elem().Field(i)
		if !fieldType.IsExported() {
			continue
		}

		if fieldType.Type.Kind() == reflect.Struct {
			c.AddFlags(field.Addr().Interface())
			continue
		}

		tag := t.Elem().Field(i).Tag
		name := tag.Get("name")
		description := tag.Get("description")
		if name == "" {
			name = strings.ToLower(t.Elem().Field(i).Name)
		}
		switch field.Kind() {
		case reflect.Bool:
			c.BoolFlag(name, description, field.Addr().Interface().(*bool))
		case reflect.String:
			c.StringFlag(name, description, field.Addr().Interface().(*string))
		case reflect.Int:
			c.IntFlag(name, description, field.Addr().Interface().(*int))
		default:
		}
	}

	return c
}

func (c *Command) OtherArgs() []string {
	return c.flags.Args()
}

func (c *Command) LongDescription(longdescription string) *Command {
	c.longdescription = longdescription
	return c
}

func (c *Command) NewSubCommandInheritFlags(name, description string) *Command {
	result := c.NewSubCommand(name, description)
	result.inheritFlags(c.flags)
	return result
}

func (c *Command) inheritFlags(inheritFlags *flag.FlagSet) {
	inheritFlags.VisitAll(func(f *flag.Flag) {
		if f.Name != "help" {
			c.flags.Var(f.Value, f.Name, f.Usage)
		}
	})
}

func (c *Command) NewSubCommandFunction(name string, description string, fn any) *Command {
	result := c.NewSubCommand(name, description)
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		panic("NewSubFunction '" + name + "' requires a function with the signature 'func(*struct) error'")
	}

	fnValue := reflect.ValueOf(fn)
	if t.NumIn() != 1 {
		panic("NewSubFunction '" + name + "' requires a function with the signature 'func(*struct) error'")
	}

	if t.In(0).Kind() != reflect.Ptr {
		panic("NewSubFunction '" + name + "' requires a function with the signature 'func(*struct) error'")
	}
	if t.In(0).Elem().Kind() != reflect.Struct {
		panic("NewSubFunction '" + name + "' requires a function with the signature 'func(*struct) error'")
	}
	if t.NumOut() != 1 {
		panic("NewSubFunction '" + name + "' requires a function with the signature 'func(*struct) error'")
	}
	if t.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		panic("NewSubFunction '" + name + "' requires a function with the signature 'func(*struct) error'")
	}
	flags := reflect.New(t.In(0).Elem())
	defaultMethod, ok := t.In(0).MethodByName("Default")

	if ok {
		if defaultMethod.Type.NumIn() != 1 {
			panic("'Default' method on struct '" + t.In(0).Elem().Name() + "' must have the signature 'Default() *" + t.In(0).Elem().Name() + "'")
		}

		if defaultMethod.Type.NumOut() != 1 {
			panic("'Default' method on struct '" + t.In(0).Elem().Name() + "' must have the signature 'Default() *" + t.In(0).Elem().Name() + "'")
		}

		if defaultMethod.Type.Out(0) != t.In(0) {
			panic("'Default' method on struct '" + t.In(0).Elem().Name() + "' must have the signature 'Default() *" + t.In(0).Elem().Name() + "'")
		}

		results := defaultMethod.Func.Call([]reflect.Value{flags})
		flags = results[0]
	}
	result.Action(func() error {
		ret := fnValue.Call([]reflect.Value{flags})[0].Interface()
		if ret != nil {
			return ret.(error)
		}
		return nil
	})
	result.AddFlags(flags.Interface())
	return result
}

func (c *Command) CommandExists(name string) (*Command, bool) {
	result, exists := c.subCommandsMap[name]
	return result, exists
}
