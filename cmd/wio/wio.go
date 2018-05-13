// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Package main contains the main code for Wio.
// Wio is a tool to make development of embedded system applications easier and simpler.
// It allows for building, testing, and uploading AVR applications for Commandline.
 package main

import (
    "wio/cmd/wio/utils/io"
    "github.com/urfave/cli"
    "time"
    "os"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/commands/create"
    "wio/cmd/wio/commands/libraries"
    "wio/cmd/wio/utils/io/log"
    "wio/cmd/wio/types"
    "wio/cmd/wio/commands/build"
)

//go:generate go-bindata -nomemcopy -prefix ../../ ../../assets/config/... ../../assets/templates/...
func main()  {
    // override help template
    cli.AppHelpTemplate =
`Wio a simplified development process for embedded applications.
Create, Build, Test, and Upload AVR projects from Commandline.

Common Commands:
    
    wio create <project type> [options] <output directory>
        Create a new Wio project in the specified directory.
    
    wio build [options]
        Build the Wio project based on all the configurations defined
    
    wio upload [options]
        Upload the Wio project to an attached embedded device
    
    wio run [options]
        Builds, Tests, and Uploads the Wio projects

Usage: {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
Global options:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
Available commands:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
Global options:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

Copyright:
   {{.Copyright}}
   {{end}}{{if .Version}}
Vesrion:
   {{.Version}}
   {{end}}
Run "wio command <help>" for more information about a command.
`

cli.CommandHelpTemplate =
`{{.Usage}}

Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}

Category:
   {{.Category}}{{end}}{{if .Description}}

Description:
   {{.Description}}{{end}}{{if .VisibleFlags}}

Available commands:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
Run "wio help" to see global options.
`

cli.SubcommandHelpTemplate =
`{{.Usage}}

Usage: {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} command{{if .VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

Available commands:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{if .VisibleFlags}}
Options:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
Run "wio help" to see global options.
`
    // get default configuration values
    defaults := types.DConfig{}
    err := io.AssetIO.ParseYml("config/defaults.yml", &defaults)
    if err != nil {
        log.Error(false, err.Error())
    }

    // command that will be executed
    var command commands.Command

    app := cli.NewApp()
    app.Name = "wio"
    app.Version = defaults.Version
    app.EnableBashCompletion = true
    app.Compiled = time.Now()
    app.Copyright = "Copyright (c) 2018 Waterloop"
    app.Usage = "Create, Build and Upload AVR projects"

    app.Flags = []cli.Flag {
        cli.BoolFlag{Name: "verbose",
            Usage: "Turns verbose mode on to show detailed errors and commands being executed",
            },
    }

    app.Commands = []cli.Command{
        {
            Name:  "create",
            Usage: "Creates and initializes a wio project.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "pkg",
                    Usage:     "Creates a wio package, intended to be used by other people",
                    UsageText: "wio create pkg <DIRECTORY> <BOARD> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: defaults.Ide},
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: defaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: defaults.Platform},
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.PKG, Update: false}
                    },
                },
                cli.Command{
                    Name:      "app",
                    Usage:     "Creates a wio application, intended to be compiled and uploaded to a device",
                    UsageText: "wio create app <DIRECTORY> <BOARD> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: defaults.Ide},
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: defaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: defaults.Platform},
                        cli.BoolFlag{Name: "tests",
                            Usage: "Creates a test folder to support unit testing",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.APP, Update: false}
                    },
                },
            },
        },
        {
            Name:  "update",
            Usage: "Updates the current project and fixes any issues.",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "pkg",
                    Usage:     "Updates a wio package, intended to be used by other people",
                    UsageText: "wio update pkg <DIRECTORY> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "board",
                            Usage: "Board being used for this project. This will use this board for the update",
                            Value: defaults.Board},
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.PKG, Update: true}
                    },
                },
                cli.Command{
                    Name:      "app",
                    Usage:     "Updates a wio application, intended to be compiled and uploaded to a device",
                    UsageText: "wio update app <DIRECTORY> [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "board",
                            Usage: "Board being used for this project. This will use this board for the update",
                            Value: defaults.Board},
                        cli.StringFlag{Name: "ide",
                            Usage: "Creates the project for a specified IDE (CLion, Eclipse, VS Code)",
                            Value: defaults.Ide},
                        cli.StringFlag{Name: "framework",
                            Usage: "Framework being used for this project. Framework contains the core libraries",
                            Value: defaults.Framework},
                        cli.StringFlag{Name: "platform",
                            Usage: "Platform being used for this project. Platform is the type of chip supported (AVR/ ARM)",
                            Value: defaults.Platform},
                        cli.BoolFlag{Name: "tests",
                            Usage: "Creates a test folder to support unit testing",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = create.Create{Context: c, Type: create.APP, Update: true}
                    },
                },
            },
        },
        {
            Name:      "build",
            Usage:     "Builds the project",
            UsageText: "wio build [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "target",
                    Usage: "Build a specified target instead of building all the targets",
                    Value: defaults.Btarget,
                },
                cli.StringFlag{Name: "directory",
                    Usage: "Directory for the project (default: current working directory)",
                    Value: getCurrDir(),
                },
            },
            Action: func(c *cli.Context) {
                command = build.Build{Context: c}
            },
        },
        {
            Name:      "clean",
            Usage:     "Cleans all the build files for the project",
            UsageText: "wio clean",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "upload",
            Usage:     "Uploads the project to a device",
            UsageText: "wio upload [command options]",
            Flags: []cli.Flag{
                cli.StringFlag{Name: "file",
                    Usage: "Hex file can be provided to upload; program will upload that file",
                    Value: defaults.File,
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to",
                    Value: defaults.Port,
                },
                cli.StringFlag{Name: "target",
                    Usage: "Uploads a specified target instead of the main/default target",
                    Value: defaults.Utarget,
                },
            },
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "run",
            Usage:     "Builds, Tests, and Uploads the project to a device",
            UsageText: "wio run [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "file",
                    Usage: "Hex file can be provided to upload; program will upload that file",
                    Value: defaults.File,
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: defaults.Port,
                },
                cli.StringFlag{Name: "target",
                    Usage: "Builds, and uploads a specified target instead of the main/default target",
                    Value: defaults.Utarget,
                },
            },
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "test",
            Usage:     "Runs unit tests available in the project",
            UsageText: "wio test",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "clean",
                    Usage: "Clean the project before building it",
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: defaults.Port,
                },
                cli.StringFlag{Name: "target",
                    Usage: "Builds, and uploads a specified target instead of the main/default target",
                    Value: defaults.Utarget,
                },
            },
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "monitor",
            Usage:     "Runs the serial monitor",
            UsageText: "wio monitor [command options]",
            Flags: []cli.Flag{
                cli.BoolFlag{Name: "gui",
                    Usage: "Runs the GUI version of the serial monitor tool",
                },
                cli.StringFlag{Name: "port",
                    Usage: "Port to upload the project to, (default: automatically select)",
                    Value: defaults.Port,
                },
            },
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "doctor",
            Usage:     "Show information about the installed tooling",
            UsageText: "wio doctor",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "configure",
            Usage:     "Configures paths for the tools used for development",
            UsageText: "wio configure [command options]",
            Flags: []cli.Flag{
                cli.StringFlag{Name: "arduino-sdk-dir",
                    Usage: "path to Arduino SDK",
                },
                cli.StringFlag{Name: "make-path",
                    Usage: "path to `make` tool",
                },
                cli.StringFlag{Name: "cmake-path",
                    Usage: "Path to `cmake` tool",
                },
                cli.StringFlag{Name: "avr-path",
                    Usage: "Path to AVR libraries",
                },
                cli.StringFlag{Name: "arm-path",
                    Usage: "Path to ARM libraries",
                },
            },
            Action: func(c *cli.Context) error {
                // If no flag provided, show current settings
                return nil
            },
        },
        {
            Name:      "analyze",
            Usage:     "Analyzes C/C++ code statically",
            UsageText: "wio analyze",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:      "doxygen",
            Usage:     "Runs doxygen tool to create documentation for the code",
            UsageText: "wio doxygen",
            Action: func(c *cli.Context) error {
                return nil
            },
        },
        {
            Name:  "libraries",
            Usage: "Libraries (package) manager for Wio projects",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:  "get",
                    Usage: "Gets all the libraries mentioned in wio.yml file and vendor folder",
                    UsageText: "wio libraries get [command options]",
                    Flags: []cli.Flag{
                        cli.BoolFlag{Name: "clean",
                            Usage: "Cleans all the current packages and re get all of them",
                        },
                        cli.StringFlag{Name: "version_control",
                            Usage: "Specify the version control tool to usage",
                            Value: "git",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = libraries.Libraries{Context: c, Type: libraries.GET}
                    },
                },
                cli.Command{
                    Name:  "update",
                    Usage: "Updates all the libraries mentioned in wio.yml file and vendor folder",
                    UsageText: "wio libraries update [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "version_control",
                            Usage: "Specify the version control tool to usage",
                            Value: "git",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = libraries.Libraries{Context: c, Type: libraries.UPDATE}
                    },
                },
                cli.Command{
                    Name:  "collect",
                    Usage: "Creates vendor folder and puts all the libraries in that folder",
                    UsageText: "wio libraries collect [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "path",
                            Usage: "Path to collect a library instead of collecting all of them",
                            Value: "none",
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = libraries.Libraries{Context: c, Type: libraries.COLLECT}
                    },
                },
            },
        },
        {
            Name:  "tool",
            Usage: "Contains various tools related to setup, initialize and upgrade of Wio",
            Subcommands: cli.Commands{
                cli.Command{
                    Name:      "setup",
                    Usage:     "When tool is newly installed, it sets up the tool for the machine",
                    UsageText: "wio setup",
                    Action: func(c *cli.Context) {
                        command = libraries.Libraries{Context: c}
                    },
                },
                cli.Command{
                    Name:      "upgrade",
                    Usage:     "Upgrades the current version of the program",
                    UsageText: "wio upgrade [command options]",
                    Flags: []cli.Flag{
                        cli.StringFlag{Name: "version",
                            Usage: "Specify the exact version to upgrade/downgrade wio to",
                            Value: defaults.Version,
                        },
                    },
                    Action: func(c *cli.Context) {
                        command = libraries.Libraries{Context: c}
                    },
                },
            },
        },
    }

    app.Action = func(c *cli.Context) error {
        app.Command("help").Run(c)
        return nil
    }

    if err = app.Run(os.Args); err != nil {
        panic(err)
    }

    // execute the command
    if command != nil {
        command.Execute()
    }
}

func getCurrDir() (string) {
    directory, err := os.Getwd()
    commands.RecordError(err, "")
    return directory
}

// Set's verbose mode on
func turnVerbose(value bool) {
    if value == true {
        log.SetVerbose()
    }
}
