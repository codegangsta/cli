package cli

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Command is a subcommand for a cli.App.
type Command struct {
	// The name of the command
	Name string
	// short name of the command. Typically one character
	ShortName string
	// A short description of the usage of this command
	Usage string
	// A longer explanation of how the command works
	Description string
	// The function to call when checking for bash command completions
	BashComplete func(context *Context)
	// An action to execute before any sub-subcommands are run, but after the context is ready
	// If a non-nil error is returned, no sub-subcommands are run
	Before func(context *Context) error
	// The function to call when this command is invoked
	Action func(context *Context)
	// List of child commands
	Subcommands []Command
	// List of flags to parse
	Flags []Flag
	// Treat all flags as normal arguments if true
	SkipFlagParsing bool
	// Boolean to hide built-in help command
	HideHelp bool
}

type boolFlag interface {
	IsBoolFlag() bool
}

// Invokes the command given the context, parses ctx.Args() to generate command-specific flags
func (c Command) Run(ctx *Context) error {

	if len(c.Subcommands) > 0 || c.Before != nil {
		return c.startApp(ctx)
	}

	if !c.HideHelp {
		// append help to flags
		c.Flags = append(
			c.Flags,
			HelpFlag,
		)
	}

	if ctx.App.EnableBashCompletion {
		c.Flags = append(c.Flags, BashCompletionFlag)
	}

	set := flagSet(c.Name, c.Flags)
	set.SetOutput(ioutil.Discard)

	var err error
	if !c.SkipFlagParsing {
		flagArgs := []string{}
		flagIdx := 0
		boolF := true
		for idx := 1; idx < len(ctx.Args()); idx++ {
			val := ctx.Args()[idx]
			if strings.HasPrefix(val, "-") {
				boolF = true
				flagArgs = append(flagArgs, "")
				copy(flagArgs[flagIdx+1:], flagArgs[flagIdx:])
				flagArgs[flagIdx] = val
				flagIdx++
				// uggly test for - and -- flags
				if len(val) > 1 && strings.HasPrefix(val, "--") {
					flag := set.Lookup(val[2:])
					if flag == nil {
						continue
					}
					fv, ok := flag.Value.(boolFlag)
					if !ok || !fv.IsBoolFlag() {
						boolF = false
					}
				} else {
					flag := set.Lookup(val[1:])
					if flag == nil {
						continue
					}
					fv, ok := flag.Value.(boolFlag)
					if !ok || !fv.IsBoolFlag() {
						boolF = false
					}
				}
			} else {
				if !boolF {
					flagArgs = append(flagArgs, "")
					copy(flagArgs[flagIdx+1:], flagArgs[flagIdx:])
					flagArgs[flagIdx] = val
					flagIdx++
				} else {
					flagArgs = append(flagArgs, val)
				}
				boolF = true
			}
		}
		err = set.Parse(flagArgs)
	} else {
		err = set.Parse(ctx.Args().Tail())
	}

	if err != nil {
		fmt.Printf("Incorrect Usage.\n\n")
		ShowCommandHelp(ctx, c.Name)
		fmt.Println("")
		return err
	}

	nerr := normalizeFlags(c.Flags, set)
	if nerr != nil {
		fmt.Println(nerr)
		fmt.Println("")
		ShowCommandHelp(ctx, c.Name)
		fmt.Println("")
		return nerr
	}
	context := NewContext(ctx.App, set, ctx.globalSet)

	if checkCommandCompletions(context, c.Name) {
		return nil
	}

	if checkCommandHelp(context, c.Name) {
		return nil
	}
	context.Command = c
	c.Action(context)
	return nil
}

// Returns true if Command.Name or Command.ShortName matches given name
func (c Command) HasName(name string) bool {
	return c.Name == name || c.ShortName == name
}

func (c Command) startApp(ctx *Context) error {
	app := NewApp()

	// set the name and usage
	app.Name = fmt.Sprintf("%s %s", ctx.App.Name, c.Name)
	if c.Description != "" {
		app.Usage = c.Description
	} else {
		app.Usage = c.Usage
	}

	// set CommandNotFound
	app.CommandNotFound = ctx.App.CommandNotFound

	// set the flags and commands
	app.Commands = c.Subcommands
	app.Flags = c.Flags
	app.HideHelp = c.HideHelp

	// bash completion
	app.EnableBashCompletion = ctx.App.EnableBashCompletion
	if c.BashComplete != nil {
		app.BashComplete = c.BashComplete
	}

	// set the actions
	app.Before = c.Before
	if c.Action != nil {
		app.Action = c.Action
	} else {
		app.Action = helpSubcommand.Action
	}

	return app.RunAsSubcommand(ctx)
}
