// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Creates and initializes a wio project. It also works as an updater when called on already created projects.
package create

import (
    "path/filepath"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils/types"
)

// Executes the create command provided configuration packet
func Execute(args types.CliArgs) {
    io.Norm.Yellow("Project Name: ")
    io.Norm.Cyan(filepath.Base(args.Directory) + "\n")
    io.Norm.Yellow("Project Type: ")
    io.Norm.Cyan(args.AppType + "\n")
    io.Norm.Yellow("Project Path: ")
    io.Norm.Cyan(args.Directory + "\n")
    io.Norm.White("\n")

    var projectType ProjectTypes = App{args: &args}

    if args.AppType == "lib" {
        projectType = Lib{args: &args}
    }

    // update only if asked to do that
    if args.Update {
        io.Norm.Yellow("Updating the project  ... ")
        err := projectType.update()

        if err != nil {
            io.Norm.Red("[failure]\n")
            io.Verb.Error(err.Error() + "\n")
        } else {
            io.Norm.Green("[success]\n")
            io.Norm.Yellow("Project has been successfully updated!!\n")
        }

        return
    }

    io.Norm.Yellow("Creating project structure ... ")
    err := projectType.createStructure()

    if err != nil {
        io.Norm.Red("[failure]\n")
        io.Verb.Error(err.Error() + "\n")
    } else {
        io.Norm.Green("[success]\n")
        projectType.printProjectStructure()
    }

    io.Norm.White("\n")
    io.Norm.Yellow("Creating template project ... ")
    err = projectType.createTemplateProject()

    if err != nil {
        io.Norm.Red("[failure]\n")
        io.Verb.Error(err.Error() + "\n")
    } else {
        io.Norm.Green("[success]\n")
        io.Norm.White("\n")
        io.Norm.Yellow("Project has been successfully created and initialized!!\n")
        io.Norm.Green("Check following commands: \n")
        projectType.printNextCommands()
    }
}
